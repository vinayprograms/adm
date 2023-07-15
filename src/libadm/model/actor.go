package model

import (
	"fmt"
	"strings"
)

type ActorType string
const (
	Attacker ActorType 	= "attacker"
	Defender ActorType 	= "defender"
)

type ModelActor struct {
	Actor ActorType
	Intents []string
	Purposes []string
}

func (p *ModelActor) Init(typ ActorType, intents []string, purposes []string) {
	p.Actor = typ
	p.Intents = intents
	p.Purposes = purposes
}

func (p *ModelActor) ParseAndInit(specification string) {
	actors := parseActors(specification)
	if len(actors) > 1 {
		fmt.Println("ModelActor: Found", len(actors), "Actors. Only the first one will be used")
	}

	p.Init(actors[0].Actor, actors[0].Intents, actors[0].Purposes)
}

func parseActors(specification string) (actors []*ModelActor) {
	lines := strings.Split(specification, "\n")
	var currentActor *ModelActor
	var currentClauseType string
	for _, line := range lines {
		if (strings.Contains(line, "As an attacker")) {
			p := ModelActor{Actor: Attacker}
			actors = append(actors, &p)
			currentActor = &p
		} else if (strings.Contains(line, "As a defender")) {
			p := ModelActor{Actor: Defender}
			actors = append(actors, &p)
			currentActor = &p
		} else if (strings.Contains(line, "I want to")) {
			currentClauseType = "intent"
			currentActor.Intents = append(currentActor.Intents, strings.TrimSpace(line))
		} else if (strings.Contains(line, "So that")) {
			currentClauseType = "purpose"
			currentActor.Purposes = append(currentActor.Purposes, strings.TrimSpace(line))
		} else if (strings.Contains(line, "And")) {
			switch currentClauseType {
			case "intent":
				currentActor.Intents = append(currentActor.Intents, strings.TrimSpace(line))
			case "purpose":
				currentActor.Purposes = append(currentActor.Purposes, strings.TrimSpace(line))
			}
		}
	}

	return actors
}