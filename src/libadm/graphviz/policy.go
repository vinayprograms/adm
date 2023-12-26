package graphviz

import "libadm/model"

func BuildPolicySubgraph(p *model.Policy, config GraphvizConfig) (s subgraph) {

	var allAssumptions []string
	for _, a := range p.Assumptions {
		allAssumptions = append(allAssumptions, CollectAssumptions(a)...)
	}
	title := buildTitleWithAssumptions(p.Title, allAssumptions, config)

	s.Init(generateID(p.Title), title, config.Policy, 1)

	for _, defense := range p.Defenses {
		// add nodes
		codeLines, defensePreConditions := GetDefenseCode(defense, config)
		for id, line := range codeLines {
			if _, present := s.Nodes[id]; !present {
				s.Nodes[id] = line
			}
		}
		s.preConditions = append(s.preConditions, defensePreConditions...)

		for _, pre := range defense.PreConditions {
			if pre.Item == nil { // Precondition
				connectAndAppend(s.Edges, "reality", generateID(pre.Step.Statement), "")
				//connectAndAppend(s.Edges, generateID(pre.Step.Statement), generateID(defense.Title), "")
			} else {
				for _, item := range pre.Item {
					connectAndAppend(s.Edges, generateID(item.Title), generateID(defense.Title), "")
				}
			}
		}
		for _, action := range defense.Actions { // connect policy defense and attacks (mitigation)
			for _, attack := range action.Item {
				connectAndAppend(s.Edges, generateID(attack.Title), generateID(defense.Title), "")
			}
		}

		// TODO: Add policy tags pointing to attacks

		for _, attacks := range defense.Tags { // if specific defense tag points to another attack/defense
			for _, attack := range attacks {
				connectAndAppend(s.Edges, generateID(attack.Title), generateID(defense.Title), "")
			}
		}
	}
	// Link each attack to policy if they are linked by tags. For this, we pick the first defense,
	// link it to that, but tell graphviz to stop the arrow at the policy box.
	var firstKey string
	for firstKey = range p.Defenses {
		break
	}
	for _, attacks := range p.Tags { // if attacks are linked to policies
		for _, attack := range attacks {
			connectAndAppend(s.Edges, generateID(attack.Title), generateID(firstKey), "lhead=cluster_"+generateID(p.Title))
		}
	}
	return
}
