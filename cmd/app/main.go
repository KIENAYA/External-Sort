package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"external-sort/config"
	"external-sort/internal/sorter"
)

func main() {
	cfg, err := config.LoadEnv(".env")
	if err != nil {
		log.Fatal(err)
	}

	os.MkdirAll(cfg.TmpDir, os.ModePerm)

	fmt.Println("Sorting chunks...")
	start := time.Now()
	chunks, err := sorter.SortChunksParallel(cfg.InputFile, cfg.TmpDir, cfg.ChunkSize)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sorted %d chunks in %v\n", len(chunks), time.Since(start))

	fmt.Println("Merging chunks...")
	start = time.Now()
	err = sorter.MergeChunksWithBatch(chunks, cfg.OutputFile, cfg.MergeBatchSize)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Merged into %s in %v\n", cfg.OutputFile, time.Since(start))

	fmt.Println("Validating output file...")
	start = time.Now()
	if err := sorter.ValidateSortedFile(cfg.OutputFile); err != nil {
		log.Fatalf("Validation failed: %v", err)
	} else {
		fmt.Printf("Validation passed in %v\n", time.Since(start))
	}

}
