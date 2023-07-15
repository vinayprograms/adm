Feature: Add additional security controls around public clients
  As a defender
  I want to have additinal security around public client
  So that attacker cannot easily compromize public client
  Scenario: User consent for public clients
    Given client is public client		
    When client requests authorization code or refresh token
    Then warn resource owner that it is a public client
    And ask for additional consent from user
  Scenario: Deployment specific client secrets
    Given client is public client		
    When client requests authorization code or refresh token
    Then secrets unique to that instance of the client are issued
  Scenario: Client secrets revocation
    When client is identified as malicious
    Then all secrets associated with the client are revoked