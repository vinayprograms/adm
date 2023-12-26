package args

import (
	"libadm/graph"
	"libadm/graphviz"
	"sources"
	"strings"
)

type graphvizCommand struct {
	admGraph    graph.Graph
	outputPath  string
	destination sources.Source
}

func (g graphvizCommand) execute() error {
	graphvizLines, err := graphviz.GenerateGraphvizCode(&g.admGraph, getGraphvizConfig())

	if err == nil {
		output := strings.Join(graphvizLines, "\n")
		g.destination.WriteContent(g.outputPath, output)
		return nil
	} else {
		return err
	}
}

// Prepare graphviz font and color specifications for different node types.
func getGraphvizConfig() graphviz.GraphvizConfig {
	return graphviz.GraphvizConfig{
		Assumption: graphviz.NodeProperties{
			Color: graphviz.ColorSet{FontColor: "brown", FillColor: "", BorderColor: "lightgray"},
			Font:  graphviz.TextProperties{FontName: "Arial", FontSize: "14"},
		},
		Policy: graphviz.NodeProperties{
			Color: graphviz.ColorSet{FontColor: "black", FillColor: "darkolivegreen3", BorderColor: "darkolivegreen3"},
			Font:  graphviz.TextProperties{FontName: "Arial", FontSize: "18"},
		},
		PreConditions: graphviz.NodeProperties{
			Color: graphviz.ColorSet{FontColor: "black", FillColor: "lightgray", BorderColor: "gray"},
			Font:  graphviz.TextProperties{FontName: "Arial", FontSize: "16"},
		},

		// Defense config
		PreEmptiveDefense: graphviz.NodeProperties{
			Color: graphviz.ColorSet{FontColor: "white", FillColor: "purple", BorderColor: "blue"},
			Font:  graphviz.TextProperties{FontName: "Arial", FontSize: "16"},
		},
		IncidentResponse: graphviz.NodeProperties{
			Color: graphviz.ColorSet{FontColor: "white", FillColor: "blue", BorderColor: "blue"},
			Font:  graphviz.TextProperties{FontName: "Arial", FontSize: "16"},
		},
		EmptyDefense: graphviz.NodeProperties{
			Color: graphviz.ColorSet{FontColor: "black", FillColor: "transparent", BorderColor: "blue"},
			Font:  graphviz.TextProperties{FontName: "Arial", FontSize: "16"},
		},

		// Attack config
		Attack: graphviz.NodeProperties{
			Color: graphviz.ColorSet{FontColor: "white", FillColor: "red", BorderColor: "red"},
			Font:  graphviz.TextProperties{FontName: "Arial", FontSize: "16"},
		},
		EmptyAttack: graphviz.NodeProperties{
			Color: graphviz.ColorSet{FontColor: "black", FillColor: "transparent", BorderColor: "red"},
			Font:  graphviz.TextProperties{FontName: "Arial", FontSize: "16"},
		},
		UnMitigatedAttack: graphviz.NodeProperties{
			Color: graphviz.ColorSet{FontColor: "red", FillColor: "yellow", BorderColor: "red"},
			Font:  graphviz.TextProperties{FontName: "Arial", FontSize: "16"},
		},

		// Start and end node config
		Reality: graphviz.NodeProperties{
			Color: graphviz.ColorSet{FontColor: "white", FillColor: "black", BorderColor: "black"},
			Font:  graphviz.TextProperties{FontName: "Arial", FontSize: "20"},
		},
		AttackerWins: graphviz.NodeProperties{
			Color: graphviz.ColorSet{FontColor: "red", FillColor: "yellow", BorderColor: "yellow"},
			Font:  graphviz.TextProperties{FontName: "Arial", FontSize: "20"},
		},

		Subgraph: graphviz.TextProperties{FontName: "Arial", FontSize: "24"},
	}
}
