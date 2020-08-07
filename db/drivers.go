package db

import (
	"database/sql"

	// Enable postgres driver
	_ "github.com/lib/pq"
)

// TrackingLocationDriver row
type TrackingLocationDriver struct {
	created      string
	modified     string
	uuid         string
	timestamp    string
	latitude     float32
	longitude    float32
	driverID     int
	accuracy     sql.NullFloat64
	activityType sql.NullString
}

// QueryDrivers executes a query for drivers and returns the results
func QueryDrivers(query string, conn *sql.DB) []TrackingLocationDriver {
	rows, err := conn.Query(query)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var results []TrackingLocationDriver

	for rows.Next() {
		var row TrackingLocationDriver
		err = rows.Scan(&row.created, &row.modified, &row.uuid, &row.timestamp, &row.latitude, &row.longitude, &row.driverID, &row.accuracy, &row.activityType)
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
