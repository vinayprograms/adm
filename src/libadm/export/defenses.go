package export

import (
	"libadm/model"
)

// exports attacks from a model as gauntlt file contents (one file per model).
// Assumptions, if present, are included.
func ExportDefenses(model *model.Model) (lines []string) {
	lines = append(lines, prepareHeaderLines(model.Title, model.Actors, "defender")...)
	for _, assumption := range model.Assumptions {
		lines = append(lines, generateAssumptionLines(assumption)...)
	}
	for _, defense := range model.Defenses {
		lines = append(lines, generateDefenseLines(defense)...)
	}
	return
}

////////////////////////////////////////
// Internal functions

func generateDefenseLines(defense *model.Defense) (lines []string) {
	if len(defense.Tags) > 0 {
		var tagLine string
		for tag := range defense.Tags {
			tagLine += tag + " "
		}
		lines = append(lines, genrateTabs(1) + tagLine)
	}
	lines = append(lines, genrateTabs(1) + "Scenario: " + defense.Title)
	lines = append(lines, generatePreConditionLinesForDefenses(defense.PreConditions)...)
	lines = append(lines, generateActionLinesForDefenses(defense.Actions)...)
	lines = append(lines, generateResultLines(defense.Results)...)
	return
}

func generatePreConditionLinesForDefenses(preConditions map[string]*model.ModelLink[[]*model.Defense]) (lines []string) {
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

func generateActionLinesForDefenses(actions map[string]*model.ModelLink[[]*model.Attack]) (lines []string) {
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