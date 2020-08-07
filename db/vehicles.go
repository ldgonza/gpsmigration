package db

import (
	"database/sql"

	// Enable postgres driver
	_ "github.com/lib/pq"
)

// TrackingLocationVehicle row
type TrackingLocationVehicle struct {
	created   string
	modified  string
	id        string
	timestamp string
	latitude  float32
	longitude float32
	vehicleID int
	alert     sql.NullString
}

// QueryVehicles executes a query for vehicles and returns the results
func QueryVehicles(query string, conn *sql.DB) []TrackingLocationVehicle {
	rows, err := conn.Query(query)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var results []TrackingLocationVehicle

	for rows.Next() {
		var row TrackingLocationVehicle
		err = rows.Scan(&row.created, &row.modified, &row.id, &row.timestamp, &row.latitude, &row.longitude, &row.vehicleID, &row.alert)
		if err != nil {
			// handle this error
			panic(err)
		}

		results = append(results, row)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return results
}
