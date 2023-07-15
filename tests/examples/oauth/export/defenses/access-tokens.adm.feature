Feature: Extract access tokens
  As a defender
  I want to secure access tokens
  So that attacker cannot extract access tokens
  @wide-scope-token 
  Scenario: Limit token lifetime
    When a token is issued to a client
    Then restrict its lifetime to a fix number of max. uses
    And limit max. use count to a sensible number.
  @wide-scope-token 
  Scenario: Don't issue token-clones
    Given a token with specific rights is issued to a client
    When that client requests a new token with same rights
    Then invalidate previous token
    And issue a new token
  @token-store 
  Scenario: Keep tokens only in transient memory
  @token-store 
  Scenario: Keep tokens only in private memory
  @token-store 
  Scenario: Apply same level of protection as refresh tokens to access tokens
  @wide-scope-token 
  Scenario: limit token scope
    When a token is issued to a client
    Then its access scope is limited to the minimum required rights for that client
  @wide-scope-token 
  Scenario: Configure token to expire after fixed period of no-use
    When a token is issued to a client
    And if client hasn't used it for a predefined period
    Then revoke that token