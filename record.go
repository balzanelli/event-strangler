package eventstrangler

import "time"

type RecordState string

const (
	ProcessingRecordState RecordState = "processing"
	CompleteRecordState   RecordState = "complete"
)

type Record struct {
	HashKey   string
	Status    RecordState
	CreatedAt time.Time
	ExpiresAt time.Time
}

func (r *Record) complete() {
	r.Status = CompleteRecordState
}

func (r *Record) isProcessing() bool {
	return r.Status == ProcessingRecordState
}

func (r *Record) isExpired() bool {
	return r.ExpiresAt.After(time.Now().UTC())
}

func getRecordExpirationTime(currentTime time.Time, timeToLive *time.Duration) time.Time {
	var duration time.Duration
	if timeToLive != nil {
		duration = time.Second * *timeToLive
	} else {
		duration = time.Hour * 1
	}
	return currentTime.Add(duration)
}

func newRecord(hashKey string, timeToLive *time.Duration) *Record {
	currentTime := time.Now().UTC()

	return &Record{
		HashKey:   hashKey,
		Status:    ProcessingRecordState,
		CreatedAt: currentTime,
		ExpiresAt: getRecordExpirationTime(currentTime, timeToLive),
	}
}
