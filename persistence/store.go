package persistence

func NewStore(configStore ConfigStorer, metadataStorer MetadataStorer) *Store {
	return &Store{
		ConfigStorer:   configStore,
		MetadataStorer: metadataStorer,
	}
}

type Store struct {
	ConfigStorer   ConfigStorer
	MetadataStorer MetadataStorer
}
