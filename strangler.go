package eventstrangler

type Config struct {
}

type Strangler struct {
	store *Store
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
	return &Strangler{store}, nil
}
