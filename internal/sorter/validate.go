package sorter

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// ValidateSortedFile kiểm tra file có được sort đúng thứ tự số học hay không.
func ValidateSortedFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("không thể mở file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var prev int
	lineNum := 0

	if scanner.Scan() {
		first := scanner.Text()
		prev, err = strconv.Atoi(first)
		if err != nil {
			return fmt.Errorf("lỗi parse dòng đầu: %v", err)
		}
		lineNum = 1
	}

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		current, err := strconv.Atoi(line)
		if err != nil {
			return fmt.Errorf("lỗi parse dòng %d: %v", lineNum, err)
		}
		if current < prev {
			return fmt.Errorf("❌ Lỗi tại dòng %d: %d < %d", lineNum, current, prev)
		}
		prev = current
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("lỗi đọc file: %v", err)
	}

	fmt.Printf("✅ File '%s' đã được sort đúng thứ tự số học.\n", path)
	return nil
}
