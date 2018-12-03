package metadata

type SyncPreparer interface {
	PrepareMetadataSync() error
}

type Puller interface {
	PullMetadata(from string) error
}
type Pusher interface {
	PushMetadata(to string) error
}
