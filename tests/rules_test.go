package test

import (
	"libadm/loaders"
	"libadm/model"
	"testing"

	"github.com/cucumber/messages-go/v16"
	"github.com/stretchr/testify/assert"
)

func TestRule_AttackChain(t *testing.T) {
	input := `
	Model: Attack Chain
		Attack: Attack 1
		Attack: Attack 2
			Given Attack 1
	`
	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var model model.Model
	model.Init(gherkinModel.Feature)
	assert.Equal(t, "Attack 1",
		model.Attacks["Attack 2"].PreConditions["Attack 1"].Item[0].Title)
}

func TestRule_DefenseChain(t *testing.T) {
	input := `
	Model: Defense Chain
		Defense: Defense 1
		Defense: Defense 2
			Given Defense 1
	`
	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var model model.Model
	model.Init(gherkinModel.Feature)
	assert.Equal(t, "Defense 1",
		model.Defenses["Defense 2"].PreConditions["Defense 1"].Item[0].Title)
}

func TestRule_PolicyMitigatesAttack(t *testing.T) {
	f := messages.Feature{
		Language: "en",
		Keyword:  "Model",
		Name:     "Password security",
		Children: []*messages.FeatureChild{
			0: {
				Rule: &messages.Rule{
					Keyword: "Policy",
					Name:    "Password security best-practices",
					Children: []*messages.RuleChild{
						0: {
							Scenario: &messages.Scenario{
								Keyword: "Defense",
								Name:    "Password length must be greater than 8 characters",
								Tags: []*messages.Tag{
									0: {
										Name: "@password-len",
									},
								},
							},
						},
					},
				},
			},
			1: {
				Scenario: &messages.Scenario{
					Keyword: "Attack",
					Name:    "Bruteforce short passwords",
					Tags: []*messages.Tag{
						0: {
							Name: "@password-len",
						},
					},
				},
			},
		},
	}

	var model model.Model
	model.Init(&f)
	assert.Equal(t,
		"Bruteforce short passwords",
		model.Policies["Password security best-practices"].
			Defenses["Password length must be greater than 8 characters"].
			Tags["@password-len"][0].
			Title)
}

func TestRule_DefenseMitigatesAttack(t *testing.T) {
	input := `
	Model: Bruteforce passwords
		
		Attack: Run a script to attempt all possible passwords
			When dictionary passwords are sent to login API
			Then one of the passwords match

		Defense: Block repeated attempts to login
			When dictionary passwords are sent to login API
			And multiple passwords attempts fail in a short period of time
			Then the connection is closed
			And IP address of the requester is blocked for a predefined period of time
	`
	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var model model.Model
	model.Init(gherkinModel.Feature)
	assert.Equal(t, model.Attacks["Run a script to attempt all possible passwords"],
		model.
			Defenses["Block repeated attempts to login"].
			Actions["dictionary passwords are sent to login API"].Item[0])
}

func TestRule_AttackThwartsDefense(t *testing.T) {
	input := `
	Model: Log into SSH using commonly used accounts
		
	Defense: Block root login over SSH
		When common root credentials are presented over SSH
		Then root login over SSH is blocked

		Attack: Use common SSH accounts to login
			When root login over SSH is blocked
			Then well known SSH accounts that have sudo privileges are used to login # Example: raspi:rasberry
	`
	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var model model.Model
	model.Init(gherkinModel.Feature)
	assert.Equal(t, model.Defenses["Block root login over SSH"],
		model.
			Attacks["Use common SSH accounts to login"].
			Actions["root login over SSH is blocked"].Item[0])
}

func TestRule_ConnectAttackDefensebyTag(t *testing.T) {
	input := `
	Model: Bruteforce passwords

		@success @rainbow-tables
		Attack: Use rainbow tables for login

		@todo @rainbow-tables
		Defense: Cancel login if password found in rainbow tables
	`
	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var model model.Model
	model.Init(gherkinModel.Feature)
	assert.Equal(t, "Use rainbow tables for login",
		model.
			Defenses["Cancel login if password found in rainbow tables"].
			Tags["@rainbow-tables"][0].Title)
}
