package work

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/magiconair/properties"

	"gitlab.com/simpliroute/gps-migration-to-json/db"
	"gitlab.com/simpliroute/gps-migration-to-json/output"
)

// Table what table to opreate on
type Table string

const (
	// Drivers export driver locations
	Drivers Table = "Drivers"
	// Vehicles export vehicle locations
	Vehicles Table = "Vehicles"
	// DriversDaily export daily driver status
	DriversDaily Table = "DriversDaily"
	// VehiclesDaily export daily vehicle status
	VehiclesDaily Table = "VehiclesDaily"
)

// StrToTable returns the Table enum for the string or panics
func StrToTable(table string) Table {
	if table == string(Drivers) {
		return Drivers
	}
	if table == string(Vehicles) {
		return Vehicles
	}
	if table == string(DriversDaily) {
		return DriversDaily
	}
	if table == string(VehiclesDaily) {
		return VehiclesDaily
	}

	panic("Unknown table!: " + table)
}

func createDir(dir string) {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(dir, 0755)
		if errDir != nil {
			log.Fatal(err)
		}
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func dolog(i int, message string) {
	ts := time.Now().UnixNano()
	fmt.Println(ts, "-", "Batch", i, "-", message)
}

// Work processes the migration and returns true if nothing else to do
func Work(i int, conn *sql.DB, uploader *s3manager.Uploader, p *properties.Properties) (done bool) {
	defer func() {
		if r := recover(); r != nil {
			dolog(i, fmt.Sprint("Fatal! ", r))
			done = true
		}
	}()

	dolog(i, "Starting")

	var (
		tableString = p.MustGetString("operation.table")
		bucketName  = p.MustGetString("bucket.name")
		outputType  = p.MustGetString("output.type")
		fileSize    = int(p.MustGetUint("output.file.size"))
		fileCount   = int(p.MustGetUint("output.files.per.batch"))
	)

	fileNumber := i * fileCount
	readSize := fileCount * fileSize

	table := StrToTable(tableString)

	dir := "results/" + string(table)
	createDir(dir)
	doneReading := true

	dolog(i, "Reading")
	offset := i * int(readSize)
	//where := " limit " + strconv.Itoa(readSize) + " offset " + strconv.Itoa(offset)
	where := "row_number > " + strconv.Itoa(offset) + " order by row_number asc limit " + strconv.Itoa(readSize)
	dolog(i, where)

	var locations []output.TrackingLocation
	var latestStatus []output.LatestTrackingStatus

	if table == Drivers {
		drivers := db.QueryDrivers(where, conn)
		locations = output.TransformDrivers(drivers)
		dolog(i, fmt.Sprintf("Length %d", len(locations)))
	} else if table == Vehicles {
		vehicles := db.QueryVehicles(where, conn)
		locations = output.TransformVehicles(vehicles)
		dolog(i, fmt.Sprintf("Length %d", len(locations)))
	} else if table == DriversDaily {
		dailyDrivers := db.QueryDriverDailyStatus(where, conn)
		latestStatus = output.TransformDailyDrivers(dailyDrivers)
		dolog(i, fmt.Sprintf("Length %d", len(latestStatus)))
	} else if table == VehiclesDaily {
		dailyVehicles := db.QueryVehicleDailyStatus(where, conn)
		latestStatus = output.TransformDailyVehicles(dailyVehicles)
		dolog(i, fmt.Sprintf("Length %d", len(latestStatus)))
	} else {
		panic("Unimplemented table!: " + string(table))
	}

	if len(latestStatus) > 0 {
		var buffer []output.LatestTrackingStatus
		filename := ""
		for curr, location := range latestStatus {
			if len(buffer) == 0 {
				filename = dir + "/" + strconv.Itoa(fileNumber) + ".json"
				dolog(i, fmt.Sprintf("Preparing JSON %s", filename))
			}

			buffer = append(buffer, location)

			if len(buffer) >= fileSize || curr >= len(latestStatus)-1 {
				if outputType == "file" {
					output.WriteLatestTrackingStatusToFile(filename, buffer)
				} else if outputType == "gcp" {
					output.WriteLatestTrackingStatusToCloudStorage(bucketName, filename, buffer)
				} else if outputType == "s3" {
					output.WriteLatestTrackingStatusToS3(uploader, bucketName, filename, buffer)
				} else if outputType == "console" {
					output.WriteLatestTrackingStatusToConsole(buffer)
				}

				buffer = nil
				fileNumber++
			}

		}
		doneReading = false
	}

	if len(locations) > 0 {
		var buffer []output.TrackingLocation
		filename := ""
		for curr, location := range locations {
			if len(buffer) == 0 {
				filename = dir + "/" + strconv.Itoa(fileNumber) + ".json"
				dolog(i, fmt.Sprintf("Preparing JSON %s", filename))
			}

			buffer = append(buffer, location)

			if len(buffer) >= fileSize || curr >= len(locations)-1 {
				if outputType == "file" {
					output.WriteLocationsToFile(filename, buffer)
				} else if outputType == "gcp" {
					output.WriteLocationsToCloudStorage(bucketName, filename, buffer)
				} else if outputType == "s3" {
					output.WriteLocationsToS3(uploader, bucketName, filename, buffer)
				} else if outputType == "console" {
					output.WriteLocationsToConsole(buffer)
				}

				buffer = nil
				fileNumber++
			}

		}
		doneReading = false
	}

	done = doneReading

	dolog(i, "Done: "+strconv.FormatBool(done))
	return done
}
