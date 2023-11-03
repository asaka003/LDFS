package storagesdk

import (
	"bytes"
	"errors"
	"io"
	"os"

	"github.com/klauspost/reedsolomon"
)

var ErrReconstruct error = errors.New("verification failed after reconstruction, data likely corrupted")

// func saveFileMeta(meta FileMetadata) error {
// 	data, err := json.Marshal(meta)
// 	if err != nil {
// 		return err
// 	}

// 	err = redis.SaveFileMeta(meta.UUID, string(data)) //保存文件元信息
// 	if err != nil {
// 		return err
// 	}
// 	err = redis.CreateDirByFileKey(meta.FileKey, meta.UUID) //保存文件目录
// 	return err
// }

// func encodeFile(filePath, outDir string, dataShards, parityShards int) (FileShardsPath []string, err error) {
// 	// Create encoding matrix.
// 	enc, err := reedsolomon.NewStream(dataShards, parityShards)
// 	f, err := os.Open(filePath)
// 	if err != nil {
// 		return
// 	}

// 	instat, err := f.Stat()
// 	if err != nil {
// 		return
// 	}

// 	shards := dataShards + parityShards
// 	out := make([]*os.File, shards)
// 	FileShardsPath = make([]string, shards)

// 	// Create the resulting files.
// 	dir, file := filepath.Split(filePath)
// 	if outDir != "" {
// 		dir = outDir
// 	}
// 	for i := range out {
// 		outfn := fmt.Sprintf("%s.%d", file, i)
// 		FileShardsPath[i] = filepath.Join(dir, outfn)
// 		out[i], err = os.Create(FileShardsPath[i])
// 		if err != nil {
// 			return
// 		}
// 	}

// 	// Split into files.
// 	data := make([]io.Writer, dataShards)
// 	for i := range data {
// 		data[i] = out[i]
// 	}
// 	// Do the split
// 	err = enc.Split(f, data, instat.Size())
// 	if err != nil {
// 		return
// 	}

// 	// Close and re-open the files.
// 	input := make([]io.Reader, dataShards)

// 	for i := range data {
// 		out[i].Close()
// 		f, err = os.Open(out[i].Name())
// 		if err != nil {
// 			return
// 		}
// 		input[i] = f
// 	}

// 	// Create parity output writers
// 	parity := make([]io.Writer, parityShards)
// 	for i := range parity {
// 		parity[i] = out[dataShards+i]
// 		defer out[dataShards+i].Close()
// 	}

// 	// Encode parity
// 	err = enc.Encode(input, parity)
// 	if err != nil {
// 		return
// 	}
// 	return
// }

func EncodeBuffer(buffer *bytes.Buffer, dataShards, parityShards int) (BufferShards []*bytes.Buffer, err error) {
	// Create encoding matrix.
	enc, err := reedsolomon.NewStream(dataShards, parityShards)
	if err != nil {
		return
	}

	shards := dataShards + parityShards
	BufferShards = make([]*bytes.Buffer, shards)
	for i := range BufferShards {
		buf := make([]byte, EC_ShardSize)
		BufferShards[i] = bytes.NewBuffer(buf)
	}

	// Split into files.
	data := make([]io.Writer, dataShards)
	for i := range data {
		data[i] = BufferShards[i]
	}
	// Do the split
	err = enc.Split(buffer, data, int64(buffer.Len()))
	if err != nil {
		return
	}

	input := make([]io.Reader, dataShards)
	for i := range input {
		input[i] = BufferShards[i]
	}

	// Create parity output writers
	parity := make([]io.Writer, parityShards)
	for i := range parity {
		parity[i] = BufferShards[dataShards+i]
	}

	// Encode parity
	err = enc.Encode(input, parity)
	if err != nil {
		return
	}
	return
}

// //将文件编码成多个数据块(未对文件数据块进行多数传输成功检测)
// func DistributeFileToNodes(UUID, filePath, fileKey string) error {
// 	totalShards := dataShards + parityShards
// 	outdir := config.FileShardDir
// 	nodes := config.FileNodeUrls
// 	// Encode the file into shards
// 	fileShardsPaths, err := encodeFile(filePath, outdir, dataShards, parityShards)
// 	if err != nil {
// 		return err
// 	}
// 	uuid := UUID
// 	var wg sync.WaitGroup
// 	//var mu sync.Mutex

