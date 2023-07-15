Feature: Friends fight
  As a defender
  I want to make sure my friends are honest
  So that I can avoid being duped
  Background: Adam and Bob are friends
    Given Adam and Bob
    And some trust exists between Adam and Bob
  Scenario: Adam's cheating is caught
    When Adam cheats Bob
    Then Bob confronts Adam
  Scenario: Bob tries to verify Adam's story
    When Adam cooks-up a story to convince Bob
    Then Bob asks probing questions to Adam
    """markdown
		Bob's first question goes here
		Bob's second question goes here
    """
  Scenario: Bob has been confrontational in the past