package graphviz

import (
	"libadm/model"
	"strings"
)

////////////////////////////////////////
// Exported functions

func CollectAssumptions(a *model.Assumption) (assumptions []string) {
	for title := range a.PreConditions {
		assumptions = append(assumptions, title)
	}
	return
}

func buildTitleWithAssumptions(title string, assumptions []string, config GraphvizConfig) (titleString string) {
	var labelLines []string

	labelLines = appendLine(labelLines, 0, "<TABLE BORDER=\"0\" CELLBORDER=\"0\" CELLSPACING=\"0\">")
	labelLines = appendLine(labelLines, 3, "<TR><TD><FONT POINT-SIZE=\""+config.Subgraph.FontSize+"\"><B>"+htmlwrap(title)+"</B></FONT></TD></TR>")
	labelLines = appendLine(labelLines, 3, "<TR><TD></TD></TR>")
	if len(assumptions) > 0 {
		labelLines = appendLine(labelLines, 3, "<TR><TD><FONT POINT-SIZE=\""+config.Assumption.Font.FontSize+"\" COLOR=\""+config.Assumption.Color.FontColor+"\"><B>Assumptions</B></FONT></TD></TR>")
		labelLines = appendLine(labelLines, 3, "<TR><TD BORDER=\"1\" SIDES=\"T\" COLOR=\""+config.Assumption.Color.BorderColor+"\"></TD></TR>")
		for _, assumption := range assumptions {
			labelLines = appendLine(labelLines, 3, "<TR><TD ALIGN=\"LEFT\"><FONT POINT-SIZE=\""+config.Assumption.Font.FontSize+"\" COLOR=\""+config.Assumption.Color.FontColor+"\">• "+assumption+"</FONT></TD></TR>")
		}
		labelLines = appendLine(labelLines, 3, "<TR><TD BORDER=\"1\" SIDES=\"T\" COLOR=\""+config.Assumption.Color.BorderColor+"\"><BR/></TD></TR>")
	}
	labelLines = appendLine(labelLines, 0, "</TABLE>")

	titleString = strings.Join(labelLines, "\n")

	return
}
