Feature: Security of refresh tokens
  As a defender
  I want to secure refresh tokens
  So that attacker cannot obtain access tokens
  Scenario: Refresh token security on web-servers
  Scenario: Authenticate client
    When a refresh token is sent from client
    Then authenticate the client
    And identify that client in the list of registered clients