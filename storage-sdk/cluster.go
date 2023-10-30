package storagesdk

import (
	"stathat.com/c/consistent"
)

type NameNodeCluster struct {
	Consistent *consistent.Consistent
}

var (
	nameNodeCluster *NameNodeCluster
)

func InitCluster(nodes []string) {
	for _, node := range nodes {
		nameNodeCluster.Consistent.Add(node)
		nameNodeCluster.Consistent.Add(node)
	}
}

/*------------------旧代码----------------------------*/
// type Cluster struct {
// 	Name       string
// 	Consistent *consistent.Consistent
// }

// var (
// 	clusters []*Cluster
// )

// func createReverseProxy(Userclusters []*Cluster) {
// 	clusters = Userclusters
// }

// //根据请求路径获取该路径对应的集群节点配置
// func getClusterByPath(path string) *Cluster {
// 	for _, cluster := range clusters {
// 		if strings.HasPrefix(path, cluster.Name) {
// 			return cluster
// 		}
// 	}
// 	return nil
// }

// //初始化一致性hash路由配置
// func InitReverseProxy(nodes []string) {
// 	// 添加集群
// 	clusters := []*Cluster{
// 		{
// 			Name:       "/LDFS/nameNode-multi/",
// 			Consistent: consistent.New(),
// 		},
// 		{
// 			Name:       "/LDFS/nameNode/",
// 			Consistent: consistent.New(),
// 		},
// 	}

// 	// 添加集群api-1的后端服务器
// 	for _, node := range nodes {
// 		clusters[0].Consistent.Add(node)
// 		clusters[1].Consistent.Add(node)
// 	}

// 	createReverseProxy(clusters)
// }