// 	// Send shards to nodes
// 	errs := make([]error, totalShards)
// 	for i, filePartsPath := range fileShardsPaths {
// 		wg.Add(1)
// 		go func(i int, filePartsPath string) {
// 			var err error
// 			defer wg.Done()
// 			nodeIndex := i % len(nodes)
// 			nodeURL := nodes[nodeIndex]
// 			if nodeURL == config.LocalFileNodeUrl {
// 				err = nil
// 			} else {
// 				err = sendCopyFileToNode(nodeURL, filePartsPath, "shard")
// 			}
// 			if err != nil {
// 				errs[i] = err
// 			} else {
// 				errs[i] = nil
// 			}
// 		}(i, filePartsPath)
// 	}
// 	wg.Wait()

// 	//----------测试阶段------暂时将errs信息输出,后续需要将errs写入指定日志中
// 	for _, err := range errs {
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	}

// 	// Save metadata as JSON
// 	fileInfo, err := os.Stat(filePath)
// 	if err != nil {
// 		return err
// 	}
// 	fileSize := fileInfo.Size()

// 	fileMeta := FileMetadata{
// 		UUID:         uuid,
// 		EncodingTime: time.Now(),
// 		FileKey:      fileKey,
// 		FileSize:     fileSize,
// 		DataShards:   dataShards,
// 		ParityShards: parityShards,
// 		Shards:       make([]Shard, totalShards),
// 		// ShardIDs:       make([]string, totalShards),
// 		// ShardNodeNames: make([]string, totalShards),
// 	}

// 	for i := 0; i < totalShards; i++ {
// 		fileMeta.Shards[i].ShardID = i
// 		fileMeta.Shards[i].NodeURL = nodes[i%len(nodes)]
// 	}
// 	err = saveFileMeta(fileMeta)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

//根据UUID获取文件meta信息
// func GetFileMetaByUUID(UUID string) (fileMeta *FileMetadata, err error) {
// 	if UUID == "" {
// 		return nil, errors.New("UUID cannot be empty")
// 	}

// 	meta, err := redis.GetFileMeta(UUID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	metaBytes := []byte(meta)

// 	fileMeta = new(FileMetadata)
// 	err = json.Unmarshal(metaBytes, fileMeta)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return fileMeta, nil
// }

// func createDir(dir string) (err error) {
// 	if _, err := os.Stat(dir); os.IsNotExist(err) {
// 		// 如果目录不存在，则创建
// 		err = os.MkdirAll(dir, 0755) // 0755 表示具有读/写/执行权限的所有者，以及读/执行权限的其他用户。
// 		if err != nil {
// 			fmt.Println("创建目录失败:", err)
// 			return err
// 		}
// 	}
// 	return
// }

// Function to decode the shards and reconstruct the file
// func reconstructFile(shardPaths []string, outputPath string, dataShardNum, parityShardNum int, fileSize int64) (err error) {
// 	// Create matrix
// 	enc, err := reedsolomon.NewStream(dataShardNum, parityShardNum)
// 	if err != nil {
// 		return
// 	}

// 	// Create shards and load the data.
// 	shards := make([]io.Reader, dataShardNum+parityShardNum)
// 	for i, sPath := range shardPaths {
// 		f, err := os.Open(sPath)
// 		if err != nil {
// 			shards[i] = nil
// 			continue
// 		} else {
// 			shards[i] = f
// 		}
// 	}

