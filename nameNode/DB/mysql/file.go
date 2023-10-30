package mysql

import (
	"LDFS/nameNode/model"
	"database/sql"
	"errors"
	"log"

	"go.uber.org/zap"
)

var (
	ErrFileExist    = errors.New("文件已存在")
	ErrFileNotExist = errors.New("文件不存在")
	ErrEmptyRows    = errors.New("sql执行行数为0")
)

/*
	mysql多个条件查询时，只会使用到一个索引，所以条件查询要么建立一个复合索引，要么只使用第一个条件的索引进行查询，其他条件则采用临时表扫描的形式
*/

//查询文件是否存在
func CheckFile(file_hash string) (file_id int64, ok bool, err error) {
	sql_str := `select ID from file_meta where file_hash = ?`
	if err := DB.Get(&file_id, sql_str, file_hash); err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return -1, false, nil
		}
		log.Println("mysql.CheckFile() failed..", zap.Error(err))
		return -1, false, err
	}
	return file_id, true, nil
}

//添加单个文件meta信息
func AddFileMeta(FileHash string, FileSize int64, FileURL string, FileType int) error {
	sql := `insert into file_meta(file_hash,file_size,file_url,file_type) values(?,?,?,?)`
	_, err := DB.Exec(sql, FileHash, FileSize, FileURL, FileType)
	if err != nil {
		log.Println("mysql.AddFileMeta() failed..", zap.Error(err))
		return err
	}
	return nil
}

//查询单个文件meta信息
func GetFileMeta(file_id int64) (*model.FileMeta, error) {
	file := new(model.FileMeta)
	sql := `select ID,file_hash,file_size,file_url from file_meta where ID = ? limit 0,1`
	if err := DB.Get(file, sql, file_id); err != nil {
		log.Println("mysql.GetFileMeta() failed..", zap.Error(err))
		return nil, err
	}
	return file, nil
}

//删除单个文件meta信息
func DeleteFileMeta(file_id int64) error {
	sql := `DELETE FROM file_meta WHERE ID = ?`
	_, err := DB.Exec(sql, file_id)
	if err != nil {
		log.Println("mysql.DeleteFileMeta() failed ..", zap.Error(err))
		return err
	}
	return nil
}

//添加用户文件信息
func AddFileMetaToUser(user_id int64, file_id int64, filename string) error {
	sql := `insert ignore into user_file (user_id,file_id,file_name) values (?,?,?)`
	res, err := DB.Exec(sql, user_id, file_id, filename)
	if err != nil {
		log.Println("mysql.UserFileUpload() failed..", zap.Error(err))
		return err
	}
	rows_count, err := res.RowsAffected()
	if err != nil {
		log.Println("mysql.UserFileUpload() failed..", zap.Error(err))
	} else if rows_count <= 0 {
		return ErrFileExist
	}
	return nil
}

/*
	在使用联结查询时，尽量将条件写在ON中，因为联立的笛卡尔积是在ON筛选之后联立的，而WHERE是在笛卡尔积之后进行筛选，
	如果数据量特别庞大，笛卡尔积所造成的计算资源消耗是难以接受的。
*/

//查询用户单个文件信息
func GetUserFileMeta(user_id int64, file_id int64) (file_meta *model.UserFileMeta, err error) {
	file_meta = new(model.UserFileMeta)
	sql_str := `
			SELECT file_id,file_name,file_hash,file_size,file_url,user_file.update_at
			FROM file_meta INNER JOIN user_file ON user_id = ? and status = 1 and id = ? and user_file.file_id = file_meta.ID
			`

	//若文件不存在则会报错
	if err := DB.Get(file_meta, sql_str, user_id, file_id); err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, err
		}
		log.Println("mysql.GetUserFileMeta() failed .. ", zap.Error(err))
		return nil, err
	}
	return
}

/*
	查询用户的所有文件信息，这里采用的是联立表的形式进行查询，采用内联方式，
	如果file_meta表使用了没有索引的条件进行where查询，就会对file_mate表进行一个全表扫描，
	这是不能够接受的，因此在内联查询时，一个表要么使用只有索引的条件进行where查询，要么不使用
	where查询，只需要内联的条件有索引就行。这样优化器在执行的过程中，就会先使用user_file的user_id
	先进行一个索引查询，然后使用file_meta的ID进行索引内联，不存在全表扫描的情况。

	另一种方式是采用多次查询的方法进行查询，例如先查询指定用户ID的所有文件ID，然后将这些ID拼接使用in
	关键字进行查询，但是in字段只有在少量范围的时候才会使用到索引，当用户in字段里面的范围比较多的时候，
	优化器就会选择全表扫描的方式查询。还有一种就是循环遍历所有文件ID，进行多次单次查询，这种查询的原理其实和
	内联查询差不多，只不过这种会有数据库和应用之间传输的资源消耗，而内联查询是占用mysql数据库的计算资源。

	但是单次查询这种方式容易使用到缓存，因为单次查询某个文件的概率很大，同时循环单次查询的方式也方便维护和扩展(因为应用服务器的负载是很好拓展的，而数据库往往是整个系统的瓶颈)，
	减少数据库的压力，同时减少加锁的时间，因为数据库联立查询会可能会对多个表进行加锁，在高并发场景下产生的影响较大。
*/
//查询用户文件信息列表
func GetUserAllFileMeta(user_id int64, status int) ([]*model.UserFileMeta, error) {
	files := make([]*model.UserFileMeta, 0)
	sql := `SELECT file_id,file_name,file_hash,file_size,user_file.update_at,file_url 
			FROM user_file INNER JOIN file_meta ON user_id = ? and status = ? and user_file.file_id = file_meta.ID
			`
	if err := DB.Select(&files, sql, user_id, status); err != nil {
		log.Println("mysql.GetUserAllFileMeta() failed...", zap.Error(err))
		return nil, err
	}
	return files, nil
}

