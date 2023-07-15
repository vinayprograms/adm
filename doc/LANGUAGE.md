# ADM Language Guide

ADM tool uses a modified version of [Gherkin](https://en.wikipedia.org/wiki/Cucumber_(software)#Gherkin_language). It adds additional rules and logic on top to help define attacks and defenses and to connect them into a decision graph.

## Model Structure

An attack-defense model is made of a model header followed by a set of assumption, attack, defense and policy blocks. Comments start with a `#` and continue till end of the line.

```text
Model: Model's title
  ...
  ...
  Assumption: Assumptions applicable to entire model
    Given ...
    And ....

  Attack: Attack Title
    Given ...
    When ...
    Then ...

  Defense: Defense Title
    Given ...
    When ...
    Then ...

  # Policies must always be specified at the end of the model
  Policy: Policy details
    Assumption: Assumptions specific to policy
      Given ...
      And ....
    Defense: A defense that implements this policy
      Given ...
      When ...
      Then ...
    Defense: Another defense that implements this policy
      Given ...
      When ...
      Then ...
```

### Model header

Model header starts with a model title statement followed by multiple blocks that capture the goals of attacker and defender. For each person, there are at-least two statements each indicating an intent and a goal -

* `As a...` statement identifies the role of the person to whom intents and goals belong.
* `I want to` captures the intent. Additional intents can specified by following this by `And` statements.
* `So that` captures the goal. Additional goals can specified by following this by `And` statements.

```text
Model: Client secret risks
  As a attacker
  I want to access client secrets
  So that I can masquerade as the original client

  As a defender
  I want to control access to client secrets
  So that it can only be used by legitimate clients

  # Rest of the model specification follows...
```

Example -

```text
Model: Risks with client secrets and keys
  As an attacker
  I want to gain root access on the target
  And access client secrets
  So that I can use the target to mine bitcoin
  And masquerade as the original client

  As a defender
  I want to control access to client secrets
  So that it can only be used for legitimate purposes
```

### Attack / Defense specification
Each attack / defense is made of 3 types of statements - `Given`, `When` and `Then`

* `Given` specifies the preconditions for that attack/defense.
* `When` is the step taken to execute the attack/defense.
* `Then` is the result expected out of that attack/defense.

Similar to model header, `And` statements can be used to specify additional `Given`, `When` and `Then` statements.

```text
Attack: Gain root shell on target
  Given target runs OpenSSH server
  And uses default root credentials supplied by the distro
  When SSH connection is established with target
  And default root credentials entered
  Then a new shell with root access is opened

Defense: Block root shells on target
  When OpenSSH server is deployed on the target
  Then change SSH configuration to disallow root logins
```

### Assumptions

A list of `Given` clauses can be added at the beginning of the model under a single `Assumption` heading. They represents all the assumptions made as part of the model. They will be applied to all attack and defense specifications in that model.

```text
Model: Security for Internet facing services
  Assumption: Minimum security infrastructure
		Given Internet facing service sits behind a firewall
		And Firewall blocks DoS attempts
  Attack: Code injection
    ...
```

### Policies

This construct can be used to capture security policies. Policies consist of two sections - Assumptions and Defenses. Assumptions, like explained in previous section, capture pre-conditions. This is followed by a list of defenses that capture security controls implemented as part of this policy. As explained earlier, defenses can be pre-emptive or incident response.

```text
Policy: Minimum Password requirements
  Assumption: Password based authentication
    Given users are identified using username and password
    And no other authentication mechanism is available
  Defense: Password must be at least 8 characters long
  Defense: Passwords must be changed upon first login
  Defense: Passwords must not contain common dictionary words
  Defense: Passwords must not be reused across multiple systems
  Defense: Passwords must not be stored using weak hashing algorithms
  Defense: Passwords must be salted before storing
```

Please note that policies must be specified at the end of the model or in a separate model.

### Partial specifications

As part of writing a model, you can start with writing title statements for assumptions, policies, attacks and defenses and skip the body (i.e., `Given`/`When`/`Then` statements). Such specifications are considered partial and will shown as empty boxes when generating a decision-graph diagram. Relations between attacks/defense won't be processed either.

Partial specifications can also be used as place-holder for other yet-to-be-specified models. This lets you iteratively build your models without breaking any relationships.

For example, you can start with a `containers.adm` file and put an empty attack to refer to another attack in kubernetes model -

```text
# File: containers.adm

Model: Secure application containers
    # External reference to an attack in "Kubernetes Security" model
    Attack: Attack on kubernetes infrastructure
```

Later, you can create `k8s.adm`, and link attacks

```text
# File: containers.adm

Model: Secure application containers
    # External reference changed to an attack chain
    Attack: Attack on kubernetes infrastructure
      Given kubernetes cluster uses docker as container runtime
      ...
```

```text
# File: k8s.adm

Model: Kubernetes Security
    Attack: Attack on Docker runtime
      Given kubernetes cluster uses docker as container runtime
      When ...
      Then ...
```

## Connecting attacks and defenses

Attacks and defenses can be connected in two ways - via statements or using tags.

### Connecting via statements

Attacks and defenses can be connected in various ways using same statement across attacks and defenses.

1. **Two attacks (or two defenses) can be chained** by using title of the first one in the `Given` statement of the next.

    ```text
    Defense: Root account is locked for SSH
      ...

    Defense: Monitor login attempts
      Given Root account is locked for SSH
      When someone attempts an SSH connection
      Then the time, user name and IP is logged
    ```

1. **A pre-emptive defense can be connected to an attack** by using the same clause in `When` statements of both the attack and defense. Pre-emptive defenses are those that are triggered as soon as an attack condition is met. Firewall & IDS rules, AppSec security controls, etc. are some examples of such defenses.

    ```text
    Attack: Scan ports on target
      Given target is exposed to internet
      When target is scanned for open TCP/UDP ports
      Then a list of standard ports are revealed

    Defense: Don't expose any TCP port apart from 80
      When target is scanned for open TCP/UDP ports
      Then all scans are blocked except port 80
    ```

1. **An incident-response can be connected to an attack** using the `Then` clause from an attack in the `When` clause of the defense. Incident-response is a defense that is triggered after an attack is successful. Password resets, various forms of alerting from security tools, etc. are some examples of such defenses.

    ```text
    Attack: Compromize user password containing dictionary words
      When passwords are leaked as part of data breach
      Then passwords are available on the internet for other attackers to use 

    Defense: Change compromized passwords
      When passwords are available on the internet for other attackers to use 
      Then all accounts using those passwords must change their password
    ```

1. **An attack can be connected to a defense** by using the `Then` clause from a defense in the `When` clause of an attack. This represents an attack that breaks a defense.

    ```text
    Defense: Use API Key to control client access
      When client tries to connect to the target
      Then API Key of that client must be passed in the request

    Attack: Steal API Key
      When API Key of that client must be passed in the request
      Then steal API Key of a client
      And use stolen API Key to connect to target
    ```

### Connecting via tags

Annotating an attack and a defense with the same tag connects them in the model. Tags allow for loose coupling between attacks and defenses. In this case, the tool automatically puts attack first followed by its defense, in the decision graph.

```text
@ssh-default-root
Attack: Gain root shell access on target
    ...

@ssh-default-root
Defense: Block root shells on target
    ...
```

The `@success` tag has special meaning. This can only be used with attacks. It indicates that an attack will succeed even if there are defenses to handle it. This can be used, for example, to represent social engineering attacks like phishing that may still succeed in the presence of defenses against credential compromize.

### Special Case - FOMOSec

When model contains a defense that doesn't mitigate any attack, the tool considers it `FOMOSec` i.e., a defense implemented out of fear, with no justifiable attack to mitigate. When generating the graph, such defendes are automatically tagged with `#fomosec` in the diagram.

### Additional Notes

The connection rules apply to the following too

1. Connecting attacks to defenses listed under a Policy
1. Chaining defenses between Model and Policy
1. Connecting pre-conditions between `Assumption` and `Attack`/`Defense` blocks.
1. Relationships across model boundaries when multiple model files are used to capture security information about a system.

## Additional Reading

Please go through the [style guide](STYLE-GUIDE.md) to learn how to write well formed ADM specifications.
