package test

import (
	"libadm/loaders"
	"libadm/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

/////////////////////////////////////////////////
// Unit tests

func TestDefenseFail(t *testing.T) {
	input := `
	Model: Model with just an attack spec
		Attack: Attack passed-off as defense
	`

	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var d model.Defense
	err = d.Init(gherkinModel.Feature.Children[0].Scenario)

	assert.Equal(t, err.Error(), "Expected 'Defense', got 'Attack'")
}

func TestDefense_TitleOnly(t *testing.T) {
	input := `
	Model: Model with just a defense spec
		Defense: We must defend our systems
	`

	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var d model.Defense
	err = d.Init(gherkinModel.Feature.Children[0].Scenario)

	assert.Nil(t, err)
	assert.Equal(t, len(d.Tags), 0)
	assert.Equal(t, d.Title, "We must defend our systems")
}

func TestDefense_ButStatementsNotSupported(t *testing.T) {
	input := `
	Model: Model with a defense spec
		Defense: We must defend our systems
			But 'but' statement is not supported in ADM
	`

	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var d model.Defense
	err = d.Init(gherkinModel.Feature.Children[0].Scenario)

	assert.Equal(t, err.Error(), "unsupported keyword - 'But'")
}

func TestDefense_DuplicateStatements(t *testing.T) {
	testVectors := map[string][]string{
		"SamePrecondition": {`Model: Model with just one defense
			Defense: We must defend our systems
			Given system is exposed to internet
			Given system is exposed to internet
	`, "precondition - 'system is exposed to internet' is already part of this defense"},
		"SameAction": {`Model: Model with just one defense
			Defense: We must defend our systems
			When nmap scan is run
			When nmap scan is run
	`, "action - 'nmap scan is run' is already part of this defense"},
		"SameResult": {`Model: Model with just one defense
			Defense: We must defend our systems
			Then block all ports except the one required for this service
			Then block all ports except the one required for this service
	`, "result - 'block all ports except the one required for this service' is already part of this defense"},
		"SamePreconditionAsAnd": {`Model: Model with just one defense
			Defense: We must defend our systems
			Given system is exposed to internet
			And system is exposed to internet
	`, "precondition - 'system is exposed to internet' is already part of this defense"},
		"SameActionAsAnd": {`Model: Model with just one defense
			Defense: We must defend our systems
			When nmap scan is run
			And nmap scan is run
	`, "action - 'nmap scan is run' is already part of this defense"},
		"SameResultAsAnd": {`Model: Model with just one defense
			Defense: We must defend our systems
			Then block all ports except the one required for this service
			And block all ports except the one required for this service
	`, "result - 'block all ports except the one required for this service' is already part of this defense"},
	}

	for name, args := range testVectors {
		t.Run(name, func(t *testing.T) {
			input := args[0]
			expected := args[1]
			gherkinModel, err := loaders.LoadGherkinContent(input)
			if err != nil {
				t.Error(err)
			}

			var d model.Defense
			err = d.Init(gherkinModel.Feature.Children[0].Scenario)

			assert.Equal(t, expected, err.Error())
		})
	}
}

func TestDefense_FullBody(t *testing.T) {
	input := `
	Model: Model with single defense spec
		Defense: We must defend our systems
			Given a precondition for defense
			When an attack step is executed
			Then the defense step must be executed
	`

	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var d model.Defense
	err = d.Init(gherkinModel.Feature.Children[0].Scenario)

	assert.Nil(t, err)
	assert.Equal(t, len(d.Tags), 0)
	assert.Equal(t, d.Title, "We must defend our systems")
	assert.NotNil(t, d.PreConditions["a precondition for defense"])
	assert.NotNil(t, d.Actions["an attack step is executed"])
	assert.NotNil(t, d.Results["the defense step must be executed"])
}

func TestDefense_Complete(t *testing.T) {
	input := `
	Model: Model with a full defense spec.
		Defense: We must defend our systems
			Given a precondition for defense
			And another precondition
			When an attack step is executed
			And another attack step is also completed
			Then the defense step must be executed
			And additional logging must be included
	`

	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var d model.Defense
	err = d.Init(gherkinModel.Feature.Children[0].Scenario)

	assert.Nil(t, err)
	assert.Equal(t, len(d.Tags), 0)
	assert.Equal(t, d.Title, "We must defend our systems")
	assert.Equal(t, len(d.PreConditions), 2)
	assert.NotNil(t, d.PreConditions["a precondition for defense"])
	assert.NotNil(t, d.PreConditions["another precondition"])
	assert.Equal(t, len(d.Actions), 2)
	assert.NotNil(t, d.Actions["an attack step is executed"])
	assert.NotNil(t, d.Actions["another attack step is also completed"])
	assert.Equal(t, len(d.Results), 2)
	assert.NotNil(t, d.Results["the defense step must be executed"])
	assert.NotNil(t, d.Results["additional logging must be included"])
}
