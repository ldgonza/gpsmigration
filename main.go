package main

import (
	"database/sql"

	"github.com/magiconair/properties"
	"gitlab.com/simpliroute/gps-migration-to-json/db"
	"gitlab.com/simpliroute/gps-migration-to-json/work"
)

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func worker(i int, conn *sql.DB, p *properties.Properties, done chan bool) {
	work.Work(i, conn, p)
	done <- true
}

func main() {
	p := properties.MustLoadFile("application.properties", properties.UTF8)

	var (
		batchCount  = int(p.MustGetUint("operation.batch.count"))
		concurrency = int(p.MustGetUint("operation.concurrency"))
		readSource  = p.MustGetBool("operation.source.read")
	)

	var conn *sql.DB = nil
	if readSource {
		conn = db.Connect()
		defer db.Close(conn)
	}

	done := make(chan bool, 1)

	currentBatch := 0
	workerCount := min(batchCount, concurrency)

	// Init workers
	for i := 0; i < workerCount; i++ {
		currentBatch = i
		go worker(currentBatch, conn, p, done)
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
			go worker(currentBatch, conn, p, done)
		}
	}
}
