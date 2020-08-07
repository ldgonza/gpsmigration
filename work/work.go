package work

import (
	"fmt"

	"gitlab.com/simpliroute/gps-migration-to-json/output"

	"gitlab.com/simpliroute/gps-migration-to-json/db"
)

// Work processes the migration.
func Work(i int) {
	conn := db.Connect()
	defer db.Close(conn)

	fmt.Println("-------------")
	fmt.Println("drivers")
	drivers := db.QueryDrivers("limit 2", conn)
	driverLocations := output.TransformDrivers(drivers)
	output.WriteLocationsToFile("output_drivers.json", driverLocations)

	fmt.Println("-------------")
	fmt.Println("vehicles")
	vehicles := db.QueryVehicles("limit 2", conn)
	vehicleLocations := output.TransformVehicles(vehicles)
	output.WriteLocationsToFile("output_vehicles.json", vehicleLocations)

	fmt.Println("-------------")
	fmt.Println("drivers daily")
	dailyDrivers := db.QueryDriverDailyStatus("limit 2", conn)
	driversLatest := output.TransformDailyDrivers(dailyDrivers)
	output.WriteLatestTrackingStatusToFile("output_daily_drivers.json", driversLatest)

	fmt.Println("-------------")
	fmt.Println("vehicles daily")
	dailyVehicles := db.QueryVehicleDailyStatus("limit 2", conn)
	vehiclesLatest := output.TransformDailyVehicles(dailyVehicles)
	output.WriteLatestTrackingStatusToFile("output_daily_vehicles.json", vehiclesLatest)
}
