Feature: Manage custom metadata

  Scenario: Git not initialised in valid spec setup
    Given I have a project directory
    But I have not initialised git
    And I have a file called "features/story1.feature" with the following content:
      """
      Feature: Story1
      """
    When I run "metadata add --story story1 key value"
    Then I should see an error message informing me "Please initialise repository first before running"

  Scenario: Attempt to add metadata to non-existent story
    Given I have a configured project directory
    When I run "metadata add --story doesnotexist key value"
    Then I should see an error message informing me "no story matching doesnotexist"

  Scenario: Successfully add metadata to extant story
    Given I have a configured project directory
    And I have a file called "features/story1.feature" with the following content:
      """
      Feature: Story1
      """
    When I run "metadata add --story story1 key value"
    Then The metadata "key" should be added to story "story1" with the value "value"
    And I should see no errors
