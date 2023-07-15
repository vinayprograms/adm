package graphviz

import (
	"libadm/model"
	"strings"
)

////////////////////////////////////////
// Exported functions

func BuildAssumptionSubgraph(a *model.Assumption, config GraphvizConfig) (s subgraph) {
	s.Init(a.Title, config.Assumption, 1)
	for title, p := range a.PreConditions {
		// add nodes
		s.preConditions = append(s.preConditions, generateID(title))
		code := generateNodeSpec(p.Statement, config.PreConditions, false)
		s.Nodes[generateID(title)] = "[" + strings.TrimSpace(code) + "]"
	}
	// add edge connecting reality to the assumption cluster.
	var firstKey string
	for firstKey = range a.PreConditions {break}
	connectAndAppend(s.Edges, "reality", generateID(a.PreConditions[firstKey].Statement), "lhead=cluster_" + generateID(a.Title))
	return
}