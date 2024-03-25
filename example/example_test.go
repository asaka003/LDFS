package example

import (
	storagesdk "LDFS/storage-sdk"
	"fmt"
	"testing"
)

func Test_upload1111121231123133(t *testing.T) {
	client := storagesdk.NewReplicasClient(
		[]string{
			"http://localhost:11001",
		},
		[]string{
			"http://localhost:11002",
			"http://localhost:11003",
		})

	err := client.SimpleUploadFile("dir/img.png", "test.png")
	if err != nil {
		fmt.Println(err)
	}
}

func Test_download1001223(t *testing.T) {
	client := storagesdk.NewReplicasClient(
		[]string{
			"http://localhost:11001",
		},
		[]string{
			"http://localhost:11002",
			"http://localhost:11003",
		})

	err := client.DownloadFile("dir/img.png", "download.png")
	if err != nil {
		fmt.Println(err)
	}
}
