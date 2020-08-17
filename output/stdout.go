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
	// Turn list to a an ID:LatestStatus map
	m := make(map[string]LatestTrackingStatus)
	for _, latest := range locations {
		m[latest.ID] = latest
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "    ")
	err := encoder.Encode(LatestTrackingStatusCollection{m})

	if err != nil {
		panic(err)
	}
}
