Feature: End-user phishing
  As a defender
  I want to control the impact phishing attacks have
  So that attacker cannot access user resources
  @external-content 
  Scenario: Never load contents from external sources in authentication page
  @phishing 
  Scenario: Client should never ask for credentials
    When a client presents authentication page instead of redirecting to OAuth server
    Then don't enter credentials
  @phishing 
  Scenario: Verify sender's email
    When email asking to login again is received
    Then verify if sender's email is from the right domain address