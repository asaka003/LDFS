package config

import (
	"LDFS/model"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

/*
	初始化配置
*/

type DataNodeclient struct {
	GetStorageInfoUrl string
}

func (dataNodeClient *DataNodeclient) GetStorageInfo(backend string) (dataNodeInfo model.DataNode, err error) {
	res, err := http.Get(backend + dataNodeClient.GetStorageInfoUrl)
	if err != nil {
		return
	}
	defer res.Body.Close()
	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	dataNodeInfo = model.DataNode{}
	err = json.Unmarshal(resBytes, &dataNodeInfo)
	if err != nil {
		return
	}
	return
}

//
