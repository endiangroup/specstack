Feature: Persist metadata
  As a Developer
  I want my metadata to persist as my specification changes
  So that my specification is always up to date with its contextual information

  Scenario: Pivot scenario metadata with file changes on git events
    Given I have a properly configured project directory
    And the pushing mode is not set to automatic
    And My story "story1" has a scenario called "scenario1" with some metadata
    And I make minor changes to scenario "scenario1" in "story1"
    When I commit and push my changes with git
    Then the metadata on "scenario1" should still exist

	## Auto mode
  Scenario: Pivot scenario metadata with file changes on spec execution
    Given I have a properly configured project directory
    And I have a story called "story1"
    And My story "story1" has a scenario called "scenario1" with some metadata
    And I make minor changes to scenario "scenario1" in "story1"
    When I run any spec command
    Then the metadata on "scenario1" should still exist
