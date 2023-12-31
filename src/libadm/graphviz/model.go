package graphviz

import (
	"libadm/model"
)

////////////////////////////////////////
// Exported functions

func BuildModelSubGraph(model *model.Model, unmitigatedAttacks map[string]interface{}, config GraphvizConfig) (modelSubgraph subgraph) {
	modelProperties := NodeProperties{
		Font:  TextProperties{FontName: "Arial", FontSize: "24"},
		Color: ColorSet{FontColor: "black", FillColor: "transparent", BorderColor: "gray"},
	}

	var allAssumptions []string
	for _, a := range model.Assumptions {
		allAssumptions = append(allAssumptions, CollectAssumptions(a)...)
	}

	title := buildTitleWithAssumptions(model.Title, allAssumptions, config)

	modelSubgraph.Init(generateID(model.Title), title, modelProperties, 1) // setup model subgraph

	//////////
	// generate subgraph for each policy
	policySubgraphs := getPolicySubgraphs(model.Policies, config)
	for _, asg := range policySubgraphs {
		modelSubgraph.AddSubgraph(&asg)
	}

	//////////
	// add attack nodes and edges to model subgraph
	attackNodes, attackEdges, attackPreConditions := getAttacks(model.Attacks, unmitigatedAttacks, config)
	for id, spec := range attackNodes {
		if _, present := modelSubgraph.Nodes[id]; !present {
			modelSubgraph.Nodes[id] = spec
		}
	}
	modelSubgraph.preConditions = append(modelSubgraph.preConditions, attackPreConditions...)
	for id, spec := range attackEdges {
		if _, present := modelSubgraph.Edges[id]; !present { // if this is a new source node
			modelSubgraph.Edges[id] = spec
		} else { // if source node exists, append only unique destinations
			for _, dest := range spec {
				if !contains(dest, modelSubgraph.Edges[id]) {
					modelSubgraph.Edges[id] = append(modelSubgraph.Edges[id], dest)
				}
			}
		}
	}

	//////////
	// add defense nodes and edges to model subgraph
	defenseNodes, defenseEdges, defensePreConditions := getDefenses(model.Defenses, config)
	for id, spec := range defenseNodes {
		if _, present := modelSubgraph.Nodes[id]; !present {
			modelSubgraph.Nodes[id] = spec
		}
	}
	modelSubgraph.preConditions = append(modelSubgraph.preConditions, defensePreConditions...)
	for id, spec := range defenseEdges {
		if _, present := modelSubgraph.Edges[id]; !present {
			modelSubgraph.Edges[id] = spec
		} else { // if source node exists, append only unique destinations
			for _, dest := range spec {
				if !contains(dest, modelSubgraph.Edges[id]) {
					modelSubgraph.Edges[id] = append(modelSubgraph.Edges[id], dest)
				}
			}
		}
	}

	return
}

////////////////////////////////////////
// Internal functions

func getPolicySubgraphs(policies map[string]*model.Policy, config GraphvizConfig) (subgraphs []subgraph) {
	for _, p := range policies {
		sg := BuildPolicySubgraph(p, config)
		_, edges, _ := getDefenses(p.Defenses, config)
		for key, value := range edges {
			sg.Edges[key] = value
		}
		sg.tabSpaceCount = 2 // Since the builder doesn't set it.
		subgraphs = append(subgraphs, sg)
	}
	return
}

