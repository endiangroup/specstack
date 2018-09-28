Feature: View and change project settings
	As a Developer
	I want to view and change settings
	I can customise them to my projects needs

	Scenario: Attempt to change settings in non-git dir
		Given I have an empty directory
		When I run "config list"
		Then I should see an error message informing me "initialise repository first"

	Scenario: View all config
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

