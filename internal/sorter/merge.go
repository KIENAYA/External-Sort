package sorter

import (
	"bufio"
	"container/heap"
	"os"
	"strconv"
	"strings"
)


type ItemBatch struct {
	value  int           
	index  int           
	batch  []int         
	file   *os.File      
	reader *bufio.Reader
}


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


func readBatchInt(reader *bufio.Reader, _ []int, n int) ([]int, error) {
	lines := make([]int, 0, n)
	for i := 0; i < n; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			if len(line) == 0 {
				return lines, err 
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

	var sb strings.Builder
	sb.Grow(12 * batchSize) 

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

	if sb.Len() > 0 {
		writer.WriteString(sb.String())
	}

	return nil
}
