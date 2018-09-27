Feature: Create new project
	As a Technical Manager
	I want to create a new project
	I can customise my settings

	Scenario: Attempt to init in non-git dir
		Given I have an empty directory
		When I run the "init" command
		Then I should see an error message informing me "initialise git first"

