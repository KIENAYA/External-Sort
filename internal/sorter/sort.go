package sorter

import (
	"fmt"
	"os"
)

// ExternalSort là entrypoint gọi chia + sort + merge
func ExternalSort(input, output, tmpDir string, chunkSize int, batchSize int) error {
	err := os.MkdirAll(tmpDir, 0755)
	if err != nil {
		return fmt.Errorf("không tạo được thư mục tạm: %w", err)
	}

	tmpFiles, err := SortChunksParallel(input, tmpDir, chunkSize)
	if err != nil {
		return fmt.Errorf("lỗi chia chunk: %w", err)
	}

	err = MergeChunksWithBatch(tmpFiles, output, batchSize)
	if err != nil {
		return fmt.Errorf("lỗi merge: %w", err)
	}
	return nil
}
