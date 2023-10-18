package test

import (
	"libadm/loaders"
	"libadm/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

/////////////////////////////////////////////////
// Unit tests

func TestAttackFail(t *testing.T) {
	input := `
	Model: Model with just one rule
		Defense: Attack passed-off as defense
	`

	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var a model.Attack
	err = a.Init(gherkinModel.Feature.Children[0].Scenario)

	assert.Equal(t, err.Error(), "Expected 'Attack', got 'Defense'")
}

func TestAttack_TitleOnly(t *testing.T) {
	input := `
	Model: Model with just one rule
		Attack: All systems are vulnerable >:D
	`

	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var a model.Attack
	err = a.Init(gherkinModel.Feature.Children[0].Scenario)

	assert.Nil(t, err)
	assert.Equal(t, len(a.Tags), 0)
	assert.Equal(t, a.Title, "All systems are vulnerable >:D")
}

func TestAttack_ButStatementsNotSupported(t *testing.T) {
	input := `
	Model: Model with just one rule
		Attack: All systems are vulnerable
			But 'but' statement is not supported in ADM
	`

	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var a model.Attack
	err = a.Init(gherkinModel.Feature.Children[0].Scenario)

	assert.Equal(t, err.Error(), "unsupported keyword - 'But'")
}

func TestAttack_DuplicateStatements(t *testing.T) {
	testVectors := map[string][]string{
		"SamePrecondition": {`Model: Model with just one attack
		Attack: All systems are vulnerable
			Given system is exposed to internet
			Given system is exposed to internet
	`, "precondition - 'system is exposed to internet' is already part of this attack"},
		"SameAction": {`Model: Model with just one attack
		Attack: All systems are vulnerable
			When nmap scan is run
			When nmap scan is run
	`, "action - 'nmap scan is run' is already part of this attack"},
		"SameResult": {`Model: Model with just one attack
		Attack: All systems are vulnerable
			Then all possible ports are listed
			Then all possible ports are listed
	`, "result - 'all possible ports are listed' is already part of this attack"},
		"SamePreconditionAsAnd": {`Model: Model with just one attack
		Attack: All systems are vulnerable
			Given system is exposed to internet
			And system is exposed to internet
	`, "precondition - 'system is exposed to internet' is already part of this attack"},
		"SameActionAsAnd": {`Model: Model with just one attack
		Attack: All systems are vulnerable
			When nmap scan is run
			And nmap scan is run
	`, "action - 'nmap scan is run' is already part of this attack"},
		"SameResultAsAnd": {`Model: Model with just one attack
		Attack: All systems are vulnerable
			Then all possible ports are listed
			And all possible ports are listed
	`, "result - 'all possible ports are listed' is already part of this attack"},
	}

	for name, args := range testVectors {
		t.Run(name, func(t *testing.T) {
			input := args[0]
			expected := args[1]
			gherkinModel, err := loaders.LoadGherkinContent(input)
			if err != nil {
				t.Error(err)
			}

			var a model.Attack
			err = a.Init(gherkinModel.Feature.Children[0].Scenario)

			assert.Equal(t, expected, err.Error())
		})
	}
}

func TestAttack_FullBody(t *testing.T) {
	input := `
	Model: Model with just one rule
		Attack: All systems are vulnerable
			Given a system with un-detected vulnerabilities
			When an attack exploits one or more of those vulnerabilities
			Then exploited vulnerabilities are called zero-days
	`

	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var a model.Attack
	err = a.Init(gherkinModel.Feature.Children[0].Scenario)

	assert.Nil(t, err)
	assert.Equal(t, len(a.Tags), 0)
	assert.Equal(t, a.Title, "All systems are vulnerable")
	assert.NotNil(t, a.PreConditions["a system with un-detected vulnerabilities"])
	assert.NotNil(t, a.Actions["an attack exploits one or more of those vulnerabilities"])
	assert.NotNil(t, a.Results["exploited vulnerabilities are called zero-days"])
}

func TestAttack_Complete(t *testing.T) {
	input := `
	Model: Model with just one rule
		Attack: All systems are vulnerable
		Given a system with un-detected vulnerabilities
		And the system is exposed to internet
		When an attacker attacks this system over the internet
		And exploits one or more of the un-detected vulnerabilities
		Then exploited vulnerabilities are called zero-days
		And attacker will stay anonymous
	`

	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var a model.Attack
	err = a.Init(gherkinModel.Feature.Children[0].Scenario)

	assert.Nil(t, err)
	assert.Equal(t, len(a.Tags), 0)
	assert.Equal(t, a.Title, "All systems are vulnerable")
	assert.Equal(t, len(a.PreConditions), 2)
	assert.NotNil(t, a.PreConditions["a system with un-detected vulnerabilities"])
	assert.NotNil(t, a.PreConditions["the system is exposed to internet"])
	assert.Equal(t, len(a.Actions), 2)
	assert.NotNil(t, a.Actions["an attacker attacks this system over the internet"])
	assert.NotNil(t, a.Actions["exploits one or more of the un-detected vulnerabilities"])
	assert.Equal(t, len(a.Results), 2)
	assert.NotNil(t, a.Results["exploited vulnerabilities are called zero-days"])
	assert.NotNil(t, a.Results["attacker will stay anonymous"])
}
