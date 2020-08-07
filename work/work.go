package work

import (
	"fmt"

	"gitlab.com/simpliroute/gps-migration-to-json/db"
)

// Work processes the migration.
func Work(i int) {
	conn := db.Connect()
	defer db.Close(conn)

	fmt.Println("-------------")
	fmt.Println("vehicles")
	vehicles := db.QueryVehicles("select * from tracking_location limit 2", conn)
	for _, vehicle := range vehicles {
		fmt.Println(vehicle)
	}

	fmt.Println("-------------")
	fmt.Println("drivers")
	drivers := db.QueryDrivers("select * from tracking_locationdriver limit 2", conn)
	for _, driver := range drivers {
		fmt.Println(driver)
	}
}
