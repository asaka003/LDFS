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
