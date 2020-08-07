package output

import (
	"time"

	"github.com/golang/protobuf/ptypes"
)

func parseTimestamp(timestampStr string) Timestamp {
	time, err := time.Parse(time.RFC3339, timestampStr)

	if err != nil {
		panic(err)
	}

	timestamp, err := ptypes.TimestampProto(time)

	if err != nil {
		panic(err)
	}

	return Timestamp{Seconds: timestamp.Seconds, Nanos: timestamp.Nanos}
}
