package example

import (
	storagesdk "LDFS/storage-sdk"
	"fmt"
	"testing"
)

func Test_example(t *testing.T) {
	client := storagesdk.NewReplicasClient([]string{
		"http://124.223.212.153:9090",
	})

	err := client.SimpleUploadFile("dir/img.png", "test.png")
	if err != nil {
		fmt.Println(err)
	}
}
