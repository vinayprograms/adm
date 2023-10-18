package graphviz

import (
	"errors"
	"strings"
)

type spec struct {
	name       string
	properties string // each property is <key>=<value>. All properties are stored as a single string
}

// Graphviz lines structured into sections.
// Each member is a map connecting the specific graphviz line to the node it represents
type subgraph struct {
	label         string
	properties    NodeProperties    // list of properties to set for this specific subgraph.
	tabSpaceCount int               // for code beautification
	subgraphs     []*subgraph       // sub-subgraphs
	preConditions []string          // list of preconditions. Required for setting rank of all properties
	Nodes         map[string]string // map node IDs to node specific properties.
	Edges         map[string][]spec // maps a single source node ID to one or more destination node IDs.
}

func (s *subgraph) Init(title string, properties NodeProperties, tabcount int) error {
	if title == "" {
		return errors.New("graphviz subgraph title cannot be empty")
	}
	s.label = title
	s.properties = properties
	s.tabSpaceCount = tabcount
	if s.Nodes == nil {
		s.Nodes = make(map[string]string)
	}
	if s.Edges == nil {
		s.Edges = make(map[string][]spec)
	}

	return nil
}

func (s *subgraph) AddSubgraph(sub *subgraph) {
	sub.tabSpaceCount = s.tabSpaceCount + 1
	s.subgraphs = append(s.subgraphs, sub)
}

func (s *subgraph) RemoveSubgraph(sub *subgraph) {
	for i, ss := range s.subgraphs {
		if ss.label == sub.label {
			s.subgraphs = append(s.subgraphs[:i-1], s.subgraphs[i+1:]...)
		}
	}
}

func (s *subgraph) GenerateGraphvizNodes() (lines []string) {
	lines = appendLine(lines, s.tabSpaceCount, "subgraph cluster_"+generateID(s.label)+" {")
	lines = appendLine(lines, s.tabSpaceCount+1,
		strings.TrimSpace(createProperty("label", "<<B>"+htmlwrap(s.label)+"</B>>", true)))
	lines = appendLine(lines, s.tabSpaceCount+1,
		"graph[style=\"filled, rounded\" rankdir=\"LR\" splines=\"true\" overlap=\"false\" nodesep=\"0.2\" ranksep=\"0.9\""+
			s.properties.Font.Apply()+
			s.properties.Color.Apply()+
			"];")

	// Generate code for each sub-subgraph
	for _, ss := range s.subgraphs {
		sgLines := ss.GenerateGraphvizNodes()
		lines = append(lines, sgLines...)
	}

	// Add nodes
	for id, node := range s.Nodes {
		lines = appendLine(lines, s.tabSpaceCount+1, id+node)
	}
	sameRankSpec := "{rank=\"same\"; "
	for _, preCond := range s.preConditions {
		sameRankSpec += preCond + "; "
	}
	sameRankSpec += "}"
	lines = appendLine(lines, s.tabSpaceCount+1, sameRankSpec)
	lines = appendLine(lines, s.tabSpaceCount, "}")

	return
}

func (s *subgraph) GenerateGraphvizEdges(includeAttackerWins bool) (lines []string) {
	edges := s.GetUniqueEdges()
	for id, destinations := range edges {
		for _, dest := range destinations {
			if dest.properties != "" {
				lines = appendLine(lines, s.tabSpaceCount, id+" -> "+dest.name+"["+dest.properties+"]")
			} else {
				lines = appendLine(lines, s.tabSpaceCount, id+" -> "+dest.name)
			}
		}
	}

	if includeAttackerWins {
		// Add invisible edges from terminal nodes to attacker-wins.
		// This aligns attacker-wins for better readability.
		var terminalNodes []string
		if len(s.subgraphs) > 0 {
			for _, destinations := range edges {
				for _, destination := range destinations {
					if _, present := edges[destination.name]; !present { // if destination node has no successor, it is a terminal node.
						terminalNodes = append(terminalNodes, destination.name)
					}
				}
			}
			for _, node := range terminalNodes {
				lines = appendLine(lines, s.tabSpaceCount, node+" -> attacker_wins[style=\"invis\"]")
			}
		}
	}

	return
}

func (s *subgraph) GetUniqueEdges() (edges map[string][]spec) {
	edges = make(map[string][]spec)
	for _, ss := range s.subgraphs {
		for id, destinations := range ss.GetUniqueEdges() {
			if _, present := edges[id]; !present { // if this a new source -> dest
				edges[id] = destinations
			} else { // if source is already present, merge new destinations to current set
				for _, dest := range destinations {
					if !contains(dest, edges[id]) {
						edges[id] = append(edges[id], dest)
					}
				}
			}
		}
	}
	for id, destinations := range s.Edges {
		if _, present := edges[id]; !present { // if this a new source -> dest
			edges[id] = destinations
		} else { // if source is already present, merge new destinations to current set
			for _, dest := range destinations {
				if !contains(dest, edges[id]) {
					edges[id] = append(edges[id], dest)
				}
			}
		}
	}
	return
}
