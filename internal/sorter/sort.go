package sorter

import (
	"fmt"
	"os"
)

func ExternalSort(input, output, tmpDir string, chunkSize int, batchSize int) error {
	err := os.MkdirAll(tmpDir, 0755)
	if err != nil {
		return fmt.Errorf("Cannot create temp dir: %w", err)
	}

	tmpFiles, err := SortChunksParallel(input, tmpDir, chunkSize)
	if err != nil {
		return fmt.Errorf("chunked fail: %w", err)
	}

	err = MergeChunksWithBatch(tmpFiles, output, batchSize)
	if err != nil {
		return fmt.Errorf("merge fail: %w", err)
	}
	return nil
}
