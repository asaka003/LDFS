package example

import (
	storagesdk "LDFS/storage-sdk"
	"fmt"
	"testing"
)

func Test_example(t *testing.T) {
	client := storagesdk.NewReplicasClient([]string{
		"http://localhost:11001",
	})

	err := client.SimpleUploadFile("dir/img.png", "test.png")
	if err != nil {
		fmt.Println(err)
	}
}

func Test_main6(t *testing.T) {
	client := storagesdk.NewReplicasClient([]string{
		"http://localhost:11001",
	})

	err := client.SimpleUploadFile("dir/img.png", "test.png")
	if err != nil {
		fmt.Println(err)
	}
}
