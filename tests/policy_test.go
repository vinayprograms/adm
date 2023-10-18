package test

import (
	"libadm/loaders"
	"libadm/model"
	"testing"

	"github.com/cucumber/messages-go/v16"
	"github.com/stretchr/testify/assert"
)

/////////////////////////////////////////////////
// Unit tests

func TestWithPolicyHeading(t *testing.T) {
	input := `
	Model: Model with just one rule
		Policy: One rule to control them all!
	`

	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var policy model.Policy
	policy.Init(gherkinModel.Feature.Children[0].Rule)
	assert.Equal(t, policy.Title, "One rule to control them all!")
	assert.Zero(t, len(policy.Tags))
}

func TestPolicyWithSingleDefense(t *testing.T) {
	input := `
	Model: Model with just one rule
		Policy: One rule to control them all!
			Defense: Only one defense in this policy
				Given a precondition
				When an attack pattern is detected
				Then a specific defense must be executed
	`

	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var policy model.Policy
	policy.Init(gherkinModel.Feature.Children[0].Rule)
	assert.Equal(t, policy.Title, "One rule to control them all!")
	assert.Zero(t, len(policy.Tags))
	assert.NotNil(t, policy.Defenses["Only one defense in this policy"])
}

func TestPolicy_Full(t *testing.T) {
	input := `
	Model: Model with just one rule
		Policy: One rule to control them all!
			Assumption: Some common assumption
				Given part-1 of the assumption
			Defense: Only one defense in this policy
				Given a precondition
				When an attack pattern is detected
				Then a specific defense must be executed
	`

	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var policy model.Policy
	policy.Init(gherkinModel.Feature.Children[0].Rule)
	assert.Equal(t, policy.Title, "One rule to control them all!")
	assert.NotNil(t, policy.Assumptions["Some common assumption"])
	assert.NotNil(t, policy.Defenses["Only one defense in this policy"])
}

func TestPolicy_Duplicates(t *testing.T) {
	testVectors := map[string][]string{
		"SameAssumption": {`Model: Model with duplicates
		Policy: Honesty is the best policy
			Assumption: Don't assume everything
				Given a lot of assumptions
			Assumption: Don't assume everything
				Given a lot of assumptions
	`, "Parser errors:\n(5:4): expected: #EOF, #TableRow, #DocStringSeparator, #StepLine, #TagLine, #ScenarioLine, #RuleLine, #Comment, #Empty, got '\t\t\tAssumption: Don't assume everything'"},
		"SameDefense": {`Model: Model with duplicates
		Policy: Honesty is the best policy
			Defense: We must defend our systems
			Defense: We must defend our systems
	`, "defense - 'We must defend our systems' is already part of this model"},
	}

	for name, args := range testVectors {
		t.Run(name, func(t *testing.T) {
			input := args[0]
			expected := args[1]
			gherkinModel, err := loaders.LoadGherkinContent(input)
			if err != nil {
				if name == "SameAssumption" {
					assert.Equal(t, expected, err.Error())
					return
				} else {
					t.Error(err)
				}
			}

			var p model.Policy
			err = p.Init(gherkinModel.Feature.Children[0].Rule)

			assert.Equal(t, expected, err.Error())
		})
	}

	// Specific test to check when same assumption structure is
	// passed twice to the policy
	p := messages.Rule{
		Keyword: "Model",
		Name:    "Policy with duplicate Assumptions",
	}
	p.Children = append(p.Children, &messages.RuleChild{
		Background: &messages.Background{
			Keyword: "Assumption",
			Name:    "Don't assume everything",
		},
	})
	var policy model.Policy
	err := policy.Init(&p)
	assert.Nil(t, err)
	err = policy.Init(&p)
	assert.Equal(t, "assumption - 'Don't assume everything' is already part of this model", err.Error())
}
