package snowflake

import (
	"time"

	"github.com/bwmarrin/snowflake"
)

var Node *snowflake.Node

func Init(stratTime string, machineID int64) (err error) {
	var start_t time.Time
	start_t, err = time.Parse("2006-01-02", stratTime)
	if err != nil {
		return
	}
	snowflake.Epoch = start_t.UnixNano() / 1000000
	Node, err = snowflake.NewNode(machineID)
	return
}

func GetID() int64 {
	return Node.Generate().Int64()
}
