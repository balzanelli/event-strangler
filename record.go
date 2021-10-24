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

func getExpirationDate(currentTime time.Time, timeToLive int) time.Time {
	return currentTime.Add(getTimeToLive(timeToLive))
}

func newRecord(hashKey string, timeToLive int) *Record {
	currentTime := time.Now().UTC()

	return &Record{
		HashKey:   hashKey,
		Status:    ProcessingRecordState,
		CreatedAt: currentTime,
		ExpiresAt: getExpirationDate(currentTime, timeToLive),
	}
}
