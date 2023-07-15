package test

import (
	"libadm/loaders"
	"libadm/model"
	"sources"
	"testing"

	"github.com/cucumber/messages-go/v16"
	"github.com/stretchr/testify/assert"
)

/////////////////////////////////////////////////
// Unit tests

func TestNullModel(t *testing.T) {
	var model model.Model
	err := model.Init(nil)
	assert.Equal(t, "cannot initialize a null model object", err.Error())
}

func TestWithEmptyModel(t *testing.T) {
	gherkinModel, err := loaders.LoadGherkinContent("")
	if err != nil {
		t.Error(err)
	}

	var model model.Model
	model.Init(gherkinModel.Feature)
}

func TestSimpleModelInit(t *testing.T) {
	input := `
	Model: A sample model
		As an attacker
		I want to compromize the target
		And maintain control over it
		So that I can use it for malicious purposes
		And extort the target

		As a defender 
		I want to secure the target
		So that it can only be used for legit purposes
	`
	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var model model.Model
	model.Init(gherkinModel.Feature)

	assert.Equal(t, model.Title, "A sample model")
	assert.Nil(t, model.Tags)
	assert.Equal(t, len(model.Actors), 2)
}

func TestManualActorAddition(t *testing.T) {
	input := `
	As an attacker
	I want to compromize the target
	And maintain control over it
	So that I can use it for malicious purposes
	And extort the target
	`
	var Actor model.ModelActor
	Actor.ParseAndInit(input)
	var model model.Model
	model.AddActor(&Actor)
}

func TestModelWithoutActors(t *testing.T) {
	input := `
	Model: A sample model
		
		Attack: Sample Attack
		Defense: Sample Defense
	`
	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var model model.Model
	model.Init(gherkinModel.Feature)

	assert.Equal(t, model.Title, "A sample model")
	assert.Nil(t, model.Tags)
	assert.Len(t, model.Attacks, 1)
	assert.Len(t, model.Defenses, 1)
	assert.Contains(t, model.Attacks, "Sample Attack")
	assert.Contains(t, model.Defenses, "Sample Defense")
}

func TestModel_Full(t *testing.T) {
	input := `
	Model: Password security analysis
		
		As an attacker
		I want to compromize the target
		And maintain control over it
		So that I can use it for malicious purposes
		And extort the target

		As a defender 
		I want to secure the target
		So that it can only be used for legit purposes

		Assumption: Some common assumption
			Given part-1 of the assumption
			And part-2 of the assumption

		@password-len
		Attack: Bruteforce short passwords
			Given password is shorter than 8 characters
			When password is enumerated using word generator
			Then one of the words match
			And the password is revealed

		Defense: Another dummy defense
			Given some precondition
			When some attack happens
			Then execute a specific defense

		Policy: Minimum password requirements
			@password-len
			Defense: password must be atleast 8 characters long
	`
	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var model model.Model
	model.Init(gherkinModel.Feature)
}

func TestModel_FromFile(t *testing.T) {
	src := sources.GetSource("./examples/others/bodyless.adm")
	input, _ := src.GetContent("./examples/others/bodyless.adm")
	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var model model.Model
	model.Init(gherkinModel.Feature)
}

func TestModelWithPolicy(t *testing.T) {
	input := `
	Model: A sample model
		Policy: Honesty is the best policy
	`
	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var model model.Model
	model.Init(gherkinModel.Feature)

	assert.Equal(t, model.Title, "A sample model")
	assert.Nil(t, model.Tags)
	assert.Len(t, model.Policies, 1)
	assert.Contains(t, model.Policies, "Honesty is the best policy")
}

/*func TestModelWithDuplicatePolicies(t *testing.T) {
	input := `
	Model: A sample model
		Policy: Honesty is the best policy
		Policy: Honesty is the best policy
	`
	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var model model.Model
	err = model.Init(gherkinModel.Feature)

	assert.Equal(t, model.Title, "A sample model")
	assert.Nil(t, model.Tags)
	assert.Equal(t, "policy - 'Honesty is the best policy' is already part of this model", err.Error())
	assert.Len(t, model.Policies, 1)
	assert.Contains(t, model.Policies, "Honesty is the best policy")
}*/

