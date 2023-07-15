package model

// RULE: An attack thwarts a defense if defense's 'Then' statement is attack's 'When' statement
func ConnectIfAttackThwartsDefense(a *Attack, d *Defense) bool {
	for aActionTitle, aActionRef := range a.Actions {
		if match := d.Results[aActionTitle]; match != nil {
			aActionRef.Item = append (aActionRef.Item, d)
			return true
		}
	}

	return false
}

// RULE: A defense mitigates an attack if attack's 'When' statement is defense's 'When' statement
func ConnectIfDefenseMitigatesAttack(a *Attack, d *Defense) bool {
	for dActionTitle, dActionRef := range d.Actions {
		if match := a.Actions[dActionTitle]; match != nil {
			dActionRef.Item = append(dActionRef.Item, a)
			return true
		}
	}

	return false
}

// RULE: A defense is an incident-response to an attack if attack's 'Then' statement is defense's 'When' statement
func ConnectIfDefenseRespondsToAttack(a *Attack, d *Defense) bool {
	for dActionTitle, dActionRef := range d.Actions {
		if match := a.Results[dActionTitle]; match != nil {
			dActionRef.Item = append(dActionRef.Item, a)
			return true
		}
	}

	return false
}

func ConnectDefenseToAttacks(d *Defense, attacks map[string]*Attack) {
	for _, attack := range attacks {
		ConnectIfAttackThwartsDefense(attack, d)
		ConnectIfDefenseMitigatesAttack(attack, d)
		ConnectIfDefenseRespondsToAttack(attack, d)
	}
}


func ConnectAttackToDefenses(a *Attack, defenses map[string]*Defense) {
	for _, defense := range defenses {
		ConnectIfAttackThwartsDefense(a, defense)
		ConnectIfDefenseMitigatesAttack(a, defense)
		ConnectIfDefenseRespondsToAttack(a, defense)
	}
}