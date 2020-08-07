package output

import (
	"gitlab.com/simpliroute/gps-migration-to-json/db"
)

func transformVehicle(location db.TrackingLocationVehicle) TrackingLocation {
	var alert *StringValue = nil
	if location.Alert.Valid {
		alert = &StringValue{Value: location.Alert.String}
	}

	timestamp := Timestamp{Seconds: 0, Nanos: 0}

	return TrackingLocation{
		ID:           nil,
		OwnerType:    "VEHICLE",
		OwnerID:      location.VehicleID,
		Accuracy:     nil,
		Alert:        alert,
		ProviderName: nil,
		BatteryLevel: nil,
		ActivityType: nil,
		Timestamp:    timestamp,
		Location:     Location{Latitude: location.Latitude, Longitude: location.Longitude}}
}

func transformDailyVehicle(latest db.TrackingVehicleDailyStatus) LatestTrackingStatus {
	date := Date{Day: 0, Month: 0, Year: 0}

	result := LatestTrackingStatus{Date: date, Location: transformVehicle(latest.Location)}

	var providerName *StringValue = nil
	if latest.ProviderName.Valid {
		providerName = &StringValue{Value: latest.ProviderName.String}
	}

	result.Location.ProviderName = providerName

	return result
}

// TransformVehicles turns driver tracking status into TrackingLocations
func TransformVehicles(locations []db.TrackingLocationVehicle) []TrackingLocation {
	var results []TrackingLocation
	for _, location := range locations {
		results = append(results, transformVehicle(location))
	}

	return results
}

// TransformDailyVehicles turns driver tracking status into LatestTrackingStatus
func TransformDailyVehicles(dailies []db.TrackingVehicleDailyStatus) []LatestTrackingStatus {
	var results []LatestTrackingStatus
	for _, daily := range dailies {
		results = append(results, transformDailyVehicle(daily))
	}

	return results
}