func TestModelWithAssumption(t *testing.T) {
	input := `
	Model: A sample model
		Assumption: Attacker can get into the network
	`
	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var model model.Model
	model.Init(gherkinModel.Feature)

	assert.Equal(t, model.Title, "A sample model")
	assert.Nil(t, model.Tags)
	assert.Len(t, model.Assumptions, 1)
	assert.Contains(t, model.Assumptions, "Attacker can get into the network")
}

/*func TestModelWithDuplicateAttacks(t *testing.T) {
	input := `
	Model: A sample model
		Attack: Sample Attack 1
		Attack: Sample Attack 1
	`
	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var model model.Model
	err = model.Init(gherkinModel.Feature)

	assert.Equal(t, model.Title, "A sample model")
	assert.Nil(t, model.Tags)
	assert.Len(t, model.Attacks, 1)
	assert.Contains(t, model.Attacks, "Sample Attack 1")
	assert.Equal(t, "attack - 'Sample Attack 1' is already part of this model", err.Error())
}*/

/*func TestModelWithDuplicateDefenses(t *testing.T) {
	input := `
	Model: A sample model
		Defense: Sample Defense 1
		Defense: Sample Defense 1
	`
	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var model model.Model
	err = model.Init(gherkinModel.Feature)

	assert.Equal(t, model.Title, "A sample model")
	assert.Nil(t, model.Tags)
	assert.Len(t, model.Defenses, 1)
	assert.Contains(t, model.Defenses, "Sample Defense 1")
	assert.Equal(t, "defense - 'Sample Defense 1' is already part of this model", err.Error())
}*/

func TestModel_Duplicates(t *testing.T) {
	testVectors := map[string][]string {
		"SamePolicy": {`Model: Model with duplicates
			Policy: Honesty is the best policy
				Defense: Be Honest
			Policy: Honesty is the best policy
				Defense: Be Honest
	`,"policy - 'Honesty is the best policy' is already part of this model"},
	"SameAssumption": {`Model: Model with duplicates
			Assumption: Don't assume everything
				Given a lot of assumptions
			Assumption: Don't assume everything
				Given a lot of assumptions
	`,"Parser errors:\n(4:4): expected: #EOF, #TableRow, #DocStringSeparator, #StepLine, #TagLine, #ScenarioLine, #RuleLine, #Comment, #Empty, got '\t\t\tAssumption: Don't assume everything'"},
	"SameDefense": {`Model: Model with duplicates
			Defense: We must defend our systems
			Defense: We must defend our systems
	`,"defense - 'We must defend our systems' is already part of this model"},
	"SameAttack": {`Model: Model with duplicates
			Attack: Master of Pwn
			Attack: Master of Pwn
	`,"attack - 'Master of Pwn' is already part of this model"},
	}

	for name, args := range testVectors {
		t.Run(name, func(t *testing.T){
			input := args[0]
			expected := args[1]
			gherkinModel, err := loaders.LoadGherkinContent(input)
			if err != nil {
				if name == "SameAssumption" {
					assert.Equal(t, expected, err.Error())
					return
				}else {
					t.Error(err)
				}
			}

			var m model.Model
			err = m.Init(gherkinModel.Feature)

			assert.Equal(t, expected, err.Error())
		})
	}

	// Specific test to check when same assumption structure is 
	// passed twice to the model
	f := messages.Feature {
		Language: "en",
		Keyword: "Model",
		Name: "Model with duplicate Assumptions",
	}
	f.Children = append(f.Children, &messages.FeatureChild{
		Background: &messages.Background {
			Keyword: "Assumption",
			Name: "Don't assume everything",
		},
	})
	var model model.Model
	err := model.Init(&f)
	assert.Nil(t, err)
	err = model.Init(&f)
	assert.Equal(t, "assumption - 'Don't assume everything' is already part of this model", err.Error())
}

func TestModel_Errors(t *testing.T) {
	testVectors := map[string][]string {
	"BadAssumption": {`Model: Model with duplicates
			Assumption: Don't assume everything
				When a lot of assumptions
	`,"Unexpected keyword - 'When' in Assumption specification"},
	}

	for name, args := range testVectors {
		t.Run(name, func(t *testing.T){
			input := args[0]
			expected := args[1]
			gherkinModel, err := loaders.LoadGherkinContent(input)
			if err != nil {
				t.Error(err)
			}

			var m model.Model
			err = m.Init(gherkinModel.Feature)

			assert.Equal(t, expected, err.Error())
		})
	}
}