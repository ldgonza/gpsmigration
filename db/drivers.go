package db

import (
	"database/sql"

	// Enable postgres driver
	_ "github.com/lib/pq"
)

// TrackingLocationDriver row
type TrackingLocationDriver struct {
	Created      string
	Modified     string
	UUID         string
	Timestamp    string
	Latitude     float64
	Longitude    float64
	DriverID     int
	Accuracy     sql.NullFloat64
	ActivityType sql.NullString
}

// QueryDrivers executes a query for drivers and returns the results
func QueryDrivers(limit string, conn *sql.DB) []TrackingLocationDriver {
	baseQuery := "select "
	baseQuery += "l.created, l.modified, l.uuid, l.timestamp AT TIME ZONE 'UTC', l.latitude, l.longitude, l.driver_id, l.accuracy, l.activity_type "
	baseQuery += "from "
	baseQuery += "tracking_location_driver_temp l "

	// where := "row_number > " + strconv.Itoa(offset) + " order by row_number asc limit " + strconv.Itoa(readSize)
	// query := baseQuery + " order by l.timestamp asc " + limit
	query := baseQuery + " where " + limit

	rows, err := conn.Query(query)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var results []TrackingLocationDriver

	for rows.Next() {
		var row TrackingLocationDriver
		err = rows.Scan(&row.Created, &row.Modified, &row.UUID, &row.Timestamp, &row.Latitude, &row.Longitude, &row.DriverID, &row.Accuracy, &row.ActivityType)
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
