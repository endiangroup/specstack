Feature: View and change project settings
	As a Developer
	I want to view and change settings
	I can customise them to my projects needs

	Scenario: Attempt to view settings in non-git dir
		Given I have an empty directory
		When I run "config list"
		Then I should see an error message informing me "initialise repository first"

	Scenario: Initialise configuration on first run
		Given I have an empty directory
		And I have initialised git
		When I run "config list"
		Then I should see some keys and values

	Scenario: Get all default configuration values
		Given I have an empty directory
		And I have initialised git
		When I run "config list"
		Then I should see the following:
		"""
project.remote=origin
project.name=test-dir
project.featuresdir=./features
project.pushingmode=auto
project.pullingmode=semi-auto
		"""

	Scenario: Attempt to get non-existing config key
		Given I have an empty directory
		And I have initialised git
		When I run "config get testkey"
		Then I should see an error message informing me "no config key 'testkey' found"

	Scenario: Get a single configuration value
		Given I have an empty directory
		And I have initialised git
		When I run "config get project.remote"
		Then I should see the following:
		"""
origin
		"""
