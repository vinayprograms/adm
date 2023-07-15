package graphviz

import "libadm/model"

func BuildPolicySubgraph(p *model.Policy, config GraphvizConfig) (s subgraph) {
	s.Init(p.Title, config.Policy, 1)
	assumptionIDs := make(map[string]string)
	for _, assumption := range p.Assumptions {
		ss := BuildAssumptionSubgraph(assumption, config)
		var key string; for key = range assumption.PreConditions {break}
		assumptionIDs["cluster_" + generateID(ss.label)] = generateID(key)
		s.AddSubgraph(&ss)
	}
	for _, defense := range p.Defenses {
		// add nodes
		codeLines, defensePreConditions := GetDefenseCode(defense, config)
		for id, line := range codeLines {
			if _, present := s.Nodes[id]; !present {
				s.Nodes[id] = line
			}
		}
		s.preConditions = append(s.preConditions, defensePreConditions...)

		// add edges
		for id, nodeKey := range assumptionIDs { // connect defense to all assumptions under this policy
			connectAndAppend(s.Edges, nodeKey, generateID(defense.Title), "ltail=" + id)
		}
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
	for firstKey = range p.Defenses { break }
	for _, attacks := range p.Tags { // if attacks are linked to policies
		for _, attack := range attacks {
			connectAndAppend(s.Edges, generateID(attack.Title), generateID(firstKey), "lhead=cluster_" + generateID(p.Title))
		}
	}
	return
}