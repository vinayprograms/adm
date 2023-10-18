package model

import (
	"log"
	"strings"
)

type ActorType string

const (
	Attacker ActorType = "attacker"
	Defender ActorType = "defender"
)

type ModelActor struct {
	Actor    ActorType
	Intents  []string
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
		log.Print("ModelActor: Found", len(actors), "Actors. Only the first one will be used")
	}

	p.Init(actors[0].Actor, actors[0].Intents, actors[0].Purposes)
}

func parseActors(specification string) (actors []*ModelActor) {
	lines := strings.Split(specification, "\n")
	var currentActor *ModelActor
	var currentClauseType string
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), "as an attacker") ||
			strings.Contains(strings.ToLower(line), "as a attacker") {
			p := ModelActor{Actor: Attacker}
			actors = append(actors, &p)
			currentActor = &p
		} else if strings.Contains(strings.ToLower(line), "as a defender") {
			p := ModelActor{Actor: Defender}
			actors = append(actors, &p)
			currentActor = &p
		} else if strings.Contains(strings.ToLower(line), "i want to") {
			currentClauseType = "intent"
			if currentActor == nil {
				log.Print("ERROR: No actor identified for the clause - \"", line, "\"")
				break
			}
			currentActor.Intents = append(currentActor.Intents, strings.TrimSpace(line))
		} else if strings.Contains(strings.ToLower(line), "so that") {
			currentClauseType = "purpose"
			if currentActor == nil {
				log.Print("ERROR: No actor identified for the clause - \"", line, "\"")
				break
			}
			currentActor.Purposes = append(currentActor.Purposes, strings.TrimSpace(line))
		} else if strings.Contains(strings.ToLower(line), "and") {
			if currentActor == nil {
				log.Print("ERROR: No actor identified for the clause - \"", line, "\"")
				break
			}
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
