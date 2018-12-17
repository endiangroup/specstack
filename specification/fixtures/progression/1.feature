Feature: https://github.com/endiangroup/specstack/pull/44/commits/e64852c3ffdce8da620186a1a3b41dbe33bbdb7d

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

  Scenario: Show metadata attached to a story
    Given I have a configured project directory
    And I have a story called "story1" in my spec with the following metadata:
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
