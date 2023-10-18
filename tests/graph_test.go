package test

import (
	"libadm/graph"
	"libadm/loaders"
	"libadm/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddModelWithNullModel(t *testing.T) {
	input := `
	Model: A sample model
	`
	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var m model.Model
	err = m.Init(gherkinModel.Feature)
	if err != nil {
		t.Error(err)
	}

	var g graph.Graph
	g.Init()
	err = g.AddModel(nil)
	assert.NotNil(t, err)
}

func TestGraphWithModelTitleOnly(t *testing.T) {
	input := `
	Model: A sample model
	`
	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var m model.Model
	err = m.Init(gherkinModel.Feature)
	if err != nil {
		t.Error(err)
	}

	var g graph.Graph
	g.Init()
	err = g.AddModel(&m)

	assert.Nil(t, err)
}

func TestGraphWithFullModel(t *testing.T) {
	input := `
	Model: Friends fight
		Assumption: Friends exist
			Given there are two friends
			And a trust relation exists

		@success
		Attack: Cheat on your friend
			When you are not sure about friend's honesty
			Then cheat on your friend

		Attack: Friend's cheating is caught
			When cheating is caught
			Then confront them

		Defense: Hide your cheat
			When cheating is caught
			Then cook up a story to convince friend

		Policy: Honesty is the best policy
			Assumption: Everyone is honest
				Given Many people surround you
				And everyone is honest
			Defense: Test honesty
				When you are not sure about friend's honesty
				Then test their honesty
	`
	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var m model.Model
	err = m.Init(gherkinModel.Feature)
	if err != nil {
		t.Error(err)
	}

	var g graph.Graph
	g.Init()
	g.AddModel(&m)
	assert.Nil(t, err)
	// Check successors to reality
	assert.Equal(t, m.Assumptions["Friends exist"], g.RealitySuccessors["Friends exist"])
	assert.Equal(t, m.Policies["Honesty is the best policy"].Assumptions["Everyone is honest"], g.RealitySuccessors["Everyone is honest"])
	assert.Equal(t, m.Attacks["Cheat on your friend"], g.RealitySuccessors["Cheat on your friend"])
	// Check predecessors to 'attacker wins'
	assert.Equal(t, m.Attacks["Cheat on your friend"], g.AttackerWinsPredecessors["Cheat on your friend"])
	// Check attack-defense relations
	assert.Equal(t, m.Policies["Honesty is the best policy"].Defenses["Test honesty"].Actions["you are not sure about friend's honesty"].Item[0], m.Attacks["Cheat on your friend"])
	assert.Equal(t, m.Defenses["Hide your cheat"].Actions["cheating is caught"].Item[0], m.Attacks["Friend's cheating is caught"])
}

func TestGraphWithTwoModels(t *testing.T) {
	model1 := `
	Model: Friends fight
		Assumption: Adam and Bob are friends
			Given Adam and Bob
			And a trust relation between Adam and Bob

		@success
		Attack: Adam cheats on Bob
			When Adam is not sure about Bob's honesty
			Then Adam cheats on Bob

		Attack: Adam's cheating is caught
			Given Adam cheats on Bob
			When Bob finds out about Adam's cheating
			Then Bob confronts Adam

		Defense: Adam hides the cheat
			When Bob finds out about Adam's cheating
			Then Adam cooks-up a story to convince Bob

		Policy: Honesty is the best policy
			Assumption: Adam and Bob are generally honest
				Given Adam and Bob are good people
				And Adam and Bob have a history of being honest with each other
			Defense: Test honesty
				When Adam is not sure about Bob's honesty
				Then Adam must test Bob's honesty
	`
	model2 := `
	Model: Divide and rule
		Assumption: Friction between friends
			Given Adam and Bob are friends
			And Charlie doesn't like Adam and Bob

		Attack: Charlie sows doubt
			When Adam cooks-up a story to convince Bob
			Then Charlie secretly point out flaws in Adam's story to Bob

		Defense: Bob tries to verify Adam's story
			Given Adam hides the cheat
			When Bob wants to confirm the story
			Then Bob asks probing questions to Adam

		# Defense that doesn't mitigate anything
		Defense: Separate Charlie from others
			When Charlie tries to interfere
			Then separate him from Adam and Bob
	`

	gherkin1, err := loaders.LoadGherkinContent(model1)
	if err != nil {
		t.Error(err)
	}
	gherkin2, err := loaders.LoadGherkinContent(model2)
	if err != nil {
		t.Error(err)
	}

	var m1, m2 model.Model
	err = m1.Init(gherkin1.Feature)
	if err != nil {
		t.Error(err)
	}
	err = m2.Init(gherkin2.Feature)
	if err != nil {
		t.Error(err)
	}

	var g graph.Graph
	g.Init()
	err = g.AddModel(&m1)
	assert.Nil(t, err)
	err = g.AddModel(&m2)
	assert.Nil(t, err)

	////////////////////////////////////////
	// Confirm graph and inter-model connections (internal connections within model are skipped)

	// Reality
	assert.Equal(t, g.RealitySuccessors["Adam and Bob are friends"], m1.Assumptions["Adam and Bob are friends"])
	assert.Equal(t, g.RealitySuccessors["Adam and Bob are generally honest"], m1.Policies["Honesty is the best policy"].Assumptions["Adam and Bob are generally honest"])
	assert.Equal(t, g.RealitySuccessors["Friction between friends"], m2.Assumptions["Friction between friends"])
	assert.Equal(t, g.RealitySuccessors["Adam cheats on Bob"], m1.Attacks["Adam cheats on Bob"])
	assert.Equal(t, g.RealitySuccessors["Separate Charlie from others"], m2.Defenses["Separate Charlie from others"])
	// Attacker Wins
	assert.Equal(t, g.AttackerWinsPredecessors["Adam cheats on Bob"], m1.Attacks["Adam cheats on Bob"])
	assert.Equal(t, g.AttackerWinsPredecessors["Charlie sows doubt"], m2.Attacks["Charlie sows doubt"])
	// Model inter-connect
	assert.Equal(t,
		m2.Attacks["Charlie sows doubt"].Actions["Adam cooks-up a story to convince Bob"].Item[0],
		m1.Defenses["Adam hides the cheat"])
	assert.Equal(t,
		m2.Defenses["Bob tries to verify Adam's story"].PreConditions["Adam hides the cheat"].Item[0],
		m1.Defenses["Adam hides the cheat"])
}

func TestGraphContainsCheck(t *testing.T) {
	model1 := `
	Model: Friends fight
		Assumption: Adam and Bob are friends
			Given Adam and Bob
			And a trust relation between Adam and Bob
`

	gherkin1, err := loaders.LoadGherkinContent(model1)
	if err != nil {
		t.Error(err)
	}

	var m1 model.Model
	err = m1.Init(gherkin1.Feature)
	if err != nil {
		t.Error(err)
	}

	var g graph.Graph
	g.Init()
	err = g.AddModel(&m1)
	assert.Nil(t, err)
	assert.True(t, g.ContainsModel("Friends fight"))
	assert.False(t, g.ContainsModel("Friends don't fight"))
	assert.Equal(t, &m1, g.GetModel("Friends fight"))
}
