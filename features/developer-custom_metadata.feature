Feature: Manage custom metadata
  As a Developer
  I want to add abitrary metadata to my specification
  So I can cleanly capture contextual information about my specification in version control

  Scenario: Git not initialised in valid spec setup
    Given I have a project directory
    But I have not initialised git
    And I have a story called "story1"
    When I run "metadata add --story story1 key=value"
    Then I should see an error message informing me "Please initialise repository first before running"

  Scenario: Attempt to add metadata to non-existent story
    Given I have a configured project directory
    And the pushing mode is not set to automatic
    When I run "metadata add --story xxx key=value"
    Then I should see an error message informing me "no story matching xxx"

  Scenario: Try to identify ambiguous story
    Given I have a configured project directory
    And I have a story called "story1"
    And I have a story called "story2"
    When I run "metadata add --story story key=value"
    Then I should see an error message informing me "story name is ambiguous. Did you mean 'story1' or 'story2'?"

  Scenario: Successfully add metadata to extant story
    Given I have a configured project directory
    And I have a story called "story1"
    When I run "metadata add --story story1 key=value"
    Then The metadata "key" should be added to story "story1" with the value "value"
    And I should see no errors

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

  Scenario: Attempt to add metadata to non-existent scenario
    Given I have a configured project directory
    And the pushing mode is not set to automatic
    When I run "metadata add --scenario xxx key=value"
    Then I should see an error message informing me "no scenario matching xxx"

  Scenario: Try to identify ambiguous story
    Given I have a configured project directory
    And I have a file called "features/story1.feature" with the following content:
      """
      Feature: Story1
      	    Scenario: Scenario1
      		    Given some setup
      			When I do something
      			Then something happens
      
      	    Scenario: Scenario2
      			When I do something else
      			Then something else happens
      """
    When I run "metadata add --scenario scenario key=value"
    Then I should see an error message informing me "scenario ID is ambiguous. Did you mean 'scenario1' or 'scenario2'?"

  Scenario: Successfully add metadata to extant scenario by name
    Given I have a configured project directory
    And I have a file called "features/story1.feature" with the following content:
      """
      Feature: Story1
      	    Scenario: Scenario1
      		    Given some setup
      			When I do something
      			Then something happens
      """
    When I run "metadata add --scenario scenario1 key=value"
    Then The metadata "key" should be added to scenario "scenario1" with the value "value"
    And I should see no errors

  Scenario: Successfully add metadata to extant scenario by deterministic ID
    Given I have a configured project directory
    And I have a file called "features/story1.feature" with the following content:
      """
      Feature: Story1
      	    Scenario: Scenario1
      		    Given some setup
      			When I do something
      			Then Something happens
      """
    When I run "metadata add --scenario <FIXME:DID> key=value"
    Then The metadata "key" should be added to scenario "DID" with the value "value"
    And I should see no errors

  Scenario: Show metadata attached to a scenario
    Given I have a configured project directory
    And I have a story called "story1"
    And My story "story1" has a scenario called "scenario1" with the following metadata:
      | Name                 | Value                                                                                       |
      | Metadata one         | Value one                                                                                   |
      | Metadata two         | Neque porro quisquam est qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit |
      | Alphabetically first | This value should be at the top                                                             |
    When I run "metadata list --scenario scenario1"
    Then I should see the following:
      """
      Alphabetically first: This value should be at the top
      Metadata one        : Value one
      Metadata two        : Neque porro quisquam est qui dolorem ipsum quia dolor sit amet, consectetur...
      """
