Feature: User redirection to malicious domains
  As a defender
  I want to restrict redirection only to approved list of domains
  So that the attacker cannot redirect users to a domain their control
  @redirection-uri 
  Scenario: Require clients to register full redirection URI and partial versions of these URIs must trigger a error response.