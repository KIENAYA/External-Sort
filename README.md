# ⚡ External Sort in Go

Efficiently sort huge files that don’t fit in memory using **external sorting**. This Go implementation splits input data into sorted chunks, then performs a **k-way merge** to produce a fully sorted output.

## 📦 Prerequisites

- Go >= 1.18

## ⚙️ Deploy

Generate 1 billion elements to run (change numbers of elements in generate-data\generate-data.go)

 ```bash
go run .\generate-data\generate-data.go
```

Run
 ```bash
go run .\cmd\app\main.go
```

## 🔨 Config (Optional)

Config parameter: chunk size, input file name, output file name, numbers of line read in each chunk for merge, temp dir to store chunk

Default:

```dotenv
INPUT_FILE=input.txt
OUTPUT_FILE=output.txt
TMP_DIR=tmp
CHUNK_SIZE=10000000
MERGE_BATCH_SIZE=50000




