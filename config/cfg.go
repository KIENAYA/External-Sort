package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	InputFile      string
	OutputFile     string
	TmpDir         string
	ChunkSize      int
	MergeBatchSize int
}

func LoadEnv(path string) (*Config, error) {
	err := godotenv.Load(path)
	if err != nil {
		log.Printf("Cannot find .env file: %v", err)
	}

	getInt := func(key string, defaultVal int) int {
		if val := os.Getenv(key); val != "" {
			if n, err := strconv.Atoi(val); err == nil {
				return n
			}
		}
		return defaultVal
	}

	cfg := &Config{
		InputFile:      os.Getenv("INPUT_FILE"),
		OutputFile:     os.Getenv("OUTPUT_FILE"),
		TmpDir:         os.Getenv("TMP_DIR"),
		ChunkSize:      getInt("CHUNK_SIZE", 1_000_000),
		MergeBatchSize: getInt("MERGE_BATCH_SIZE", 100),
	}

	return cfg, nil
}
