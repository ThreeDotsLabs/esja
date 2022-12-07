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
)

func defaultInsertMarkers(count int) string {
	result := strings.Builder{}

	index := 1
	for i := 0; i < count; i++ {
		result.WriteString(
			fmt.Sprintf(
				"($%d,$%d,$%d,$%d,$%d),",
				index,
				index+1,
				index+2,
				index+3,
				index+4,
			),
		)
		index += 5
	}

	return strings.TrimRight(result.String(), ",")
}
