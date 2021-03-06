package db

import (
	"database/sql"

	// Enable postgres driver
	_ "github.com/lib/pq"
)

// TrackingLocationVehicle row
type TrackingLocationVehicle struct {
	Created   string
	Modified  string
	ID        string
	Timestamp string
	Latitude  float64
	Longitude float64
	VehicleID int
	Alert     sql.NullString
}

// QueryVehicles executes a query for vehicles and returns the results
func QueryVehicles(limit string, conn *sql.DB) []TrackingLocationVehicle {
	baseQuery := "select l.created, l.modified, l.id, l.timestamp AT TIME ZONE 'UTC', l.latitude, l.longitude, l.vehicle_id, l.alert "
	baseQuery += "from "
	baseQuery += "tracking_location_temp l "

	// ./query.sh "select * from tracking_location_temp where row_number > 60000000 order by row_number asc limit 300000"
	// query := baseQuery + " order by l.timestamp asc  " + limit
	query := baseQuery + " where " + limit

	rows, err := conn.Query(query)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var results []TrackingLocationVehicle

	for rows.Next() {
		var row TrackingLocationVehicle
		err = rows.Scan(&row.Created, &row.Modified, &row.ID, &row.Timestamp, &row.Latitude, &row.Longitude, &row.VehicleID, &row.Alert)
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
