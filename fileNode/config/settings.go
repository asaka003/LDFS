package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var (
	//Limit     bool   //视频文件限流配置
	//SystemPre string //服务节点编号
	NodeName string //服务节点名称

	AddrConsule []string
	//AddrRpcLocalServer  string
	AddrHttpLocalServer string
	//AddrHttpApiServer   string
	LocalServerName string
	//NodeNumber      int

	//TempFileDir     string
	FilesDir        string
	VideoFaceImgDir string
	VideoFilesDir   string
	FileShardDir    string

	NameRpcVideoServer string

	FileNodeUrls     []string
	LocalFileNodeUrl string
)

func config_init() {
	//配置视频文件限流信息
	//Limit = viper.GetBool("video_ini.limit")
	NodeName = viper.GetString("node.name") //加载服务节点名称

	//AddrHttpApiServer = viper.GetString("Servers.ApiServer.IP")
	//初始化服务地址以及服务名称
	//SystemPre = viper.GetString("node.system_num") //服务节点编号
	AddrConsule = viper.GetStringSlice("Servers.Consule.IP")
	AddrHttpLocalServer = viper.GetString("node.http_ip")
	//AddrRpcLocalServer = viper.GetString("node.rpc_ip")
	LocalServerName = viper.GetString("node.name")
	//NodeNumber = viper.GetInt("node.number")
	NameRpcVideoServer = viper.GetString("Servers.VideoServer.Name")

	//读取所有存储节点的http服务URL地址信息
	FileNodeUrls = viper.GetStringSlice("Nodes.Urls")
	LocalFileNodeUrl = viper.GetString("Nodes.LocalUrl")
}

// //初始化配置读取器viper
// func ViperInit() error {
// 	viper.SetConfigName("config")
// 	viper.SetConfigType("yaml")
// 	viper.AddConfigPath("./config")
// 	err := viper.ReadInConfig()
// 	if err != nil {
// 		fmt.Printf("viper.ReadConfig() failed : %s", err)
// 		return err
// 	}

// 	viper.SetConfigName("ServerAddr")
// 	viper.SetConfigType("yaml")
// 	viper.AddConfigPath("./config")
// 	err = viper.MergeInConfig()
// 	if err != nil {
// 		fmt.Printf("viper.ReadConfig() failed : %s", err)
// 		return err
// 	}

// 	viper.WatchConfig()
// 	viper.OnConfigChange(func(in fsnotify.Event) {
// 		config_init()
// 		//设置配置文件修改后的动态绑定事件......
// 		fmt.Println("配置文件改变")
// 	})

// 	//加载配置信息
// 	config_init()
// 	FileDirInit()
// 	return err
// }

func createDir(dir string) (err error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// 如果目录不存在，则创建
		err = os.MkdirAll(dir, 0755) // 0755 表示具有读/写/执行权限的所有者，以及读/执行权限的其他用户。
		if err != nil {
			fmt.Println("创建目录失败:", err)
			return err
		}
	}
	return
}

//文件存储路径初始化
func FileDirInit() {
	// TempFileDir = viper.GetString("dir.tempFileDir")
	FilesDir = viper.GetString("dir.fileDir")
	VideoFaceImgDir = viper.GetString("dir.videoFaceImgDir")
	VideoFilesDir = viper.GetString("dir.videoDir")
	FileShardDir = viper.GetString("dir.fileShardDir")

	//createDir(TempFileDir)
	createDir(FilesDir)
	createDir(VideoFaceImgDir)
	createDir(VideoFilesDir)
	createDir(FileShardDir)
}
