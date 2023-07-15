package args

import (
	"fmt"
	"libadm/model"
	"strings"
)

type attackSummaryCommand struct {
	model *model.Model
}

type defenseSummaryCommand struct {
	model *model.Model
}

type preemtiveDefenseSummaryCommand struct {
	model *model.Model
}

type incidentResponseSummaryCommand struct {
	model *model.Model
}

func (a attackSummaryCommand) execute() error {
	for attackTitle := range a.model.Attacks {
		fmt.Println("\tATTACK: " + attackTitle)
	}

	return nil
}

func (d defenseSummaryCommand) execute() error {
	for defenseTitle := range d.model.Defenses {
		fmt.Println("\tDEFENSE: " + defenseTitle)
	}
	
	for policyTitle, policy := range d.model.Policies {
		if len(policy.Defenses) > 0 {
			fmt.Println("\tPOLICY: " + policyTitle)
			for defenseTitle := range policy.Defenses {
				fmt.Println("\t\tDEFENSE: " + defenseTitle)
			}
		}
	}

	return nil
}

func (d preemtiveDefenseSummaryCommand) execute() error {
	lines := getPreemtiveDefenseLines(d.model.Defenses, 1, d.model.Attacks)
	
	for policyTitle, policy := range d.model.Policies {
		pdLines := getPreemtiveDefenseLines(policy.Defenses, 2, d.model.Attacks)
		if len(pdLines) > 0 {
			lines = append(lines, "\tPOLICY: " + policyTitle)
			lines = append(lines, pdLines...)
		}
	}
	output := strings.Join(lines, "\n")
	fmt.Println(output)

	return nil
}

func (d incidentResponseSummaryCommand) execute() error {
	lines := getIncidentResponseLines(d.model.Defenses, 1, d.model.Attacks)
	
	for policyTitle, policy := range d.model.Policies {
		irLines := getIncidentResponseLines(policy.Defenses, 2, d.model.Attacks)
		if len(irLines) > 0 {
			lines = append(lines, "\tPOLICY: " + policyTitle)
			lines = append(lines, irLines...)
		}
	}
	output := strings.Join(lines, "\n")
	fmt.Println(output)

	return nil
}

/////////////////////////////////////////////////
// Helper functions

func getPreemtiveDefenseLines(defenses map[string]*model.Defense, tabspaces int, attacks map[string]*model.Attack) (lines []string) {
	for defenseTitle, defense := range defenses {
		for _, attack := range attacks {
			for dActionTitle := range defense.Actions {
				if match := attack.Actions[dActionTitle]; match != nil {
					var begin string
					for i := 0; i < tabspaces; i++ { begin += "\t" }
					lines = append(lines, begin + "DEFENSE: " + defenseTitle)
				}
			}
		}
	}
	return
}

func getIncidentResponseLines(defenses map[string]*model.Defense, tabspaces int, attacks map[string]*model.Attack) (lines []string) {
	for defenseTitle, defense := range defenses {
		for _, attack := range attacks {
			for dActionTitle := range defense.Actions {
				if match := attack.Results[dActionTitle]; match != nil {
					var begin string
					for i := 0; i < tabspaces; i++ { begin += "\t" }
					lines = append(lines, begin + "DEFENSE: " + defenseTitle)
				}
			}
		}
	}
	return
}