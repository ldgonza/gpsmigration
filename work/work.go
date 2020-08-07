package work

import (
	"fmt"

	"gitlab.com/simpliroute/gps-migration-to-json/db"
)

// Work processes the migration.
func Work(i int) {
	conn := db.Connect()
	defer db.Close(conn)
	vehicles := db.QueryVehicles("select * from tracking_location limit 2", conn)

	for _, vehicle := range vehicles {
		fmt.Println(vehicle)
	}
}
