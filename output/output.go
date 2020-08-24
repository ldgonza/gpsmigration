package output

import (
	"encoding/json"
	"os"
	"strconv"
)

// StringValue a string value
type StringValue struct {
	Value string `json:"value"`
}

// FloatValue a float value
type FloatValue struct {
	Value float64 `json:"value"`
}

// Date represents a date
type Date struct {
	Day   int `json:"day"`
	Month int `json:"month"`
	Year  int `json:"year"`
}

// Timestamp represents a unix epoch timestamp
type Timestamp struct {
	Seconds int64 `json:"seconds"`
	Nanos   int32 `json:"nanos"`
}

// LatestTrackingStatus represents last tracking status for the day
type LatestTrackingStatus struct {
	ID       string           `json:"-"`
	Date     Date             `json:"date"`
	Location TrackingLocation `json:"latest_location"`
}

// Location represents a GPS location
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// ActivityTypeNotSet null value for activity type
var ActivityTypeNotSet string = "ACTIVITY_TYPE_NOT_SET"

// TrackingLocation represents a tracking location
type TrackingLocation struct {
	ID           *int         `json:"id,omitempty"`
	OwnerType    string       `json:"owner_type"`
	OwnerID      int          `json:"owner_id"`
	Accuracy     *FloatValue  `json:"accuracy,omitempty"`
	Alert        *StringValue `json:"alert,omitempty"`
	ProviderName *StringValue `json:"provider_name,omitempty"`
	BatteryLevel *FloatValue  `json:"battery_level,omitempty"`
	ActivityType *string      `json:"activity_type,omitempty"`
	Timestamp    Timestamp    `json:"timestamp"`
	Location     Location     `json:"location"`
}

// TrackingLocationCollection represents the collection contents with the collection name
type TrackingLocationCollection struct {
	TrackingLocations []TrackingLocation `json:"tracking_locations"`
}

// LatestTrackingStatusCollection represents the collection contents with the collection name
type LatestTrackingStatusCollection struct {
	LatestTrackingStatus map[string]LatestTrackingStatus `json:"latest_tracking_status"`
}

func getDate(datestr string) Date {
	year, err := strconv.Atoi(datestr[0:4])
	if err != nil {
		panic(err)
	}

	month, err := strconv.Atoi(datestr[5:7])

	if err != nil {
		panic(err)
	}

	day, err := strconv.Atoi(datestr[8:10])

	if err != nil {
		panic(err)
	}

	return Date{Day: day, Month: month, Year: year}
}

// WriteLocationsToFile writes outputs to files
func WriteLocationsToFile(filename string, locations []TrackingLocation) {
	file, _ := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer file.Close()
	encoder := json.NewEncoder(file)

	err := encoder.Encode(TrackingLocationCollection{locations})
	if err != nil {
		panic(err)
	}
}

// WriteLatestTrackingStatusToFile writes outputs to files
func WriteLatestTrackingStatusToFile(filename string, locations []LatestTrackingStatus) {
	// Turn list to a an ID:LatestStatus map
	m := make(map[string]LatestTrackingStatus)
	for _, latest := range locations {
		m[latest.ID] = latest
	}

	file, _ := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer file.Close()
	encoder := json.NewEncoder(file)

	err := encoder.Encode(LatestTrackingStatusCollection{m})
	if err != nil {
		panic(err)
	}
}
