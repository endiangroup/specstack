Feature: Persist metadata
  As a Developer
  I want my metadata to persist as my specification changes
  So that my specification is always up to date with its contextual information

  Scenario: Pivot scenario metadata with file changes on spec execution
    Given I have a properly configured project directory
    And My story "story1" has a scenario called "scenarioA" with some metadata
    And I make minor changes to scenario "scenario1" in "story1"
    When I run any spec metadata command
    Then the metadata on "scenario1" should still exist

  Scenario: Pivot scenario metadata with git commits
    Given I have a properly configured project directory
    And My story "story1" has a scenario called "scenario1" with some metadata
    And I make minor changes to scenario "scenario1" in "story1"
    When I make a commit
    Then the metadata on "scenario1" should still exist

  Scenario: Pivot scenario metadata with remote git changes
    Given I have a properly configured project directory
    And My story "story1" has a scenario called "scenario1" with some metadata
    When there are minor changes to scenario "scenario1" on the remote git server
    And I pull from the remote git server
    Then the metadata on "scenario1" should still exist
