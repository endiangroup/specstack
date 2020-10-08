Feature: https://github.com/endiangroup/specstack/pull/44/commits/e64852c3ffdce8da620186a1a3b41dbe33bbdb7d

  Scenario: Git not initialised in valid spec setup
    Given I have a project directory
    But I have not initialised git
    And I have a file called "features/story1.feature" with the following content:
      """
      Feature: Story1
      """
    When I run "metadata add --story story1 key=value"
    Then I should see an error message informing me "Please initialise repository first before running"

  Scenario: Git not initialised in valid spec setup
    Given I have a project directory
    But I have not initialised git
    And I have a story called "story1"
    When I run "metadata add --story story1 key=value"
    Then I should see an error message informing me "Please initialise repository first before running"
