package main

import (
	"LDFS/dataNode/util"
	"bytes"
	"fmt"
	"testing"
)

func TestGetSystemDiskSize(t *testing.T) {
	size, err := util.GetSystemDiskSize()
	if err != nil {
		t.Errorf("Error: %v", err)
	} else {
		fmt.Printf("System Disk Size: %d bytes\n", size)
	}
}

func TestReaderWriter(t *testing.T) {
	//var reader *bytes.Buffer = bytes.NewBuffer(make([]byte, 1024))
	var dst []*bytes.Buffer = []*bytes.Buffer{
		bytes.NewBuffer([]byte("New Data1")),
		bytes.NewBuffer([]byte("New Data2")),
	}
	for _, buf := range dst {
		// 保存原始数据
		originalData := buf.Bytes()
		// 重置 *bytes.Buffer
		buf.Reset()
		// 现在可以重新使用 buf，例如写入数据
		buf.Write([]byte("New Data"))
		fmt.Println(string(originalData))
		fmt.Println(string(buf.Bytes()))
	}
}
