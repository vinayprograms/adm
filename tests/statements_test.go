package test

import (
	"libadm/loaders"
	"libadm/model"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

/////////////////////////////////////////////////
// Unit tests

func TestStep_NilInitialization(t *testing.T) {
	var s model.Step
	err := s.Init(nil)
	assert.Equal(t, err.Error(), "expected a 'Given'/'When'/'Then'/'And' statement. Got 'nil'")
}

func TestStepWithDocString(t *testing.T) {
	input := `
	Model: Model with just one attack
		Attack: First step
			When attack step is executed
			"""markdown
			Some docstring
			"""
	`

	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var s model.Step
	err = s.Init(gherkinModel.Feature.Children[0].Scenario.Steps[0])
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, s.Keyword, "When")
	assert.Equal(t, s.Statement, "attack step is executed")
	assert.Equal(t, strings.TrimSpace(s.DocStringType), "markdown")
	assert.Equal(t, strings.TrimSpace(s.DocString), "Some docstring")
}

func TestStepWithDataTable(t *testing.T) {
	input := `
	Model: Model with just one attack
		Attack: First step
			When <attack-step> is executed
			| attack-step |
			| step-1      |
			| step-2      |
			| step-3      |
	`

	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var s model.Step
	err = s.Init(gherkinModel.Feature.Children[0].Scenario.Steps[0])
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, s.Keyword, "When")
	assert.Equal(t, s.Statement, "<attack-step> is executed")
	assert.Equal(t, s.DataTable.Heads[0], "attack-step")
	assert.Equal(t, s.DataTable.Rows[0][0], "step-1")
	assert.Equal(t, s.DataTable.Rows[1][0], "step-2")
	assert.Equal(t, s.DataTable.Rows[2][0], "step-3")
}