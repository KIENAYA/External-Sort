package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	const total = 1_000
	const output = "input1.txt"

	f, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	writer := bufio.NewWriterSize(f, 4*1024*1024) // 4MB buffer
	defer writer.Flush()

	rand.Seed(time.Now().UnixNano())

	start := time.Now()
	for i := 0; i < total; i++ {
		n := rand.Intn(2_000) // ví dụ giới hạn trong 2 tỷ
		fmt.Fprintln(writer, n)

		// if i > 0 && i%10_000_000 == 0 {
		// 	fmt.Printf("✅ Đã ghi %d dòng...\n", i)
		// }
	}
	fmt.Println("🎉 Hoàn tất! Thời gian:", time.Since(start))
}
