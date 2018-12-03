package metadata

func PrepareSync(sp SyncPreparer) error {
	return sp.PrepareMetadataSync()
}

func Pull(puller Puller, from string) error {
	return puller.PullMetadata(from)
}

func Push(pusher Pusher, to string) error {
	return pusher.PushMetadata(to)
}
