package db

import (
	"database/sql"

	// Enable postgres Vehicle
	_ "github.com/lib/pq"
)

// TrackingVehicleDailyStatus row
type TrackingVehicleDailyStatus struct {
	ID           int
	Created      string
	Modified     string
	Date         string
	ProviderName sql.NullString
	LocationID   sql.NullString
	VehicleID    int
	Location     TrackingLocationVehicle
}

// QueryVehicleDailyStatus executes a query for Vehicles and returns the results
func QueryVehicleDailyStatus(where string, conn *sql.DB) []TrackingVehicleDailyStatus {
	baseQuery := "select d.id "
	baseQuery += ", d.created, d.modified, d.date, d.provider_name, d.location_id, d.vehicle_id "
	baseQuery += ", l.created, l.modified, l.id, l.timestamp AT TIME ZONE 'UTC', l.latitude, l.longitude, l.vehicle_id, l.alert "
	baseQuery += "from "
	baseQuery += "tracking_vehicledailystatus d "
	baseQuery += "inner join tracking_location l on l.id = d.location_id "

	query := baseQuery + " " + where

	rows, err := conn.Query(query)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var results []TrackingVehicleDailyStatus

	for rows.Next() {
		var row TrackingVehicleDailyStatus
		err = rows.Scan(
			&row.ID, &row.Created, &row.Modified, &row.Date, &row.ProviderName, &row.LocationID, &row.VehicleID,
			&row.Location.Created, &row.Location.Modified, &row.Location.ID, &row.Location.Timestamp, &row.Location.Latitude, &row.Location.Longitude, &row.Location.VehicleID, &row.Location.Alert)

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
