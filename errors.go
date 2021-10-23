package eventstrangler

import (
	"fmt"
	"time"
)

type EventAlreadyProcessingError struct {
	HashKey   string
	StartedAt time.Time
}

func (e EventAlreadyProcessingError) Error() string {
	return fmt.Sprintf("Event With Hash Key '%s' Is Already Being Processed At '%s'", e.HashKey,
		e.StartedAt.String())
}

type EventAlreadyProcessedError struct {
	HashKey string
}

func (e EventAlreadyProcessedError) Error() string {
	return fmt.Sprintf("Event With Hash Key '%s' Has Already Been Processed", e.HashKey)
}
