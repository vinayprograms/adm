package export

import (
	"libadm/model"
	"strings"
)

// exports attacks from a model as gauntlt file contents (one file per model).
// Assumptions, if present, are included.
func ExportAttacks(model *model.Model) (lines []string) {
	lines = append(lines, prepareHeaderLines(model.Title, model.Actors, "attacker")...)
	for _, assumption := range model.Assumptions {
		lines = append(lines, generateAssumptionLines(assumption)...)
	}
	for _, attack := range model.Attacks {
		lines = append(lines, generateAttackLines(attack)...)
	}
	return
}

////////////////////////////////////////
// Internal functions

func generateAttackLines(attack *model.Attack) (lines []string) {
	if len(attack.Tags) > 0 {
		lines = append(lines, genrateTabs(1) + strings.Join(attack.Tags, " "))
	}
	lines = append(lines, genrateTabs(1) + "Scenario: " + attack.Title)
	lines = append(lines, generatePreConditionLinesForAttacks(attack.PreConditions)...)
	lines = append(lines, generateActionLinesForAttacks(attack.Actions)...)
	lines = append(lines, generateResultLines(attack.Results)...)
	return
}

func generatePreConditionLinesForAttacks(preConditions map[string]*model.ModelLink[[]*model.Attack]) (lines []string) {
	var additionalLines []string
	for _, p := range preConditions {
		if p.Step.Keyword != "Given" {
			additionalLines = append(additionalLines, prepareStatement(p.Step)...)
		} else {
			lines = append(lines, prepareStatement(p.Step)...)
		}
	}
	lines = append(lines, additionalLines...)
	return
}

func generateActionLinesForAttacks(actions map[string]*model.ModelLink[[]*model.Defense]) (lines []string) {
	var additionalLines []string
	for _, a := range actions {
		if a.Step.Keyword != "When" {
			additionalLines = append(additionalLines, prepareStatement(a.Step)...)
		} else {
			lines = append(lines, prepareStatement(a.Step)...)
		}
	}
	lines = append(lines, additionalLines...)
	return
}