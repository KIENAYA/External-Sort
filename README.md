# ⚡ External Sort in Go

Efficiently sort huge files that don’t fit in memory using **external sorting**. This Go implementation splits input data into sorted chunks, then performs a **k-way merge** to produce a fully sorted output.

## 📦 Prerequisites

- Go >= 1.18

## ⚙️ Deploy

Generate 1 billion elements to run

 ```bash
go run .\generate-data\generate-data.go
```

Run
 ```bash
go run .\cmd\app\main.go
```




