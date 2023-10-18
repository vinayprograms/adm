package graphviz

import (
	"libadm/model"
	"strings"
)

////////////////////////////////////////
// Exported functions

func GetDefenseCode(d *model.Defense, config GraphvizConfig) (codeLines map[string]string, preConditions []string) {
	codeLines = make(map[string]string)
	for title, p := range d.PreConditions {
		if p.Item == nil { // add it only if precondition doesn't refer to another object (attack, defense, etc.)
			preConditions = append(preConditions, generateID(title))
			code := generateNodeSpec(p.Step.Statement, config.PreConditions, false)
			codeLines[generateID(title)] = "[" + strings.TrimSpace(code) + "]"
		}
	}
	if len(d.PreConditions) == 0 && len(d.Actions) == 0 && len(d.Results) == 0 {
		codeLines[generateID(d.Title)] = "[" + strings.TrimSpace(generateNodeSpec(d.Title, config.EmptyDefense, true)) + "]"
	} else if isIncidentResponse(d) {
		codeLines[generateID(d.Title)] = "[" + strings.TrimSpace(generateNodeSpec(d.Title, config.IncidentResponse, false)) + "]"
	} else {
		codeLines[generateID(d.Title)] = "[" + strings.TrimSpace(generateNodeSpec(d.Title, config.PreEmptiveDefense, false)) + "]"
	}

	return
}

////////////////////////////////////////
// Internal functions

func isIncidentResponse(d *model.Defense) bool {
	for _, a := range d.Actions {
		for _, other := range a.Item {
			for clause := range other.Results {
				if a.Step.Statement == clause {
					return true
				}
			}
		}
	}
	return false
}
