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

	/* 	for _, driver := range drivers {
		fmt.Println(driver)
	} */
	/*
		fmt.Println("-------------")
		fmt.Println("vehicles")
		vehicles := db.QueryVehicles("limit 2", conn)
		for _, vehicle := range vehicles {
			fmt.Println(vehicle)
		}

		fmt.Println("-------------")
		fmt.Println("drivers daily")
		driversDaily := db.QueryDriverDailyStatus("limit 2", conn)
		for _, driverDaily := range driversDaily {
			fmt.Println(driverDaily)
		}

		fmt.Println("-------------")
		fmt.Println("vehicles daily")
		vehiclesDaily := db.QueryVehicleDailyStatus("limit 2", conn)
		for _, vehicleDaily := range vehiclesDaily {
			fmt.Println(vehicleDaily)
		} */
}
