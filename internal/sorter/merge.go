package sorter

import (
	"bufio"
	"container/heap"
	"os"
	"strconv"
	"strings"
)

// ItemBatch đại diện cho batch đang xử lý từ một file
type ItemBatch struct {
	value  int           // Giá trị hiện tại
	index  int           // Index hiện tại trong batch
	batch  []int         // Batch dữ liệu
	file   *os.File      // File đang đọc
	reader *bufio.Reader // Reader đang sử dụng
}

// MinHeap là heap tối thiểu theo giá trị
type MinHeap []*ItemBatch

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i].value < h[j].value }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x any) {
	*h = append(*h, x.(*ItemBatch))
}

func (h *MinHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

// Đọc batch n dòng từ reader, trả về slice []int
func readBatchInt(reader *bufio.Reader, _ []int, n int) ([]int, error) {
	lines := make([]int, 0, n)
	for i := 0; i < n; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			if len(line) == 0 {
				return lines, err // EOF hoặc lỗi
			}
			line = strings.TrimSpace(line)
			if num, parseErr := strconv.Atoi(line); parseErr == nil {
				lines = append(lines, num)
			}
			return lines, err
		}
		line = strings.TrimSpace(line)
		if num, parseErr := strconv.Atoi(line); parseErr == nil {
			lines = append(lines, num)
		}
	}
	return lines, nil
}

// MergeChunksWithBatch thực hiện merge k-way từ nhiều file tạm thành 1 file output
func MergeChunksWithBatch(files []string, output string, batchSize int) error {
	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()

	writer := bufio.NewWriterSize(out, 16*1024*1024) // buffer 16MB
	defer writer.Flush()

	minHeap := &MinHeap{}
	heap.Init(minHeap)

	// Mỗi file tạo batch riêng
	for _, path := range files {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		reader := bufio.NewReaderSize(f, 16*1024*1024)

		batch, err := readBatchInt(reader, nil, batchSize)
		if err != nil && len(batch) == 0 {
			f.Close()
			continue
		}

		if len(batch) == 0 {
			f.Close()
			continue
		}

		heap.Push(minHeap, &ItemBatch{
			value:  batch[0],
			index:  0,
			batch:  batch,
			file:   f,
			reader: reader,
		})
	}

	// Sử dụng strings.Builder để gom nội dung trước khi ghi
	var sb strings.Builder
	sb.Grow(12 * batchSize) // sơ bộ cấp phát

	for minHeap.Len() > 0 {
		item := heap.Pop(minHeap).(*ItemBatch)

		sb.WriteString(strconv.Itoa(item.value))
		sb.WriteByte('\n')

		if sb.Len() > 64*1024 {
			writer.WriteString(sb.String())
			sb.Reset()
		}

		item.index++
		if item.index < len(item.batch) {
			item.value = item.batch[item.index]
			heap.Push(minHeap, item)
		} else {
			newBatch, err := readBatchInt(item.reader, nil, batchSize)
			if err != nil && len(newBatch) == 0 {
				item.file.Close()
				continue
			}
			if len(newBatch) == 0 {
				item.file.Close()
				continue
			}
			item.batch = newBatch
			item.index = 0
			item.value = newBatch[0]
			heap.Push(minHeap, item)
		}
	}

	// Ghi phần còn lại
	if sb.Len() > 0 {
		writer.WriteString(sb.String())
	}

	return nil
}
