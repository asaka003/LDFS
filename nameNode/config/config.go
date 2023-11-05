package config

import (
	"LDFS/model"
	"fmt"
	"math/rand"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	//读取文件个各节点的IP地址
	DataNodeList []model.DataNode
	DataNodeCli  *DataNodeclient

	Mysql  *MysqlConfig
	Redis  *RedisConfig
	Sqlite *SqliteConfig

	SystemDB string

	MultiUploadDir string
	FileMetaDir    string

	ECDataShardNum   int64
	ECParityShardNum int64
	CopyReplicasNum  int64
)

const (
	RemainSize int64 = 100 * 1024 * 1024
)

//初始化配置读取器viper
func viperInit() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("viper.ReadConfig() failed : %s", err)
		return err
	}

	viper.SetConfigName("ServerAddr")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err = viper.MergeInConfig()
	if err != nil {
		fmt.Printf("viper.ReadConfig() failed : %s", err)
		return err
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		//设置配置文件修改后的动态绑定事件......
		fmt.Println("配置文件改变")
	})
	return nil
}

//初始化所有配置信息
func ConfigInit() (err error) {
	err = viperInit()
	if err != nil {
		return err
	}
	Mysql = &MysqlConfig{
		Host:         viper.GetString("mysql.host"),
		Port:         viper.GetString("mysql.port"),
		User:         viper.GetString("mysql.user"),
		Password:     viper.GetString("mysql.password"),
		DB:           viper.GetString("mysql.dbname"),
		MaxOpenConns: viper.GetInt("mysql.maxOpenConns"),
		MaxIdleConns: viper.GetInt("mysql.maxIdleConns"),
	}
	Redis = &RedisConfig{
		Host:     viper.GetString("redis.host"),
		Password: viper.GetString("redis.password"),
		Port:     viper.GetString("redis.port"),
		PoolSize: viper.GetInt("redis.poolSize"),
		DB:       viper.GetInt("redis.db"),
	}

	Sqlite = &SqliteConfig{
		Dialect: viper.GetString("sqlite.dialect"),
		DbFile:  viper.GetString("sqlite.dbfile"),
	}

	SystemDB = viper.GetString("SYSTEM_DB")
	MultiUploadDir = viper.GetString("MultiUploadDir")
	FileMetaDir = viper.GetString("FileMetaDir")
	DataNodeCli = &DataNodeclient{}
	DataNodeUrls := viper.GetStringSlice("Nodes.List")
	for _, url := range DataNodeUrls {
		//请求DataNode 磁盘存储情况
		dataNode, err := DataNodeCli.GetStorageInfo(url)
		if err != nil {
			panic(err)
		}
		DataNodeList = append(DataNodeList, dataNode)
	}

	ECDataShardNum = viper.GetInt64("EC.dataShards")
	ECParityShardNum = viper.GetInt64("EC.parityShards")
	CopyReplicasNum = viper.GetInt64("Copy.replicasNum")

	// 种子随机数生成器
	rand.Seed(time.Now().UnixNano())

	return
}
