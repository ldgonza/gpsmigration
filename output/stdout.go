package output

import (
	"encoding/json"
	"os"
)

// WriteLocationsToConsole writes outputs to console
func WriteLocationsToConsole(locations []TrackingLocation) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "    ")
	err := encoder.Encode(TrackingLocationCollection{locations})

	if err != nil {
		panic(err)
	}
}

// WriteLatestTrackingStatusToConsole writes outputs to console
func WriteLatestTrackingStatusToConsole(locations []LatestTrackingStatus) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "    ")
	err := encoder.Encode(LatestTrackingStatusCollection{locations})

	if err != nil {
		panic(err)
	}
}
