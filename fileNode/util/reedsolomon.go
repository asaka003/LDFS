package util

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/klauspost/reedsolomon"
)

var ( // 默认是3+3的数据块和验证快模式
	dataShards   int = 3
	parityShards int = 3

	tempFileShardsDir string = "/tmp/temp_FileShards"
	tempFileDir       string = "/tmp/temp_File"
)

func SetShardsMod(DataShards int, ParityShards int) {
	dataShards = DataShards
	parityShards = ParityShards
}

type FileMetadata struct {
	UUID         string    `json:"uuid"`
	FileKey      string    `json:"file_key"`
	FileHash     string    `json:"file_hash"`
	ContentType  string    `json:"content-type"`
	Etag         string    `json:"etag"`
	FileSize     int64     `json:"file_size"`
	DataShards   int       `json:"data_shards"`
	ParityShards int       `json:"parity_shards"`
	Shards       []Shard   `json:"shards"`
	EncodingTime time.Time `json:"encoding_time"`
}

type Shard struct {
	ShardID  int    `json:"shard_id"`
	NodeName string `json:"node_name"`
	NodeURL  string `json:"node_url"`
	Hash     string `json:"hash"`
}

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

func encodeFile(filePath, outDir string, dataShards, parityShards int) (FileShardsPath []string, err error) {
	// Create encoding matrix.
	enc, err := reedsolomon.NewStream(dataShards, parityShards)
	f, err := os.Open(filePath)
	if err != nil {
		return
	}

	instat, err := f.Stat()
	if err != nil {
		return
	}

	shards := dataShards + parityShards
	out := make([]*os.File, shards)
	FileShardsPath = make([]string, shards)

	// Create the resulting files.
	dir, file := filepath.Split(filePath)
	if outDir != "" {
		dir = outDir
	}
	for i := range out {
		outfn := fmt.Sprintf("%s.%d", file, i)
		FileShardsPath[i] = filepath.Join(dir, outfn)
		out[i], err = os.Create(FileShardsPath[i])
		if err != nil {
			return
		}
	}

	// Split into files.
	data := make([]io.Writer, dataShards)
	for i := range data {
		data[i] = out[i]
	}
	// Do the split
	err = enc.Split(f, data, instat.Size())
	if err != nil {
		return
	}

	// Close and re-open the files.
	input := make([]io.Reader, dataShards)

	for i := range data {
		out[i].Close()
		f, err = os.Open(out[i].Name())
		if err != nil {
			return
		}
		input[i] = f
		defer f.Close()
	}

	// Create parity output writers
	parity := make([]io.Writer, parityShards)
	for i := range parity {
		parity[i] = out[dataShards+i]
		defer out[dataShards+i].Close()
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

func createDir(dir string) (err error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// 如果目录不存在，则创建
		err = os.MkdirAll(dir, 0755) // 0755 表示具有读/写/执行权限的所有者，以及读/执行权限的其他用户。
		if err != nil {
			fmt.Println("创建目录失败:", err)
			return err
		}
	}
	return
}

// Function to decode the shards and reconstruct the file
// func ReconstructFile(UUID string) (outputPath string, err error) { // outputPath可以用作缓存处理
// 	// Load file metadata
// 	fileMeta, err := GetFileMetaByUUID(UUID)
// 	if err != nil {
// 		return
// 	}

// 	createDir(tempFileDir)
// 	createDir(tempFileShardsDir)

// 	var wg sync.WaitGroup
// 	// var mu sync.Mutex

// 	// Download shards
// 	shardsPaths := make([]string, fileMeta.DataShards+fileMeta.ParityShards)
// 	errs := make([]error, fileMeta.DataShards+fileMeta.ParityShards)
// 	for i, shard := range fileMeta.Shards {
// 		wg.Add(1)
// 		go func(i int, shard Shard) {
// 			defer wg.Done()
// 			uuidShardID := fmt.Sprintf("%s.%d", fileMeta.UUID, shard.ShardID)
// 			resp, err := http.Get(shard.NodeURL + "/getFileData/shard/" + uuidShardID)
// 			if err != nil {
// 				errs[i] = err
// 				return
// 			}
// 			defer resp.Body.Close()

// 			tempFilePath := path.Join(tempFileShardsDir, uuidShardID)
// 			tempFile, err := os.Create(tempFilePath)
// 			if err != nil {
// 				errs[i] = err
// 				return
// 			}
// 			_, err = io.Copy(tempFile, resp.Body)
// 			if err != nil {
// 				tempFile.Close()
// 				errs[i] = err
// 				return
// 			}
// 			tempFile.Close()
// 			shardsPaths[i] = tempFilePath
// 		}(i, shard)
// 	}
// 	wg.Wait()

// 	//----------测试阶段------暂时将errs信息输出,后续需要将errs写入指定日志中
// 	for _, err := range errs {
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	}

