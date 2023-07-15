package export

import (
	"fmt"
	"libadm/model"
)

////////////////////////////////////////
// Common internal functions

func prepareHeaderLines(title string, actors []*model.ModelActor, typeFilter model.ActorType) (lines []string) {
	lines = append(lines, "Feature: " + title)
	for _, actor := range actors {
		if actor.Actor != typeFilter {
			continue
		}
		lines = append(lines, genrateTabs(1) + "As a " + fmt.Sprint(actor.Actor))
		for _, intent := range actor.Intents {
			lines = append(lines, genrateTabs(1) + intent)
		}
		for _, goal := range actor.Purposes {
			lines = append(lines, genrateTabs(1) + goal)
		}
	}
	return lines
}

func generateAssumptionLines(assumption *model.Assumption) (lines []string) {
	lines = append(lines, genrateTabs(1) + "Background: " + assumption.Title)
	for _, preCondition := range assumption.PreConditions {
		lines = append(lines, prepareStatement(preCondition)...)
	}
	return
}

func generateResultLines(results map[string]*model.Step) (lines []string) {
	var additionalLines []string
	for _, a := range results {
		if a.Keyword != "Then" {
			additionalLines = append(additionalLines, prepareStatement(a)...)
		} else {
			lines = append(lines, prepareStatement(a)...)
		}
	}
	lines = append(lines, additionalLines...)
	return
}

func prepareStatement(step *model.Step) (lines []string) {
	lines = append(lines, genrateTabs(2) + step.Keyword + " " + step.Statement)
	if step.DocString != "" {
		lines = append(lines, genrateTabs(2) + "\"\"\"" + step.DocStringType)
		lines = append(lines, step.DocString)
		lines = append(lines, genrateTabs(2) + "\"\"\"")
	}
	return lines
}

////////////////////////////////////////
// Helper functions

func genrateTabs(tabCount uint) (result string) {
	var i uint
	for i = 0; i < tabCount; i++ {
		result += "  "
	}
	return
}