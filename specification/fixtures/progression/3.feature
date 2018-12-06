Feature: https://github.com/endiangroup/specstack/pull/44/commits/e64852c3ffdce8da620186a1a3b41dbe33bbdb7d

  Scenario: Successfully add metadata to extant story
    Given I have a configured project directory
    And I have a file called "features/story1.feature" with the following content:
      """
      				  Feature: Story1
      """
    When I run "metadata add --story story1 key=value"
    Then The metadata "key" should be added to story "story1" with the value "value"
    And I should see no errors

  Scenario: Successfully add metadata to extant story
    Given I have a configured project directory
    And I have a story called "story1"
    When I run "metadata add --story story1 key=value"
    Then The metadata "key" should be added to story "story1" with the value "value"
    And I should see no errors
