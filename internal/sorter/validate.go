package sorter

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)


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
			return fmt.Errorf("parse error first line: %v", err)
		}
		lineNum = 1
	}

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		current, err := strconv.Atoi(line)
		if err != nil {
			return fmt.Errorf("Parse fail at line %d: %v", lineNum, err)
		}
		if current < prev {
			return fmt.Errorf("Error at %d: %d < %d", lineNum, current, prev)
		}
		prev = current
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Cannot read file: %v", err)
	}

	fmt.Printf("File '%s' is sorted.\n", path)
	return nil
}
