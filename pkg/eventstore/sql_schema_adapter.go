package eventstore

import (
	"fmt"
	"strings"
)

const (
	defaultEventsTableName = "events"
	defaultSelectQuery     = `
SELECT 
	aggregate_id, 
	aggregate_version, 
	event_name, 
	event_payload
FROM %s
WHERE aggregate_id = $1 AND aggregate_type = $2
ORDER BY aggregate_version ASC;
`
	defaultInsertQuery = `
INSERT INTO %s (
	aggregate_id, 
	aggregate_version, 
	aggregate_type, 
	event_name, 
	event_payload
)
VALUES %s
`
	defaultInsertMarkersCount  = 5
	defaultInsertMarkersPatter = "($%d,$%d,$%d,$%d,$%d),"
)

func defaultInsertMarkers(count int) string {
	result := strings.Builder{}

	var indices []any
	for i := 1; i <= count*defaultInsertMarkersCount; i++ {
		indices = append(indices, i)
		if i%defaultInsertMarkersCount == 0 {
			result.WriteString(
				fmt.Sprintf(
					defaultInsertMarkersPatter,
					indices...,
				),
			)
			indices = nil
		}
	}

	return strings.TrimRight(result.String(), ",")
}
