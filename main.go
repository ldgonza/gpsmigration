package main

import (
	"flag"

	"gitlab.com/simpliroute/gps-migration-to-json/work"
)

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func worker(i int, done chan bool) {
	work.Work(i)
	done <- true
}

func main() {
	var batchSize, batchCount, firstID, lastID, concurrency int

	flag.IntVar(&batchSize, "batchSize", 10, "batch size")
	flag.IntVar(&batchCount, "batchCount", 0, "total number of batches to run")
	flag.IntVar(&firstID, "firstId", 0, "first ID")
	flag.IntVar(&lastID, "lastId", 0, "last ID")
	flag.IntVar(&concurrency, "concurrency", 1, "number of parallel batches")

	flag.Parse()

	done := make(chan bool, 1)

	currentBatch := 0
	workerCount := min(batchCount, concurrency)

	// Init workers
	for i := 0; i < workerCount; i++ {
		currentBatch = i
		go worker(currentBatch, done)
	}

	// Restart when done until complete
	doneBatchCount := 0
	complete := false
	for !complete {
		<-done
		doneBatchCount++

		if batchCount > 0 && doneBatchCount == batchCount {
			complete = true
			break
		}

		currentBatch++
		if batchCount == 0 || currentBatch < batchCount {
			go worker(currentBatch, done)
		}
	}
}
