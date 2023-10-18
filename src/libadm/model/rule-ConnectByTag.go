package model

// RULE: All attacks and defenses that have the same tag are connected. Defense always follows attack.
func ConnectByTags(a *Attack, d *Defense) {
	for _, aTag := range a.Tags {
		if foundSpecialTag(aTag) {
			continue
		}
		for dTag := range d.Tags {
			if foundSpecialTag(dTag) {
				continue
			} else if aTag == dTag {
				d.Tags[dTag] = append(d.Tags[dTag], a)
			}
		}
	}
}

func ConnectSingleAttackToMultiDefensesByTag(a *Attack, defenses map[string]*Defense) {
	for _, defense := range defenses {
		ConnectByTags(a, defense)
	}
}

func ConnectMultiAttacksToSingleDefensByTag(attacks map[string]*Attack, defense *Defense) {
	for _, a := range attacks {
		ConnectByTags(a, defense)
	}
}

func MultiConnectAttacksAndDefensesByTag(attacks map[string]*Attack, defenses map[string]*Defense) {
	for _, attack := range attacks {
		for _, defense := range defenses {
			ConnectByTags(attack, defense)
		}
	}
}

func foundSpecialTag(tag string) bool {
	specialTags := []string{
		"@yolosec", // Bad pre-condition / defense choice
		"@success", // Attack that will succeed even in the presence of defenses
		"@todo",    // Defense / Policy that are yet to be implemented
		"@done",    // Defense / Policy that has been implemented
	}

	for _, t := range specialTags {
		if t == tag {
			return true
		}
	}

	return false
}
