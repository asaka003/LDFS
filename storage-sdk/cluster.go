package storagesdk

import (
	"stathat.com/c/consistent"
)

type NameNodeCluster struct {
	Consistent *consistent.Consistent
}

var (
	nameNodeLeaderCluster   *NameNodeCluster //leader集群
	nameNodeFollowerCluster *NameNodeCluster //follower集群
)

func InitCluster(nameNodeLeaderUrls, nameNodeFollowerUrls []string) {
	nameNodeLeaderCluster = &NameNodeCluster{
		Consistent: consistent.New(),
	}
	for _, node := range nameNodeLeaderUrls {
		nameNodeLeaderCluster.Consistent.Add(node)
	}

	nameNodeFollowerCluster = &NameNodeCluster{
		Consistent: consistent.New(),
	}
	for _, node := range nameNodeFollowerUrls {
		nameNodeFollowerCluster.Consistent.Add(node)
	}
}
