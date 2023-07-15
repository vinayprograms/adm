package model

// RULE: Attacks can be chained if one's 'Given' statement matches another's title
func ChainAttacks(newAttack *Attack, currentAttacks map[string]*Attack) {
	for title, currAttack := range currentAttacks {
		// forward connections
		for _, p := range newAttack.PreConditions {
			if p.Step.Statement == title {
				p.Item = append(p.Item, currAttack)
			}
		}
		// reverse connections
		for _, p := range currAttack.PreConditions {
			if p.Step.Statement == newAttack.Title {
				p.Item = append(p.Item, newAttack)
			}
		}
	}
}

// RULE: Defenses can be chained if one's 'Given' statement matches another's title
func ChainDefenses(newDefense *Defense, currentDefenses map[string]*Defense) {
	for title, currDefense := range currentDefenses {
		// forward connections
		for _, p := range newDefense.PreConditions {
			if p.Step.Statement == title {
				p.Item = append(p.Item, currDefense)
			}
		}
		// reverse connections
		for _, p := range currDefense.PreConditions {
			if p.Step.Statement == newDefense.Title {
				p.Item = append(p.Item, newDefense)
			}
		}
	}
}