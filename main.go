package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/magiconair/properties"
	"gitlab.com/simpliroute/gps-migration-to-json/db"
	"gitlab.com/simpliroute/gps-migration-to-json/output"
	"gitlab.com/simpliroute/gps-migration-to-json/work"
)

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func worker(firstBatch int, i int, conn *sql.DB, uploader *s3manager.Uploader, p *properties.Properties, done chan bool) {
	if i < firstBatch {
		done <- false
		return
	}

	done <- work.Work(i, conn, uploader, p)
}

func main() {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		fmt.Println("Total Time Elapsed: ", duration.Seconds(), "s")
	}()

	p := properties.MustLoadFile("application.properties", properties.UTF8)

	var (
		batchCount  = int(p.MustGetUint("operation.batch.count"))
		concurrency = int(p.MustGetUint("operation.concurrency"))
		firstBatch  = int(p.MustGetUint("operation.first.batch"))
		rampUpDelay = int(p.MustGetUint("operation.ramp.up.delay"))
		outputType  = p.MustGetString("output.type")
	)

	var uploader *s3manager.Uploader

	if outputType == "s3" {
		var err error
		uploader, err = output.GetS3Uploader()
		if err != nil {
			panic(err)
		}
	}

	var conn *sql.DB = nil
	conn = db.Connect()
	defer db.Close(conn)

	done := make(chan bool, 1)

	currentBatch := 0

	workerCount := concurrency

	if batchCount > 0 {
		workerCount = min(batchCount, concurrency)
	}

	// Init workers
	for i := 0; i < workerCount; i++ {
		currentBatch = i
		go worker(firstBatch, currentBatch, conn, uploader, p, done)
		time.Sleep(time.Duration(rampUpDelay) * time.Second)
	}

	// Restart when done until complete
	doneBatchCount := 0
	complete := false
	wrapUp := false
	result := false

	wrappedUpCount := 0
	for !complete {
		if wrapUp && wrappedUpCount >= workerCount {
			complete = true
			break
		}

		result = <-done
		wrapUp = result || wrapUp

		doneBatchCount++
		if wrapUp {
			wrappedUpCount++
		}

		if batchCount > 0 && doneBatchCount >= batchCount {
			complete = true
			break
		}

		if !wrapUp {
			currentBatch++
			if batchCount == 0 || currentBatch < batchCount {
				go worker(firstBatch, currentBatch, conn, uploader, p, done)
			}
		}
	}
}
