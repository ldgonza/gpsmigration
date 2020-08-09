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

// Work processes the migration.
func Work(i int, conn *sql.DB, p *properties.Properties) {
	var (
		batchSize   = int(p.MustGetUint("operation.batch.size"))
		read        = p.MustGetBool("operation.source.read")
		write       = p.MustGetBool("operation.destination.write")
		delete      = p.MustGetBool("operation.file.delete")
		tableString = p.MustGetString("operation.table")
	)

	table := StrToTable(tableString)

	dir := "results/" + string(table)
	createDir(dir)
	filename := dir + "/" + strconv.Itoa(i) + ".json"

	if read {
		offset := i * int(batchSize)
		where := " limit " + strconv.Itoa(batchSize) + " offset " + strconv.Itoa(offset)

		if table == Drivers {
			drivers := db.QueryDrivers(where, conn)
			driverLocations := output.TransformDrivers(drivers)
			output.WriteLocationsToFile(filename, driverLocations)
		} else if table == Vehicles {
			vehicles := db.QueryVehicles(where, conn)
			vehicleLocations := output.TransformVehicles(vehicles)
			output.WriteLocationsToFile(filename, vehicleLocations)
		} else if table == DriversDaily {
			dailyDrivers := db.QueryDriverDailyStatus(where, conn)
			driversLatest := output.TransformDailyDrivers(dailyDrivers)
			output.WriteLatestTrackingStatusToFile(filename, driversLatest)
		} else if table == VehiclesDaily {
			dailyVehicles := db.QueryVehicleDailyStatus(where, conn)
			vehiclesLatest := output.TransformDailyVehicles(dailyVehicles)
			output.WriteLatestTrackingStatusToFile(filename, vehiclesLatest)
		} else {
			panic("Unimplemented table!: " + string(table))
		}
	}

	if write {
		fmt.Println("writing")

		if delete {
			fmt.Println("deleting file")
		}
	}
}
