package eventstore

import (
	"fmt"
	"strings"
)

const eventsTableName = "events"

func defaultInsertMarkers(count int) string {
	result := strings.Builder{}

	index := 1
	for i := 0; i < count; i++ {
		result.WriteString(fmt.Sprintf("($%d,$%d,$%d,$%d,$%d),", index, index+1, index+2, index+3, index+4))
		index += 5
	}

	return strings.TrimRight(result.String(), ",")
}
