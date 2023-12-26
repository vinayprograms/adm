package graphviz

import (
	"libadm/graph"
)

////////////////////////////////////////
// Exported functions

func GenerateGraphvizCode(g *graph.Graph, config GraphvizConfig) ([]string, error) {
	// Generate graphviz code
	var lines []string

	lines = append(lines, generateHeader("top", config.Reality, config.AttackerWins)...)

	lines = append(lines, generateBody(g, config)...)

	lines = append(lines, generateFooter()...)

	return lines, nil
}

////////////////////////////////////////
// Internal functions

func generateHeader(id string, realityProperties NodeProperties, attackerWinsProperties NodeProperties) (header []string) {
	header = appendLine(header, 0, "digraph \""+id+"\" {")
	// digraph properties are indented by 1 tabspace
	header = appendLine(header, 1, "// Base Styling")
	header = appendLine(header, 1, "compound=true")
	header = appendLine(header, 1, "graph[style=\"filled, rounded\" rankdir=\"LR\" splines=\"true\" overlap=\"false\" nodesep=\"0.2\" ranksep=\"0.9\"];")
	header = appendLineSpacer(header)
	header = appendLine(header, 1, "// Start and end nodes")
	header = appendLine(header, 1, realityProperties.Apply("reality", "Reality", "box"))
	return
}

func generateBody(g *graph.Graph, config GraphvizConfig) (lines []string) {
	for _, m := range g.Models {
		gvizModel := BuildModelSubGraph(m, g.UnmitigatedAttacks, config)
		// generate code for the model.
		lines = append(lines, gvizModel.GenerateGraphvizNodes()...)
		lines = append(lines, gvizModel.GenerateGraphvizEdges()...)
	}

	return
}

func generateFooter() (footer []string) {
	footer = appendLine(footer, 1, "subgraph cluster_Legend {")
	footer = appendLine(footer, 2, "label=\"Legend\"")
	footer = appendLine(footer, 2, "graph[style=\"filled, rounded\" rankdir=\"LR\" fontsize=\"16\" splines=\"true\" overlap=\"false\" nodesep=\"0.1\" ranksep=\"0.2\" fontname=\"Courier\" fillcolor=\"lightyellow\" color=\"yellow\"];")
	footer = appendLine(footer, 2, "A[label=\"Pre-\\nCondition\"  shape=\"box\"  style=\"filled, rounded\"  margin=\"0.2\"  fontname=\"Arial\"  fontsize=\"12\"  fontcolor=\"black\"  fillcolor=\"lightgray\"  color=\"gray\"]")
	footer = appendLine(footer, 2, "B[label=\"Attack\"  shape=\"box\"  style=\"filled, rounded\"  margin=\"0.2\"  fontname=\"Arial\"  fontsize=\"12\"  fontcolor=\"white\"  fillcolor=\"red\"  color=\"red\"]")
	footer = appendLine(footer, 2, "C[label=\"Unmitigated\\nAttack\"  shape=\"box\"  style=\"filled, rounded\"  margin=\"0.2\"  fontname=\"Arial\"  fontsize=\"12\"  fontcolor=\"red\"  fillcolor=\"yellow\"  color=\"red\" penwidth=\"4\"]")
	footer = appendLine(footer, 2, "D[label=\"Pre-emptive\\nDefense\"  shape=\"box\"  style=\"filled, rounded\"  margin=\"0.2\"  fontname=\"Arial\"  fontsize=\"12\"  fontcolor=\"white\"  fillcolor=\"purple\"  color=\"blue\"]")
	footer = appendLine(footer, 2, "E[label=\"Incident\\nResponse\"  shape=\"box\"  style=\"filled, rounded\"  margin=\"0.2\"  fontname=\"Arial\"  fontsize=\"12\"  fontcolor=\"white\"  fillcolor=\"blue\"  color=\"blue\"]")
	footer = appendLine(footer, 2, "F[label=\"Policy\"  shape=\"box\"  style=\"filled, rounded\"  margin=\"0.2\"  fontname=\"Arial\"  fontsize=\"12\"  fontcolor=\"black\"  fillcolor=\"darkolivegreen3\"  color=\"darkolivegreen3\"]")
	footer = appendLine(footer, 1, "}")
	footer = appendLine(footer, 1, "A -> reality [style=\"invis\" ltail=\"cluster_Legend\"]")
	footer = appendLine(footer, 1, "B -> reality [style=\"invis\" ltail=\"cluster_Legend\"]")
	footer = appendLine(footer, 1, "C -> reality [style=\"invis\" ltail=\"cluster_Legend\"]")
	footer = appendLine(footer, 1, "D -> reality [style=\"invis\" ltail=\"cluster_Legend\"]")
	footer = appendLine(footer, 1, "E -> reality [style=\"invis\" ltail=\"cluster_Legend\"]")
	footer = appendLine(footer, 1, "F -> reality [style=\"invis\" ltail=\"cluster_Legend\"]")
	footer = appendLine(footer, 0, "}")
	return
}
