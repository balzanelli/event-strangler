package eventstrangler

// Config defines the user configurable options for event-strangler.
type Config struct {
	// HashKey contains the options used to generate the idempotency
	// key that is persisted into the DB.
	HashKey *HashKeyOptions
	// TimeToLive specifies how long idempotency keys remain in the DB in seconds
	// before they are deemed expired.
	//
	// Defaults to 1 hour.
	TimeToLive int
	// Store is an interface that is implemented for different DB providers,
	// currently examples include: LevelDB.
	Store Store
}

type Strangler struct {
	// config user defined configuration for event-strangler.
	config *Config
	// store implementation.
	store Store
}

// Once ensures that only one invocation of a method can be performed. The idempotency store
// is queried to ensure the idempotency key has not been previously complete.
//
// object should receive a json representation of the event.
func (s *Strangler) Once(object map[string]interface{}) (string, error) {
	hashKey, err := getHashKey(object, s.config.HashKey)
	if err != nil {
		return "", err
	}

	exists, err := s.store.Exists(hashKey)
	if err != nil {
		return hashKey, err
	}
	if exists {
		record, err := s.store.Get(hashKey)
		if err != nil {
			return hashKey, err
		}
		if record.isProcessing() {
			return hashKey, EventAlreadyProcessingError{
				HashKey:   hashKey,
				StartedAt: record.CreatedAt,
			}
		}
		return hashKey, EventAlreadyProcessedError{HashKey: hashKey}
	}

	record := newRecord(hashKey, s.config.TimeToLive)
	if err = s.store.Put(hashKey, record, s.config.TimeToLive); err != nil {
		return hashKey, err
	}
	return hashKey, nil
}

// Complete marks an idempotency key as completed, so it cannot be processed by
// subsequent invocations.
func (s *Strangler) Complete(hashKey string) error {
	record, err := s.store.Get(hashKey)
	if err != nil {
		return err
	}

	record.complete()
	if err = s.store.Put(hashKey, record, s.config.TimeToLive); err != nil {
		return err
	}
	return nil
}

// Purge expunges the idempotency key from the idempotency store. This allows the
// event to be processed by future invocations.
func (s *Strangler) Purge(hashKey string) error {
	return s.store.Delete(hashKey)
}

// New initializes a new strangler instance. `config` can be used to provide user
// defined configuration, such as configuring the DB storage backend and the
// idempotency key generation strategy.
func New(config *Config) (*Strangler, error) {
	if config == nil {
		return nil, InvalidConfigurationError{
			Message: "'Config' cannot be empty",
		}
	}
	if config.Store == nil {
		return nil, InvalidConfigurationError{
			Message: "'Store' cannot be empty",
		}
	}
	return &Strangler{
		config: config,
		store:  config.Store,
	}, nil
}
