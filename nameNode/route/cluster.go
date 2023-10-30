package route

import (
	"strings"

	"stathat.com/c/consistent"
)

type Cluster struct {
	Name       string
	Consistent *consistent.Consistent
}

var (
	clusters []*Cluster
)

func CreateReverseProxy(Userclusters []*Cluster) {
	clusters = Userclusters
}

//根据请求路径获取该路径对应的集群节点配置
func getClusterByPath(path string) *Cluster {
	for _, cluster := range clusters {
		if strings.HasPrefix(path, cluster.Name) {
			return cluster
		}
	}
	return nil
}
