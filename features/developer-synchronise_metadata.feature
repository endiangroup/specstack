Feature: Synchronise metadata
  As a Developer
  I want to synchronise my specification's metadata
  So that I can stay up to date with my team
  But I know that git has limitations, so I will accept three modes of operation
  for pushes and pulls, which can be set independently:
  1. Manual mode, where I can push and pull metadata explicity. I will use this
  when I want to push or pull metadata without making changes to my source code,
  or in unusual tech setups.
  2. Semi-automatic mode, where my metadata is pulled after a merge (generally a
  git pull) or pushed when I push my code changes to git. This is the way that
  I will generally fetch metadata.
  3. Automatic push mode, where my metadata are pushed to the git server as soon
  as I add them, without me having to do anything. I want this to be the default
  behaviour.

  Scenario: Git not initialised for manual pull
    Given I have a project directory
    But I have not initialised git
    When I run "pull"
    Then I should see an error message informing me "initialise repository first"

  Scenario: Git not initialised for manual push
    Given I have a project directory
    But I have not initialised git
    When I run "push"
    Then I should see an error message informing me "initialise repository first"

  Scenario: Git remote not configured for manual pull
    Given I have a git-initialised project directory
    But I have not configured a project remote
    When I run "pull"
    Then I should see an error message informing me "configure a project remote first"

  Scenario: Project remote not configured for manual push
    Given I have a git-initialised project directory
    But I have not configured a project remote
    When I run "push"
    Then I should see an error message informing me "configure a project remote first"

  Scenario: Project remote not configured automatic push
    Given I have a git-initialised project directory
    And I have set the pulling mode to automatic
    But I have not configured a project remote
    When I add some metadata
    Then I should see an error message informing me "configure a project remote first"

  Scenario: Project remote not configured for manual pull
    Given I have a git-initialised project directory
    But I have not configured a project remote
    When I run "pull"
    Then I should see an error message informing me "configure a project remote first"

  Scenario: Git remote not set for manual push
    Given I have a git-initialised project directory
    But I have not set a git remote
    When I run "push"
    Then I should see an error message informing me "set git remote 'origin' first"

  Scenario: Git remote not set for automatic push
    Given I have a git-initialised project directory
    And I have set the pushing mode to automatic
    But I have not set a git remote
    When I add some metadata
    Then I should see an error message informing me "set git remote 'origin' first"

  Scenario: Unexpected error for manual pull
    Given I have a properly configured project directory
    But The remote git server isn't responding properly
    When I run "pull"
    Then I should see an appropriate error from git

  Scenario: Unexpected error for manual push
    Given I have a properly configured project directory
    But The remote git server isn't responding properly
    And I have added some metadata
    When I run "push"
    Then I should see an appropriate error from git

  Scenario: Unexpected error for automatic push
    Given I have a properly configured project directory
    And I have set the pushing mode to automatic
    But The remote git server isn't responding properly
    When I add some metadata
    Then I should see an appropriate error from git

  @next
  Scenario: Successful manual pull
    Given I have a properly configured project directory
    And there are new metadata on the remote git server
    When I run "pull"
    Then my metadata should be fetched from the remote git server
    And I should see no errors

  Scenario: Successful manual push
    Given I have a properly configured project directory
    And the pushing mode is not set to automatic
    And I have added some metadata
    When I run "push"
    Then my metadata should be pushed to the remote git server
    And I should see no errors

  @next
  Scenario: Successful semi-automatic pull
    Given I have a properly configured project directory
    And I have set the pulling mode to semi-automatic
    And there are new metadata on the remote git server
    When I run a git pull
    Then my metadata should be fetched from the remote git server
    And I should see no errors

  Scenario: Successful semi-automatic push
    Given I have a properly configured project directory
    And I have set the pushing mode to semi-automatic
    And I have added some metadata
    When I run a git push
    Then my metadata should be pushed to the remote git server
    And I should see no errors

  Scenario: Successful automatic push
    Given I have a properly configured project directory
    And I have set the pushing mode to automatic
    When I add some metadata
    Then my metadata should be pushed to the remote git server
    And I should see no errors
