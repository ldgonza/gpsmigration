package db

import (
	"database/sql"
	"fmt"

	// Enable postgres driver
	_ "github.com/lib/pq"

	"github.com/magiconair/properties"
)

func getProperty(p *properties.Properties, name string) string {
	val, ok := p.Get(name)

	if !ok {
		panic(fmt.Sprintf("Missing property! %s", name))
	}

	return val
}

// Connect returns a DB connection
func Connect() *sql.DB {
	p := properties.MustLoadFile("connection.properties", properties.UTF8)

	var (
		host     = getProperty(p, "db.host")
		port     = getProperty(p, "db.port")
		user     = getProperty(p, "db.user")
		password = getProperty(p, "db.pass")
		dbname   = getProperty(p, "db.name")
	)

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	return db
}

// Close closes a connection
func Close(conn *sql.DB) {
	conn.Close()
}

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

// QueryVehicles executes a query and returns the results
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
