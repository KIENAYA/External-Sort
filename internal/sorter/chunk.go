package sorter 

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"time"
)

func SortChunksParallel(inputFile, tmpDir string, chunkSize int) ([]string, error) {
	in, err := os.Open(inputFile)
	if err != nil {
		return nil, err
	}
	defer in.Close()

	scanner := bufio.NewScanner(in)

	chunkChan := make(chan []string, runtime.NumCPU())
	resultChan := make(chan string, 1024)
	var wg sync.WaitGroup

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for chunk := range chunkChan {
				filePath, err := sortAndWriteChunkParallel(chunk, tmpDir, workerID)
				if err != nil {
					panic(err)
				}
				resultChan <- filePath
			}
		}(i)
	}

	chunk := []string{}
	chunkID := 0
	for scanner.Scan() {
		chunk = append(chunk, scanner.Text())
		if len(chunk) >= chunkSize {
			copyChunk := make([]string, len(chunk))
			copy(copyChunk, chunk)
			chunkChan <- copyChunk
			chunk = chunk[:0]
			chunkID++
		}
	}
	if len(chunk) > 0 {
		copyChunk := make([]string, len(chunk))
		copy(copyChunk, chunk)
		chunkChan <- copyChunk
	}
	close(chunkChan)

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var tmpFiles []string
	for path := range resultChan {
		tmpFiles = append(tmpFiles, path)
	}

	return tmpFiles, nil
}

func sortAndWriteChunkParallel(chunk []string, tmpDir string, id int) (string, error) {
	nums := make([]int, 0, len(chunk))
	for _, line := range chunk {
		n, err := strconv.Atoi(line)
		if err != nil {
			return "", fmt.Errorf("không thể parse '%s': %w", line, err)
		}
		nums = append(nums, n)
	}

	sort.Ints(nums)

	fileName := fmt.Sprintf("chunk_%d_%d.txt", id, time.Now().UnixNano())
	path := filepath.Join(tmpDir, fileName)
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	for _, n := range nums {
		fmt.Fprintln(writer, n)
	}
	writer.Flush()

	return path, nil
}
