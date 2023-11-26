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
	nameNodeCluster = &NameNodeCluster{
		Consistent: consistent.New(),
	}
	for _, node := range nodes {
		nameNodeCluster.Consistent.Add(node)
	}
}
