package logic

import (
	"LDFS/model"
	db "LDFS/nameNode/DB"
)

//解析是否为视频文件
// func is_video(file_ext string) bool {
// 	video_ext := []string{
// 		".mp4",
// 		".avi",
// 		".mkv",
// 		".mov",
// 		".flv",
// 	}
// 	for _, v := range video_ext {
// 		if file_ext == v {
// 			return true
// 		}
// 	}
// 	return false
// }

// func AddUserUploadFileInfo(user_id int64, file_size int64, file_name string, file_hash string, FileKey string) error {
// 	var file_type int
// 	if is_video(filepath.Ext(file_name)) {
// 		file_type = 1
// 	} else {
// 		file_type = 0
// 	}
// 	err := mysql.AddFileMeta(file_hash, file_size, FileKey, file_type)
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	file_id, _, err := mysql.CheckFile(file_hash)
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	err = mysql.AddFileMetaToUser(user_id, file_id, file_name)
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	return nil
// }

//保存加单上传文件DataNode存储信息
func SaveSampleUploadInfo(fileKey string, list []*model.Shard) (err error) {
	return db.DB.SaveFileMetaInfo(fileKey, list)
}

//查询所有的FileKeys列表
func GetAllFileKeys() {

}
