package graphviz

////////////////////////////////////////
// Structures used to specify graphviz properties

// Set of configuration parameters for entire graph.
type GraphvizConfig struct {
	Assumption				NodeProperties
	Policy						NodeProperties
	PreConditions			NodeProperties
	PreEmptiveDefense	NodeProperties
	IncidentResponse	NodeProperties
	Attack						NodeProperties
	
	EmptyDefense			NodeProperties
	EmptyAttack				NodeProperties

	Reality						NodeProperties
	AttackerWins			NodeProperties

	Subgraph					TextProperties
}

// Colors to use when generating graphviz code. These should be colors that your
// graphviz rendering engine supports. Values are not validated here.
type ColorSet struct {
	FontColor 	string
	FillColor 	string
	BorderColor	string
}
func (c ColorSet) Apply() (code string) {
	code += createProperty("fontcolor", c.FontColor, false)
	code += createProperty("fillcolor", c.FillColor, false)
	code += createProperty("color", c.BorderColor, false)
	return
}

// Font to use when generating graphviz code. This font  must be supported by your
// graphviz rendering engine. Values are not validated here.
type TextProperties struct {
	FontName	string
	FontSize	string
}
func (t TextProperties) Apply() (code string) {
	code += createProperty("fontname", t.FontName, false)
	code += createProperty("fontsize", t.FontSize, false)
	return
}

// Full set of properties for a specific node in the graph.
type NodeProperties struct {
	Color	ColorSet
	Font TextProperties
}
func (n NodeProperties) Apply(id string, label string, shape string) (code string) {
	code += id + "["
	code += createProperty("label", label, false)
	code += n.Font.Apply()
	code += n.Color.Apply()
	code += createProperty("shape", shape, false)
	code += createProperty("style", "filled, rounded", false)
	code += "]"
	return
}

////////////////////////////////////////
// Helper functions

// create a graphviz key-value pair (with or without quotes)
func createProperty(key string, value string, hasHTML bool) string {
	if !hasHTML { //Result: key="value"
		return " " + key + "=\"" + value + "\" "
	} else { //Result: key=value. This is required when 'value' has HTML style tags (B/U/I)
		return " " + key + "=" + value + " "
	}
}