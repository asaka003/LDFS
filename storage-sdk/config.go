package storagesdk

import (
	"LDFS/nodeClient"
)

var (
	DataNodeUrls   []string
	NameNodeUrls   []string
	NameNodeClient *nodeClient.NameNodeHttpClient
	DataNodeClient *nodeClient.DataNodeHttpClient
)

const EC_OutDir = "/tmp/ec_out/"
const EC_InputDir = "/tmp/ec_input/"
const (
	Copy_BlockSize int64 = 16 * 1024 * 1024
	EC_BlockSize   int64 = 128 * 1024 * 1024
	EC_ShardSize   int64 = 16 * 1024 * 1024
)
