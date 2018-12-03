Feature: Manage custom metadata

  Scenario: Git not initialised in valid spec setup
    Given I have a project directory
    But I have not initialised git
    And I have a file called "features/story1.feature" with the following content:
      """
      Feature: Story1
      """
    When I run "metadata add --story story1 key=value"
    Then I should see an error message informing me "Please initialise repository first before running"

  Scenario: Attempt to add metadata to non-existent story
    Given I have a configured project directory
    And the pushing mode is not set to automatic
    When I run "metadata add --story xxx key=value"
    Then I should see an error message informing me "no story matching xxx"

  Scenario: Successfully add metadata to extant story
    Given I have a configured project directory
    And I have a file called "features/story1.feature" with the following content:
      """
      Feature: Story1
      """
    When I run "metadata add --story story1 key=value"
    Then The metadata "key" should be added to story "story1" with the value "value"
    And I should see no errors

  Scenario: Show metadata attached to a story
    Given I have a configured project directory
    And I have a file called "features/story1.feature" with the following content:
      """
      Feature: Story1
      """
    And My story "story1" has the following metadata:
      | Name                 | Value                                                                                       |
      | Metadata one         | Value one                                                                                   |
      | Metadata two         | Neque porro quisquam est qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit |
      | Alphabetically first | This value should be at the top                                                             |
    When I run "metadata list --story story1"
    Then I should see the following:
      """
      Alphabetically first: This value should be at the top
      Metadata one        : Value one
      Metadata two        : Neque porro quisquam est qui dolorem ipsum quia dolor sit amet, consectetur...
      """
