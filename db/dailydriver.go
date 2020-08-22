package db

import (
	"database/sql"

	// Enable postgres driver
	_ "github.com/lib/pq"
)

// TrackingDriverDailyStatus row
type TrackingDriverDailyStatus struct {
	ID           int
	Created      string
	Modified     string
	Date         string
	BatteryLevel sql.NullFloat64
	AccountID    int
	DriverID     int
	LocationID   sql.NullString
	Location     TrackingLocationDriver
}

// QueryDriverDailyStatus executes a query for drivers and returns the results
func QueryDriverDailyStatus(limit string, conn *sql.DB) []TrackingDriverDailyStatus {
	baseQuery := "select d.id "
	baseQuery += ", d.created, d.modified, d.date, d.battery_level, d.account_id, d.driver_id, d.location_id "
	baseQuery += ", l.created, l.modified, l.uuid, l.timestamp AT TIME ZONE 'UTC', l.latitude, l.longitude, l.driver_id, l.accuracy, l.activity_type "
	baseQuery += "from "
	baseQuery += "tracking_driverdailystatus d "
	baseQuery += "inner join tracking_locationdriver l on l.uuid = d.location_id "

	query := baseQuery + " order by d.id asc " + limit

	rows, err := conn.Query(query)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var results []TrackingDriverDailyStatus

	for rows.Next() {
		var row TrackingDriverDailyStatus
		err = rows.Scan(
			&row.ID, &row.Created, &row.Modified, &row.Date, &row.BatteryLevel, &row.AccountID, &row.DriverID, &row.LocationID,
			&row.Location.Created, &row.Location.Modified, &row.Location.UUID, &row.Location.Timestamp, &row.Location.Latitude, &row.Location.Longitude, &row.Location.DriverID, &row.Location.Accuracy, &row.Location.ActivityType)

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
