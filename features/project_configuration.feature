Feature: View and change project configuration
	As a Developer
	I want to view and change project settings
	I can customise them to my projects needs

	Background:
		Given I have an empty directory

	Scenario: Attempt to view settings in non-git dir
		When I run "config list"
		Then I should see an error message informing me "initialise repository first"

	Scenario: Initialise configuration on first run
		Given I have initialised git
		And I have set my user details
		When I run "config list"
		Then I should see some configuration keys and values

	Scenario: Get all default configuration values
		Given I have initialised git
		And I have set my user details
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
		Given I have initialised git
		And I have set my user details
		When I run "config get testkey"
		Then I should see an error message informing me "no config key 'testkey' found"

	Scenario: Get a single configuration value
		Given I have initialised git
		And I have set my user details
		When I run "config get project.remote"
		Then I should see the following:
		"""
origin
		"""

	Scenario: Set config value with invalid format
		Given I have initialised git
		And I have set my user details
		When I run "config set testvalue"
		Then I should see an error message informing me "invalid argument format, expected: key=value"

	Scenario: Set a non-existing configuration key
		Given I have initialised git
		And I have set my user details
		When I run "config set testkey=testvalue"
		Then I should see an error message informing me "no config key 'testkey' found"

	Scenario: Set a config value
		Given I have initialised git
		And I have set my user details
		When I run "config set project.name=TestProject"
		Then The config key "project.name" should equal "TestProject"
