Feature: Lengthy attacks and defenses
  Scenario: Incident Response
    Given Second Pre-condition
    And Third Pre-condition
    When Second result is written to file on the target
    Then File monitor notifies defender of changes
    And Defender starts investigation
  Scenario: First Defense
    Given First Pre-condition
    And Third Pre-condition
    When Attacker performs first action
    Then Defender takes first step towards stopping the attack
  Scenario: Second Defense
    Given Second Pre-condition
    And First Defense
    When Attacker performs first action
    And Attacker performs second action
    Then Defender takes another step towards stopping the attack