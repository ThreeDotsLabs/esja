package eventstore

import (
	"fmt"
	"strings"
)

const (
	defaultEventsTableName = "events"
	defaultSelectQuery     = `
SELECT 
	stream_id, 
	stream_version, 
	stream_type, 
	event_name, 
	event_payload
FROM %s
WHERE stream_id = $1
ORDER BY stream_version ASC;
`
	defaultInsertQuery = `
INSERT INTO %s (
	stream_id, 
	stream_version, 
	stream_type, 
	event_name, 
	event_payload
)
VALUES %s
`
	defaultInsertMarkersCount   = 5
	defaultInsertMarkersPattern = "($%d,$%d,$%d,$%d,$%d),"
)

func defaultInsertMarkers(count int) string {
	result := strings.Builder{}

	var indices []any
	for i := 1; i <= count*defaultInsertMarkersCount; i++ {
		indices = append(indices, i)
		if i%defaultInsertMarkersCount == 0 {
			result.WriteString(
				fmt.Sprintf(
					defaultInsertMarkersPattern,
					indices...,
				),
			)
			indices = nil
		}
	}

	return strings.TrimRight(result.String(), ",")
}
