package output

import (
	"strings"

	"gitlab.com/simpliroute/gps-migration-to-json/db"
)

func transformDriver(location db.TrackingLocationDriver) TrackingLocation {
	var accuracy *FloatValue = nil
	if location.Accuracy.Valid {
		accuracy = &FloatValue{Value: location.Accuracy.Float64}
	}

	var activityType *string = nil
	if location.ActivityType.Valid {
		upper := strings.ToUpper(location.ActivityType.String)
		activityType = &upper
	}

	timestamp := Timestamp{Seconds: 0, Nanos: 0}

	return TrackingLocation{
		ID:           nil,
		OwnerType:    "DRIVER",
		OwnerID:      location.DriverID,
		Accuracy:     accuracy,
		Alert:        nil,
		ProviderName: nil,
		BatteryLevel: nil,
		ActivityType: activityType,
		Timestamp:    timestamp,
		Location:     Location{Latitude: location.Latitude, Longitude: location.Longitude}}
}

func transformDailyDriver(latest db.TrackingDriverDailyStatus) LatestTrackingStatus {
	date := Date{Day: 0, Month: 0, Year: 0}

	result := LatestTrackingStatus{Date: date, Location: transformDriver(latest.Location)}
	
	var batteryLevel *FloatValue = nil
	if latest.BatteryLevel.Valid {
		batteryLevel = &FloatValue{Value: latest.BatteryLevel.Float64}
	}

	result.Location.BatteryLevel = batteryLevel

	return result
}

// TransformDrivers turns driver tracking status into TrackingLocations
func TransformDrivers(locations []db.TrackingLocationDriver) []TrackingLocation {
	var results []TrackingLocation
	for _, location := range locations {
		results = append(results, transformDriver(location))
	}

	return results
}

// TransformDailyDrivers turns driver latest  tracking status into LatestTrackingStatus
func TransformDailyDrivers(dailies []db.TrackingDriverDailyStatus) []LatestTrackingStatus {
	var results []LatestTrackingStatus
	for _, daily := range dailies {
		results = append(results, transformDailyDriver(daily))
	}

	return results
}