// 	// Verify the shards
// 	ok, err := enc.Verify(shards)
// 	if err != nil {
// 		return
// 	}
// 	//close input file
// 	for i := range shards {
// 		if shards[i] != nil {
// 			err = shards[i].(*os.File).Close()
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	if !ok { // Verification failed. Reconstructing data
// 		for i, sPath := range shardPaths {
// 			f, err := os.Open(sPath)
// 			if err != nil {
// 				shards[i] = nil
// 				continue
// 			} else {
// 				shards[i] = f
// 			}
// 		}
// 		// Create out destination writers
// 		out := make([]io.Writer, len(shards))
// 		for i := range out {
// 			if shards[i] == nil {
// 				outfn := shardPaths[i]
// 				out[i], err = os.Create(outfn)
// 				if err != nil {
// 					return err
// 				}
// 			}
// 		}
// 		err = enc.Reconstruct(shards, out)
// 		if err != nil {
// 			return err
// 		}
// 		// Close output.
// 		for i := range out { //---------------测试阶段--------只输出最后遇到的错误
// 			if out[i] != nil {
// 				err = out[i].(*os.File).Close()
// 				if err != nil {
// 					return err
// 				}
// 			}
// 		}

// 		//close input file
// 		for i := range shards {
// 			if shards[i] != nil {
// 				err = shards[i].(*os.File).Close()
// 				if err != nil {
// 					return err
// 				}
// 			}
// 		}

// 		for i, sPath := range shardPaths {
// 			f, err := os.Open(sPath)
// 			if err != nil {
// 				shards[i] = nil
// 				continue
// 			} else {
// 				shards[i] = f
// 			}
// 		}
// 		ok, err = enc.Verify(shards)
// 		if !ok {
// 			return errors.New("Verification failed after reconstruction, data likely corrupted")
// 		}
// 		if err != nil {
// 			return
// 		}

// 		//close input file
// 		for i := range shards {
// 			if shards[i] != nil {
// 				err = shards[i].(*os.File).Close()
// 				if err != nil {
// 					return err
// 				}
// 			}
// 		}
// 	}

// 	f, err := os.Create(outputPath)
// 	if err != nil {
// 		return
// 	}
// 	defer f.Close()

// 	for i, sPath := range shardPaths {
// 		f, err := os.Open(sPath)
// 		if err != nil {
// 			shards[i] = nil
// 			continue
// 		} else {
// 			shards[i] = f
// 		}
// 		defer f.Close()
// 	}

// 	// We don't know the exact filesize.
// 	err = enc.Join(f, shards, fileSize)
// 	if err != nil {
// 		return
// 	}
// 	return
// }

//如果出现错误的shard在则在相应的shardBuffer设置为nil
func ReconstructBuffer(shardsBuffer []*bytes.Buffer, outputPath string, dataShardNum, parityShardNum int, BlockSize int64) (err error) {
	// Create matrix
	enc, err := reedsolomon.NewStream(dataShardNum, parityShardNum)
	if err != nil {
		return
	}

	// Create shards and load the data.
	shards := make([]io.Reader, dataShardNum+parityShardNum)
	originalBytes := make([][]byte, dataShardNum+parityShardNum)
	for i, buf := range shardsBuffer {
		shards[i] = shardsBuffer[i]
		originalBytes[i] = buf.Bytes()
	}

	// Verify the shards
	ok, err := enc.Verify(shards)
	if err != nil {
		return
	}
	if !ok { // Verification failed. Reconstructing data
		for i := range shardsBuffer {
			if shardsBuffer[i] != nil {
				shardsBuffer[i].Reset()
				shardsBuffer[i].Write(originalBytes[i])
				shards[i] = shardsBuffer[i]
			}
		}
		// Create out destination writers
		out := make([]io.Writer, len(shards))
		for i := range out {
			if shards[i] == nil {
				shardsBuffer[i] = bytes.NewBuffer(make([]byte, EC_ShardSize))
				out[i] = shardsBuffer[i]
			}
		}
		err = enc.Reconstruct(shards, out)
		if err != nil {
			return err
		}

		//把恢复成功的数据写回[]byte保存
		for i := range originalBytes {
			if originalBytes[i] == nil {
				originalBytes[i] = shardsBuffer[i].Bytes()
				shards[i] = shardsBuffer[i]
			} else {
				shardsBuffer[i].Reset()
				shardsBuffer[i].Write(originalBytes[i])
				shards[i] = shardsBuffer[i]
			}
		}

		ok, err = enc.Verify(shards)
		if !ok {
			return ErrReconstruct
		}
		if err != nil {
			return
		}
	}

	f, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	for i := range originalBytes {
		shards[i] = bytes.NewBuffer(originalBytes[i])
	}

	// We don't know the exact filesize.
	err = enc.Join(f, shards, BlockSize)
	if err != nil {
		return
	}
	return
}
