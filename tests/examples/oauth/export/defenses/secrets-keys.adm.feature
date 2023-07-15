Feature: Client secrets and keys
  As a defender
  I want to block access to client secrets
  Scenario: Don't publish app binary containing secrets
    When binary is scanned
    Then block release of app binary
  @yolosec 
  Scenario: Don't store secrets in config file
    Given application uses config files for secrets
    When config file is scanned
    And secrets are revealed
    Then block release of app binary
  @apache 
  Scenario: use Apache Web Server v2.5.51 or later
    Given web-server is Apache (version <= v2.5.50)
    Then upgrade Apache Web Server to v2.5.51 or later
  Scenario: Encrypt secrets before storing in file
  Scenario: Store secrets and keys in vault
  Scenario: Don't publish code containing secrets
    When repo is scanned
    And secrets are revealed
    Then release of code is blocked
  Scenario: Block encoded paths
    When request contains encoded URL of config file
    Then decode encoded URL
    And block URL
  Scenario: Don't allow path traversal to configuration files
    Given secrets are stored in configuration file
    When request contains URL of config file
    Then block the request
  Scenario: Web-server shouldn't serve config files
    Given secrets are stored in configuration file
    When request contains URL of config file
    Then web-server is configured to not serve such file types