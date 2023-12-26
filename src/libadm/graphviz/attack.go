package graphviz

import (
	"libadm/model"
	"strings"
)

////////////////////////////////////////
// Exported functions

func GetAttackCode(a *model.Attack, unmitigatedAttacks map[string]interface{}, config GraphvizConfig) (codeLines map[string]string, preConditions []string) {
	codeLines = make(map[string]string)
	for title, p := range a.PreConditions {
		if p.Item == nil { // add it only if precondition doesn't refer to another object (attack, defense, etc.)
			preConditions = append(preConditions, generateID(title))
			code := generateNodeSpec(p.Step.Statement, config.PreConditions, false, false)
			codeLines[generateID(title)] = "[" + strings.TrimSpace(code) + "]"
		}
	}

	if _, present := unmitigatedAttacks[a.Title]; present {
		codeLines[generateID(a.Title)] = "[" + strings.TrimSpace(generateNodeSpec(a.Title, config.UnMitigatedAttack, false, true)) + "]"
	} else if len(a.PreConditions) == 0 && len(a.Actions) == 0 && len(a.Results) == 0 {
		codeLines[generateID(a.Title)] = "[" + strings.TrimSpace(generateNodeSpec(a.Title, config.EmptyAttack, true, false)) + "]"
	} else {
		codeLines[generateID(a.Title)] = "[" + strings.TrimSpace(generateNodeSpec(a.Title, config.Attack, false, false)) + "]"
	}

	return
}
