package persistence

func NewStore(configStore ConfigStorer) *Store {
	return &Store{
		ConfigStorer: configStore,
	}
}

type Store struct {
	ConfigStorer ConfigStorer
}
