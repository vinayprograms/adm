# ADM Style Guide

Following are some style recommendations for `Given`, `When` or `Then` statements. These improve readability of statements connecting attacks and defenses.

1. Don't start the statement with pronouns (*I*, *We*, etc.). Addressing someone is part of model header. Attack and defense specifications must be person agnostic. In the example below, while addressing oneself helps in attack, it affects the readability of the connected defense.

    ```text
    # The identity of 'I' changes between attack and defense

    Attack: Gain root shell on target
      Given target runs OpenSSH server
      And uses default root credentials supplied by the distro
      When I connect to SSH on the target
      And I enter default root credentials
      Then I get a new shell with root access

    Defense: Block root shells on target
      When I connect to SSH on the target
      Then target detects a root login
      And blocks it
    ```

1. Don't start the statement with a verb (i.e., action word). For example, the statement `Then block source code from release` in a defense, when used in an attack (to break that defense) will read as `When block source code from release`. In comparison, `Then source code is blocked from release` and `When source code is blocked from release` look more natural from a grammatical standpoint.
