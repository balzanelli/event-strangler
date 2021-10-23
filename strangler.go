package eventstrangler

type Config struct {
	InstanceName string
}

type Strangler struct {
	config *Config
	store  *Store
}

func (s *Strangler) Once() (bool, error) {
	return false, nil
}

func (s *Strangler) Complete() error {
	return nil
}

func NewStrangler(config *Config) (*Strangler, error) {
	store, err := NewStore()
	if err != nil {
		return nil, err
	}
	return &Strangler{
		config: config,
		store:  store,
	}, nil
}