// 	// Create matrix
// 	enc, err := reedsolomon.NewStream(fileMeta.DataShards, fileMeta.ParityShards)
// 	if err != nil {
// 		return
// 	}

// 	// Open the inputs
// 	shards, _, err := openInput(fileMeta.DataShards, fileMeta.ParityShards, fileMeta.UUID, tempFileShardsDir)
// 	if err != nil {
// 		return
// 	}

// 	// Verify the shards
// 	ok, err := enc.Verify(shards)
// 	if ok {
// 		fmt.Println("No reconstruction needed")
// 	} else {
// 		fmt.Println("Verification failed. Reconstructing data", err)
// 		shards, _, err = openInput(fileMeta.DataShards, fileMeta.ParityShards, fileMeta.UUID, tempFileShardsDir)
// 		if err != nil {
// 			return
// 		}
// 		// Create out destination writers
// 		out := make([]io.Writer, len(shards))
// 		for i := range out {
// 			if shards[i] == nil { //---------------测试阶段---------这里不应该只在tmp目录下修复文件,还需要将文件修复到原节点中
// 				outfn := filepath.Join(tempFileShardsDir, fmt.Sprintf("%s.%d", fileMeta.UUID, i))
// 				out[i], err = os.Create(outfn)
// 				if err != nil {
// 					return
// 				}
// 			}
// 		}
// 		err = enc.Reconstruct(shards, out)
// 		if err != nil {
// 			fmt.Println("Reconstruct failed -", err)
// 			return
// 		}
// 		// Close output.
// 		for i := range out { //---------------测试阶段--------只输出最后遇到的错误
// 			if out[i] != nil {
// 				err = out[i].(*os.File).Close()
// 				// if err != nil {
// 				// 	return err
// 				// }
// 			}
// 		}
// 		if err != nil {
// 			return
// 		}

// 		shards, _, _ = openInput(fileMeta.DataShards, fileMeta.ParityShards, fileMeta.UUID, tempFileShardsDir)
// 		ok, err = enc.Verify(shards)
// 		if !ok {
// 			fmt.Println("Verification failed after reconstruction, data likely corrupted:", err)
// 			return
// 		}
// 		if err != nil {
// 			return
// 		}
// 	}

// 	// Join the shards and write them
// 	outputPath = filepath.Join(tempFileDir, fileMeta.UUID)
// 	f, err := os.Create(outputPath)
// 	if err != nil {
// 		return
// 	}

// 	shards, _, err = openInput(fileMeta.DataShards, fileMeta.ParityShards, fileMeta.UUID, tempFileShardsDir)
// 	if err != nil {
// 		return
// 	}

// 	// We don't know the exact filesize.
// 	err = enc.Join(f, shards, fileMeta.FileSize)
// 	if err != nil {
// 		return
// 	}

// 	// Cleanup temporary files
// 	// for _, shardPath := range shardsPaths {
// 	// 	err = os.Remove(shardPath)
// 	// 	if err != nil {
// 	// 		return err
// 	// 	}
// 	// }
// 	return
// }

func openInput(dataShards, parShards int, uuid, ShardsDir string) (r []io.Reader, size int64, err error) {
	// Create shards and load the data.
	shards := make([]io.Reader, dataShards+parShards)
	for i := range shards {
		infn := filepath.Join(ShardsDir, fmt.Sprintf("%s.%d", uuid, i))
		f, err := os.Open(infn)
		if err != nil {
			shards[i] = nil
			continue
		} else {
			shards[i] = f
		}
		stat, err := f.Stat()
		if err == nil && stat.Size() > 0 {
			size = stat.Size()
		} else {
			shards[i] = nil
		}
	}
	return shards, size, nil
}
