package util

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// GetDirectorySize 到所有文件的大小，以字节（Bytes）为单位
func GetDirectorySize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

//获取整个系统的磁盘大小
func GetSystemDiskSize() (int64, error) {
	var cmd *exec.Cmd

	if runtime.GOOS == "linux" {
		cmd = exec.Command("df", "-h")
	} else if runtime.GOOS == "windows" {
		cmd = exec.Command("wmic", "logicaldisk", "get", "size")
	} else {
		return 0, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	if runtime.GOOS == "linux" {
		return parseLinuxDiskSize(string(output))
	} else if runtime.GOOS == "windows" {
		return parseWindowsDiskSize(string(output))
	}

	return 0, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
}

// 解析Linux系统下的磁盘大小
func parseLinuxDiskSize(output string) (int64, error) {
	lines := strings.Split(output, "\n")
	if len(lines) < 2 {
		return 0, fmt.Errorf("failed to parse df output")
	}

	fields := strings.Fields(lines[1])
	if len(fields) < 2 {
		return 0, fmt.Errorf("failed to parse df output")
	}

	// 解析磁盘总大小（以字节为单位）
	diskSizeStr := fields[1]
	diskSize, err := parseSizeString(diskSizeStr)
	if err != nil {
		return 0, err
	}

	return diskSize, nil
}

// 解析Windows系统下的磁盘大小
func parseWindowsDiskSize(output string) (diskSize int64, err error) {
	lines := strings.Split(output, "\n")
	if len(lines) < 2 {
		return 0, fmt.Errorf("failed to parse wmic output")
	}

	// 解析磁盘总大小（以字节为单位）
	for i := 1; i < len(lines); i++ {
		diskSizeStr := strings.TrimSpace(lines[i])
		if diskSizeStr == "" {
			continue
		}
		diskSize, err = strconv.ParseInt(diskSizeStr, 10, 64)
		if err != nil {
			return 0, err
		}
	}
	return diskSize, nil
}

// 解析带有单位的大小字符串，如 "10G"、"200M" 等
func parseSizeString(sizeStr string) (int64, error) {
	sizeUnit := sizeStr[len(sizeStr)-1]
	sizeValueStr := sizeStr[:len(sizeStr)-1]
	sizeValue, err := strconv.ParseFloat(sizeValueStr, 64)
	if err != nil {
		return 0, err
	}

	switch sizeUnit {
	case 'G':
		return int64(sizeValue * 1024 * 1024 * 1024), nil
	case 'M':
		return int64(sizeValue * 1024 * 1024), nil
	case 'K':
		return int64(sizeValue * 1024), nil
	default:
		return 0, fmt.Errorf("unsupported size unit: %c", sizeUnit)
	}
}
