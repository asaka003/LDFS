package example

import (
	storagesdk "LDFS/storage-sdk"
	"fmt"
	"testing"
)

func Test_ReplicasUpload(t *testing.T) {
	client := storagesdk.NewReplicasClient(
		[]string{
			"http://localhost:11001",
		},
		[]string{
			"http://localhost:11002",
			"http://localhost:11003",
		})

	err := client.SimpleUploadFile("dir/2.png", "1231.jpg")
	if err != nil {
		fmt.Println(err)
	}
}

func Test_ECUpload(t *testing.T) {
	client := storagesdk.NewECClient(
		[]string{
			"http://localhost:11001",
		},
		[]string{
			"http://localhost:11002",
			"http://localhost:11003",
		})

	err := client.SimpleUploadFile("dir/2.png", "1231.jpg")
	if err != nil {
		fmt.Println(err)
	}
}

func Test_ReplicasDownload(t *testing.T) {
	client := storagesdk.NewReplicasClient(
		[]string{
			"http://localhost:11001",
		},
		[]string{
			"http://localhost:11002",
			"http://localhost:11003",
		})

	err := client.DownloadFile("dir/2.png", "download2.png")
	if err != nil {
		fmt.Println(err)
	}
}

func Test_ECDownload(t *testing.T) {
	client := storagesdk.NewECClient(
		[]string{
			"http://localhost:11001",
		},
		[]string{
			"http://localhost:11002",
			"http://localhost:11003",
		})

	err := client.DownloadFile("dir/2.png", "download5.png")
	if err != nil {
		fmt.Println(err)
	}
}