//查询用户视频文件信息列表( 待优化...     file_meta查询使用的索引是file_type,尽量减少三表以上的join)
func GetUserVideoFiles(user_id int64) ([]*model.UserVideoMeta, error) {
	files := make([]*model.UserVideoMeta, 0)
	sql := `
		SELECT file_id,file_meta.file_hash,file_name,file_size,face_url,user_file.update_at
		FROM user_file 
		INNER JOIN file_meta ON user_id = ? and status = 1 and file_type = 1 and file_meta.ID = user_file.file_id
		INNER JOIN videos_img ON file_meta.file_hash = videos_img.file_hash
		`
	if err := DB.Select(&files, sql, user_id); err != nil {
		log.Println("mysql.GetUserVideoFiles() failed...", zap.Error(err))
		return nil, err
	}
	return files, nil
}

//更新用户文件信息
func UpdateFileInfo(user_id int64, file_id int64, file_name string) (err error) {
	sql := `UPDATE user_file set file_name = ? where user_id = ? and file_id = ?`
	result, err := DB.Exec(sql, file_name, user_id, file_id)
	if err != nil {
		log.Println("", zap.Error(err))
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("", zap.Error(err))
		return
	}
	if rows <= 0 {
		log.Println("", zap.Error(ErrFileNotExist))
		return ErrFileNotExist
	}
	return
}

//添加视频文件封面信息
func InsertFaceImg(img_url string, file_hash string) (err error) {
	sql := `insert into videos_img (file_hash,face_url) values (?,?)`
	if _, err = DB.Exec(sql, file_hash, img_url); err != nil {
		log.Println("插入视频封面信息失败", zap.Error(err))
	}
	return
}

//更新视频文件封面信息
func UpdateFaceImg(img_url string, file_hash string) (err error) {
	sql := `update videos_img set face_url = ? where file_hash = ?`
	if _, err = DB.Exec(sql, img_url, file_hash); err != nil {
		log.Println("插入视频封面信息失败", zap.Error(err))
	}
	return
}

//删除用户单个文件meta信息
func DeleteUserFileMeta(file_id int64, user_id int64) error {
	sql := "update user_file set status=0 where user_id =? and file_id=?"
	_, err := DB.Exec(sql, user_id, file_id)
	if err != nil {
		log.Println("mysql.DeleteFileMeta() failed ..", zap.Error(err))
		return err
	}

	return nil
}

//回收用户文件
func RecycleUserFile(file_id int64, user_id int64) (err error) {
	sql := `update user_file set status = 1 where user_id = ? and file_id = ?`
	_, err = DB.Exec(sql, user_id, file_id)
	if err != nil {
		log.Println("", zap.Error(err))
	}
	return
}

//用户提交上传文件信息
func UserSubmitUpload(user_id int64, file_size int64, file_name string, file_url string, file_hash string, file_type int) (err error) {
	//开启事务
	tx, err := DB.Beginx()
	if err != nil {
		log.Println(err)
	}

	//事务执行提交
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			log.Panicln(err)
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	sql1 := `INSERT INTO file_meta (file_hash,file_size,file_url,file_type) VALUES(?,?,?,?)`
	rs, err := tx.Exec(sql1, file_hash, file_size, file_url, file_type)
	if err != nil {
		log.Println(err)
		return err
	}
	n, err := rs.RowsAffected()
	if err != nil {
		log.Println(err)
		return err
	}
	if n != 1 {
		return ErrEmptyRows
	}

	sql2 := `SELECT file_id from file_meta where file_hash = ?`
	var file_id int64
	err = tx.Get(&file_id, sql2, file_hash)
	if err != nil {
		log.Println(err)
		return err
	}

	sql3 := `INSERT INTO user_file (user_id,file_id,file_name) VALUES(?,?,?)`
	rs, err = tx.Exec(sql3, user_id, file_id, file_name)
	if err != nil {
		log.Println(err)
		return err
	}
	n, err = rs.RowsAffected()
	if err != nil {
		log.Println(err)
		return err
	}
	if n != 1 {
		log.Println(err)
		return ErrEmptyRows
	}
	return err
}

//检查用户是否持有该文件
// func CheckUserFile(user_id int64, file_id int64) (bool, error) {
// 	sql := `select count(user_id) from user_file where user_id = ? and file_id = ?`
// 	var count int
// 	if err := DB.Get(&count, sql, user_id, file_id); err != nil {
// 		log.Println("mysql.CheckUserFile() failed..", zap.Error(err))
// 		return false, err
// 	}
// 	if count == 1 {
// 		return true, nil
// 	} else {
// 		return false, nil
// 	}
// }

//统计某文件用户持有量
func CountFileUserNum(file_id int64) (num int64, err error) {
	sql := `select count(user_id) from user_file where file_id = ?`
	err = DB.Get(&num, sql, file_id)
	if err != nil {
		log.Println("统计用户持有量失败", zap.Error(err))
	}
	return
}

//删除用户持有文件记录
func DeepDelUserFileMeta(file_id int64, user_id int64) (err error) {
	sql := `delete from user_file where user_id = ? and file_id = ?`
	_, err = DB.Exec(sql, user_id, file_id)
	if err != nil {
		log.Println("删除用户文件记录失败", zap.Error(err))
	}
	return
}
