Feature: Synchronise metadata

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
    Then I should see a warning message informing me "configure a project remote first"

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
    Then I should see a warning message informing me "set git remote 'origin' first"

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
    Then I should see an appropriate warning from git

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
