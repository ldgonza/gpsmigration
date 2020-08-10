package work

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/magiconair/properties"
	"gitlab.com/simpliroute/gps-migration-to-json/output"

	"gitlab.com/simpliroute/gps-migration-to-json/db"
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
	fmt.Println("Batch", i, "-", message)
}

// Work processes the migration and returns true if nothing else to do
func Work(i int, conn *sql.DB, p *properties.Properties) (done bool) {
	defer func() {
		if r := recover(); r != nil {
			dolog(i, fmt.Sprint("Fatal! ", r))
		}

		done = true
	}()

	dolog(i, "Starting")

	var (
		batchSize   = int(p.MustGetUint("operation.batch.size"))
		tableString = p.MustGetString("operation.table")
		bucketName  = p.MustGetString("bucket.name")
	)

	table := StrToTable(tableString)

	dir := "results/" + string(table)
	createDir(dir)

	filename := dir + "/" + strconv.Itoa(i) + ".json"
	doneReading := true

	dolog(i, "Reading")
	offset := i * int(batchSize)
	where := " limit " + strconv.Itoa(batchSize) + " offset " + strconv.Itoa(offset)

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
		dolog(i, fmt.Sprintf("Length %d", len(locations)))
	} else if table == VehiclesDaily {
		dailyVehicles := db.QueryVehicleDailyStatus(where, conn)
		latestStatus = output.TransformDailyVehicles(dailyVehicles)
		dolog(i, fmt.Sprintf("Length %d", len(locations)))
	} else {
		panic("Unimplemented table!: " + string(table))
	}

	if len(latestStatus) > 0 {
		fmt.Println("Preparing JSON " + filename)
		//output.WriteLatestTrackingStatusToFile(filename, latestStatus)
		output.WriteLatestTrackingStatusToCloudStorage(bucketName, filename, latestStatus)
		doneReading = false
	}

	if len(locations) > 0 {
		fmt.Println("Preparing JSON " + filename)
		// output.WriteLocationsToFile(filename, locations)
		output.WriteLocationsToCloudStorage(bucketName, filename, locations)
		doneReading = false
	}

	done = doneReading

	dolog(i, "Done: "+strconv.FormatBool(done))
	return done
}
