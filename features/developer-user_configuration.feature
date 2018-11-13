Feature: View and change user configuration

  Scenario: No git user name set
    Given I have an empty directory
    And I have initialised git
    When I run "config list"
    Then I should see an error message informing me "no user.name set"

  Scenario: No git user email set
    Given I have an empty directory
    And I have initialised git
    And I have set the git user name to "Spec Stack"
    When I run "config list"
    Then I should see an error message informing me "no user.email set"

  Scenario: Get all default configuration values
    Given I have an empty directory
    And I have initialised git
    And I have set the git user name to "Spec Stack"
    And I have set the git user email to "dev@specstack.io"
    When I run "config list"
    Then I should see the following:
      """
      user.name=Spec Stack
      user.email=dev@specstack.io
      """
