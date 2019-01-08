package specstack

/*
func Test_Initialise_ReturnsErrorIfRepositoryIsntInitialised(t *testing.T) {
	mockRepo := &repository.MockRepository{}
	mockDeveloper := &personas.MockDeveloper{}
	repoStore := persistence.NewStore(&persistence.MockConfigStorer{}, &persistence.MockMetadataStorer{})
	ctrl := New("", mockRepo, mockDeveloper, repoStore, os.Stdout, os.Stderr)

	mockRepo.On("IsInitialised").Return(false)

	assert.Equal(t, ErrUninitialisedRepo, app.Initialise())
}

func Test_Initialise_CreatesConfigOnFirstRun(t *testing.T) {
	mockRepo := &repository.MockRepository{}
	mockDeveloper := &personas.MockDeveloper{}
	mockConfigStore := &persistence.MockConfigStorer{}
	repoStore := persistence.NewStore(mockConfigStore, &persistence.MockMetadataStorer{})
	app := New("", mockRepo, mockDeveloper, repoStore, os.Stdout, os.Stderr)

	mockRepo.On("IsInitialised").Return(true)
	mockRepo.On("GetConfig", "user.name").Return("username", nil)
	mockRepo.On("GetConfig", "user.email").Return("user@email", nil)
	mockRepo.On("PrepareMetadataSync").Return(nil)
	mockConfigStore.On("AllConfig").Return(map[string]string{}, nil)

	assert.NoError(t, app.Initialise())

	mockConfigStore.AssertExpectations(t)
}

func Test_Initialise_ReturnsErrorWhenMissingUsername(t *testing.T) {
	mockRepo := &repository.MockRepository{}
	mockDeveloper := &personas.MockDeveloper{}
	mockConfigStore := &persistence.MockConfigStorer{}
	repoStore := persistence.NewStore(mockConfigStore, &persistence.MockMetadataStorer{})
	app := New("/testing/test-dir", mockRepo, mockDeveloper, repoStore, os.Stdout, os.Stderr)

	mockRepo.On("IsInitialised").Return(true)
	mockRepo.On("GetConfig", "user.name").Return("", persistence.ErrNoConfigFound)
	mockConfigStore.On("AllConfig").Return(map[string]string{}, persistence.ErrNoConfigFound)

	err := app.Initialise()

	assert.IsType(t, MissingRequiredConfigValueErr(""), err)
}

func Test_Initialise_ReturnsErrorWhenMissingEmail(t *testing.T) {
	mockRepo := &repository.MockRepository{}
	mockDeveloper := &personas.MockDeveloper{}
	mockConfigStore := &persistence.MockConfigStorer{}
	repoStore := persistence.NewStore(mockConfigStore, &persistence.MockMetadataStorer{})
	app := New("/testing/test-dir", mockRepo, mockDeveloper, repoStore, os.Stdout, os.Stderr)

	mockRepo.On("IsInitialised").Return(true)
	mockRepo.On("GetConfig", "user.name").Return("username", nil)
	mockRepo.On("GetConfig", "user.email").Return("", persistence.ErrNoConfigFound)
	mockConfigStore.On("AllConfig").Return(map[string]string{}, persistence.ErrNoConfigFound)

	err := app.Initialise()

	assert.IsType(t, MissingRequiredConfigValueErr(""), err)
}

func Test_Initialise_SetsConfigDefaults(t *testing.T) {
	mockRepo := &repository.MockRepository{}
	mockDeveloper := &personas.MockDeveloper{}
	mockConfigStore := &persistence.MockConfigStorer{}
	repoStore := persistence.NewStore(mockConfigStore, &persistence.MockMetadataStorer{})
	app := New("/testing/test-dir", mockRepo, mockDeveloper, repoStore, os.Stdout, os.Stderr).(*appController)
	expectedConfig := config.NewWithDefaults()
	expectedConfig.Project.Name = "test-dir"
	expectedConfig.User.Name = "username"
	expectedConfig.User.Email = "user@email"

	mockRepo.On("IsInitialised").Return(true)
	mockRepo.On("GetConfig", "user.name").Return("username", nil)
	mockRepo.On("GetConfig", "user.email").Return("user@email", nil)
	mockRepo.On("PrepareMetadataSync").Return(nil)
	mockConfigStore.On("SetConfig", mock.Anything, mock.Anything).Return(nil)
	mockConfigStore.On("AllConfig").Return(map[string]string{}, persistence.ErrNoConfigFound)
	mockConfigStore.On("StoreConfig", mock.Anything).Return(nil, nil)

	assert.NoError(t, app.Initialise())

	assert.Equal(t, expectedConfig, app.config)
}

func Test_Initialise_LoadsExistingConfigIfNotFirstRun(t *testing.T) {
	mockRepo := &repository.MockRepository{}
	mockDeveloper := &personas.MockDeveloper{}
	mockConfigStore := &persistence.MockConfigStorer{}
	repoStore := persistence.NewStore(mockConfigStore, &persistence.MockMetadataStorer{})
	app := New("/testing/test-dir", mockRepo, mockDeveloper, repoStore, os.Stdout, os.Stderr).(*appController)
	expectedConfig := config.NewWithDefaults()

	mockRepo.On("IsInitialised").Return(true)
	mockRepo.On("PrepareMetadataSync").Return(nil)
	mockDeveloper.On("ListConfiguration", mock.Anything).Return(nil, nil)
	mockConfigStore.On("AllConfig").Return(config.ToMap(expectedConfig), nil)

	assert.NoError(t, app.Initialise())

	mockConfigStore.AssertExpectations(t)
	assert.Equal(t, expectedConfig, app.config)
}
*/
