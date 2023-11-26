package util

import (
	"fmt"
	"testing"

	"github.com/shirou/gopsutil/disk"
)

func TestXxx(t *testing.T) {
	//GetDiskUsageInfo()
	// size, err := GetDirectorySize("D:\\材料文档\\南大作业")
	// fmt.Println(size, err)
	// str := filepath.VolumeName("/etc/passwd/test/666")
	// fmt.Println(str == "")
	var total int64
	var free int64
	var used int64

	usage, err := disk.Usage("D:\\材料文档\\南大作业")
	if err != nil {
		fmt.Println(err)
	}
	total += int64(usage.Total)
	free += int64(usage.Free)
	used += int64(usage.Used)
	fmt.Println(total, free, used)
}
