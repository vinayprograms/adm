package test

import (
	"libadm/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActor(t *testing.T) {
	input := `
	As an attacker
	I want to compromize the target
	And maintain control over it
	So that I can use it for malicious purposes
	And extort the target
	`
	
	var Actor model.ModelActor
	Actor.ParseAndInit(input)

	assert.Equal(t, Actor.Actor, model.Attacker)
	assert.Equal(t, Actor.Intents, []string{"I want to compromize the target", "And maintain control over it"})
	assert.Equal(t, Actor.Purposes, []string{"So that I can use it for malicious purposes", "And extort the target"})
}

func TestWith2Actors(t *testing.T) {
	input := `
	As an attacker
	I want to compromize the target
	And maintain control over it
	So that I can use it for malicious purposes
	And extort the target

	As a defender 
	I want to secure the target
	So that it can only be used for legit purposes
	`
	
	var Actor model.ModelActor
	Actor.ParseAndInit(input)

	assert.Equal(t, Actor.Actor, model.Attacker)
	assert.Equal(t, Actor.Intents, []string{"I want to compromize the target", "And maintain control over it"})
	assert.Equal(t, Actor.Purposes, []string{"So that I can use it for malicious purposes", "And extort the target"})
}