func getAttacks(attacks map[string]*model.Attack, unmitigatedAttacks map[string]interface{}, config GraphvizConfig) (nodes map[string]string, edges map[string][]spec, preConditions []string) {
	nodes = make(map[string]string)
	edges = make(map[string][]spec)

	for _, a := range attacks {
		ac, attackPreConditions := GetAttackCode(a, unmitigatedAttacks, config)
		for id, properties := range ac {
			if _, present := nodes[id]; !present {
				nodes[id] = properties
			}
		}
		preConditions = append(preConditions, attackPreConditions...)

		if len(a.PreConditions) == 0 { // if attacks have no preconditions
			// does any action point to a defense
			doesThwartDefense := false
			for _, link := range a.Actions {
				if link.Item != nil {
					doesThwartDefense = true
					break
				}
			}
			// if attack doesn't succeed a defense and if there are no assumptions, connect attack to 'reality'
			if !doesThwartDefense {
				connectAndAppend(edges, "reality", generateID(a.Title), "")
			}
		}
		for _, pre := range a.PreConditions {
			if pre.Item == nil {
				connectAndAppend(edges, "reality", generateID(pre.Step.Statement), "")
				connectAndAppend(edges, generateID(pre.Step.Statement), generateID(a.Title), "")
			} else {
				for _, item := range pre.Item {
					connectAndAppend(edges, generateID(item.Title), generateID(a.Title), "")
				}
			}
		}
		for _, action := range a.Actions {
			if action.Item != nil {
				for _, item := range action.Item {
					connectAndAppend(edges, generateID(item.Title), generateID(a.Title), "")
				}
			}
		}
	}

	return
}

func getDefenses(defenses map[string]*model.Defense, config GraphvizConfig) (nodes map[string]string, edges map[string][]spec, preConditions []string) {
	nodes = make(map[string]string)
	edges = make(map[string][]spec)

	for _, d := range defenses {
		ac, defensePreConditions := GetDefenseCode(d, config)

		//////////
		// add nodes
		for id, properties := range ac {
			if _, present := nodes[id]; !present {
				nodes[id] = properties
			}
		}
		preConditions = append(preConditions, defensePreConditions...)

		////////////////////
		// add edges
		fomosecProperties := createProperty("label", "#fomosec", false)
		fomosecProperties += createProperty("penwidth", "2", false)
		fomosecProperties += createProperty("color", "red", false)
		fomosecProperties += createProperty("fontname", config.Subgraph.FontName, false)
		fomosecProperties += createProperty("fontcolor", "red", false)

		if len(d.PreConditions) == 0 && len(d.Actions) == 0 && len(d.Results) == 0 { // empty specifications (title only)
			connectAndAppend(edges, "reality", generateID(d.Title), "")
			continue
		}
		for _, pre := range d.PreConditions { // link defense precondition to defense node
			if pre.Item == nil { // if it is just a precondition clause
				if isFomosec(d) {
					connectAndAppend(edges, "reality", generateID(pre.Step.Statement), fomosecProperties)
				} else {
					connectAndAppend(edges, "reality", generateID(pre.Step.Statement), "")
				}
				connectAndAppend(edges, generateID(pre.Step.Statement), generateID(d.Title), "")
			} else { // if precondition points to another attack/defense
				for _, item := range pre.Item {
					connectAndAppend(edges, generateID(item.Title), generateID(d.Title), "")
				}
			}
		}
		for _, action := range d.Actions {
			if action.Item != nil { // if action points to another attack/defense
				for _, item := range action.Item {
					connectAndAppend(edges, generateID(item.Title), generateID(d.Title), "")
				}
			} else if len(d.PreConditions) == 0 && isFomosec(d) { // if action doesn't point to anything, and there are no preconditions and if its a #fomosec
				connectAndAppend(edges, "reality", generateID(d.Title), fomosecProperties)
			}
		}
		for _, attacks := range d.Tags { // if tag points to another attack/defense
			for _, attack := range attacks {
				connectAndAppend(edges, generateID(attack.Title), generateID(d.Title), "")
			}
		}
	}

	return
}

func isFomosec(d *model.Defense) bool {
	isfomosec := true
	for _, pre := range d.PreConditions {
		if pre.Item != nil {
			isfomosec = false // even if one precondition points to a defense/attack, its not #fomosec
			break
		}
	}
	for _, pre := range d.Actions {
		if pre.Item != nil {
			isfomosec = false // even if one action points to a defense/attack, its not #fomosec
			break
		}
	}
	for _, attacks := range d.Tags {
		if len(attacks) > 0 {
			isfomosec = false // even if one action points to a defense/attack, its not #fomosec
			break
		}
	}

	return isfomosec
}
