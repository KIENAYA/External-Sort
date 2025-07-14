package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	const total = 1_000_000_000
	const output = "input.txt"

	f, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	writer := bufio.NewWriterSize(f, 4*1024*1024) 
	defer writer.Flush()

	rand.Seed(time.Now().UnixNano())

	start := time.Now()
	for i := 0; i < total; i++ {
		n := rand.Intn(2_000_000_000) 
		fmt.Fprintln(writer, n)

		// if i > 0 && i%10_000_000 == 0 {
		// 	fmt.Printf(" Đã ghi %d dòng...\n", i)
		// }
	}
	fmt.Println("Completed in:", time.Since(start))
}
