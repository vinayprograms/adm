package graphviz

import "strings"

//////////////////////////////////////////////////////////////////////////////
// Helper functions

// create an edge from source to destination and add it to edges-map
func connectAndAppend(edges map[string][]spec, source string, destination string, additionalProperties string) {
	if _, present := edges[source]; present { // if the source and destination are already listed
		for _, dest := range edges[source] {
			if dest.name == destination {
				return
			}
		}
	}
	edges[source] = append(edges[source], spec{name: destination, properties: additionalProperties})
}

// Get graphviz code for a specific node (assumption/policy/pre-condition/attack/defense)
func generateNodeSpec(title string, properties NodeProperties, isEmpty bool) (code string) {
	code = createProperty("label", wrap(title), false)
	if isEmpty {
		code += createProperty("shape", "box3d", false)
		code += createProperty("style", "filled, dashed", false)
	} else {
		code += createProperty("shape", "box", false)
		code += createProperty("style", "filled, rounded", false)
	}
	code += createProperty("margin", "0.2", false)
	code += properties.Font.Apply()
	code += properties.Color.Apply()
	return
}

func appendLine(document []string, tabs int, line string) []string {
	return append(document, genrateTabs(tabs) + line)
}
func appendLineSpacer(document []string) []string {
	return append(document, "")
}

func genrateTabs(tabCount int) (result string) {
	for i := 0; i < tabCount; i++ {
		result += "  "
	}
	return
}

func generateID(s string) string {
	return strings.ReplaceAll(cleanup(s), " ", "_")
}

// Replace symbols for use in ID generator
func cleanup(str string) (cleanedStr string) {
	cleanedStr = str
	// Symbols to remove
	for _, s := range []string{".", "(", ")", "[", "]", "{", "}", "'", "`", "-", "+", "?", ",", ":"} {
		cleanedStr = strings.ReplaceAll(cleanedStr, s, "")
	}
	// Replace with alternate
	replacements := map[string]string{
		"<":  "lt",
		">":  "gt",
		"=":  "eq",
		"\"": "\\\"",
	}
	for k, v := range replacements {
		cleanedStr = strings.ReplaceAll(cleanedStr, k, v)
	}

	return cleanedStr
}

// Wrap long lines exceeding 15 chars, along word boundaries
func wrap(s string) (wrapString string) {
	temp := strings.Fields(s)
	length := 0
	for _, str := range temp {
		length += len(str)
		if length > 15 {
			wrapString += "\\n" + str
			length = 0
		} else {
			wrapString += " " + str
		}
	}
	return strings.TrimSpace(wrapString)
}

// Wrap long lines, but use HTML line breaks.
func htmlwrap(s string) (wrapString string) {
	temp := strings.Fields(s)
	length := 0
	for _, str := range temp {
		length += len(str)
		if length > 15 {
			wrapString += "<br></br>" + str
			length = 0
		} else {
			wrapString += " " + str
		}
	}
	return strings.TrimSpace(wrapString)
}

// Check array membership
func contains[T comparable](item T, array []T) bool {
	for _, x := range array {
		if item == x {
			return true
		}
	}

	return false
}