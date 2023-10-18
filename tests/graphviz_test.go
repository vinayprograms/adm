package test

import (
	"fmt"
	"libadm/graph"
	"libadm/graphviz"
	"libadm/loaders"
	"libadm/model"
	"os"
	"sources"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGrapvizhWithHumanOnlyModel(t *testing.T) {

	filepath := "./examples/friends.adm"
	outputpath := "./examples/friends.dot"

	file, err := os.ReadFile(filepath)
	assert.Nil(t, err)

	gherkin1, err := loaders.LoadGherkinContent(string(file))
	assert.Nil(t, err)

	var m1 model.Model
	err = m1.Init(gherkin1.Feature)
	assert.Nil(t, err)

	var g graph.Graph
	g.Init()
	err = g.AddModel(&m1)
	assert.Nil(t, err)

	code, err := graphviz.GenerateGraphvizCode(&g, getConfig())
	assert.Nil(t, err)

	output := strings.Join(code, "\n")
	err = os.WriteFile(outputpath, []byte(output), 0777)
	assert.Nil(t, err)
}

func TestGrapvizhModelWithBodylessSpecifications(t *testing.T) {
	filepath := "./examples/others/bodyless.adm"
	outputpath := "./examples/others/bodyless.dot"

	file, err := os.ReadFile(filepath)
	assert.Nil(t, err)

	gherkin1, err := loaders.LoadGherkinContent(string(file))
	assert.Nil(t, err)

	var m1 model.Model
	err = m1.Init(gherkin1.Feature)
	assert.Nil(t, err)

	var g graph.Graph
	g.Init()
	err = g.AddModel(&m1)
	assert.Nil(t, err)

	code, err := graphviz.GenerateGraphvizCode(&g, getConfig())
	assert.Nil(t, err)

	output := strings.Join(code, "\n")
	err = os.WriteFile(outputpath, []byte(output), 0777)
	assert.Nil(t, err)
}

func TestGrapvizhModelWithLotsOfStatements(t *testing.T) {
	filepath := "./examples/others/lengthy.adm"
	outputpath := "./examples/others/lengthy.dot"

	file, err := os.ReadFile(filepath)
	assert.Nil(t, err)

	gherkin1, err := loaders.LoadGherkinContent(string(file))
	assert.Nil(t, err)

	var m1 model.Model
	err = m1.Init(gherkin1.Feature)
	assert.Nil(t, err)

	var g graph.Graph
	g.Init()
	err = g.AddModel(&m1)
	assert.Nil(t, err)

	code, err := graphviz.GenerateGraphvizCode(&g, getConfig())
	assert.Nil(t, err)

	output := strings.Join(code, "\n")
	err = os.WriteFile(outputpath, []byte(output), 0777)
	assert.Nil(t, err)
}

func TestGrapvizhWithFullModel(t *testing.T) {
	filepath := "./examples/oauth/secrets-keys.adm"
	outputpath := "./examples/oauth/secrets-keys.dot"

	file, err := os.ReadFile(filepath)
	assert.Nil(t, err)

	gherkin1, err := loaders.LoadGherkinContent(string(file))
	assert.Nil(t, err)

	var m1 model.Model
	err = m1.Init(gherkin1.Feature)
	assert.Nil(t, err)

	var g graph.Graph
	g.Init()
	err = g.AddModel(&m1)
	assert.Nil(t, err)

	code, err := graphviz.GenerateGraphvizCode(&g, getConfig())
	assert.Nil(t, err)

	output := strings.Join(code, "\n")
	err = os.WriteFile(outputpath, []byte(output), 0777)
	assert.Nil(t, err)
}

func TestGrapvizhWithMultipleModels(t *testing.T) {
	dirpath := "./examples/oauth/"
	outputpath := "./examples/oauth/full-oauth.dot"

	models, err := getContentFromDirectory(sources.LocalSource{}, dirpath)
	assert.Nil(t, err)

	var g graph.Graph
	g.Init()
	for _, modelText := range models {
		gherkin1, err := loaders.LoadGherkinContent(modelText)
		assert.Nil(t, err)

		var m1 model.Model
		err = m1.Init(gherkin1.Feature)
		assert.Nil(t, err)

		err = g.AddModel(&m1)
		assert.Nil(t, err)
	}

	code, err := graphviz.GenerateGraphvizCode(&g, getConfig())
	assert.Nil(t, err)

	output := strings.Join(code, "\n")
	err = os.WriteFile(outputpath, []byte(output), 0777)
	assert.Nil(t, err)
}

func TestGrapvizhWithMultipleModelsLinkedViaPolicy(t *testing.T) {
	dirpath := "./examples/others/multifile"
	outputpath := "./examples/others/multifile.dot"

	models, err := getContentFromDirectory(sources.LocalSource{}, dirpath)
	assert.Nil(t, err)

	var g graph.Graph
	g.Init()
	for _, modelText := range models {
		gherkin1, err := loaders.LoadGherkinContent(modelText)
		assert.Nil(t, err)

		var m1 model.Model
		err = m1.Init(gherkin1.Feature)
		assert.Nil(t, err)

		err = g.AddModel(&m1)
		assert.Nil(t, err)
	}

	code, err := graphviz.GenerateGraphvizCode(&g, getConfig())
	assert.Nil(t, err)

	output := strings.Join(code, "\n")
	err = os.WriteFile(outputpath, []byte(output), 0777)
	assert.Nil(t, err)
}

////////////////////////////////////////
// Helper functions

func getContentFromDirectory(source sources.Source, path string) ([]string, error) {
	files, err := source.GetFiles(path)
	if err != nil {
		return nil, err
	}
	if len(files) > 0 {
		fmt.Println("Found", len(files), "file(s)")
	}
	var content []string
	for _, file := range files {
		newContent, err := source.GetContent(file)
		if err != nil {
			return content, err
		}

		content = append(content, newContent)
	}

	return content, nil
}

func getConfig() graphviz.GraphvizConfig {
	return graphviz.GraphvizConfig{
		Assumption: graphviz.NodeProperties{
			Color: graphviz.ColorSet{FontColor: "white", FillColor: "dimgray", BorderColor: "dimgray"},
			Font:  graphviz.TextProperties{FontName: "Times", FontSize: "18"},
		},
		Policy: graphviz.NodeProperties{
			Color: graphviz.ColorSet{FontColor: "black", FillColor: "darkolivegreen3", BorderColor: "darkolivegreen3"},
			Font:  graphviz.TextProperties{FontName: "Times", FontSize: "18"},
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

		// Start and end node config
		Reality: graphviz.NodeProperties{
			Color: graphviz.ColorSet{FontColor: "white", FillColor: "black", BorderColor: "black"},
			Font:  graphviz.TextProperties{FontName: "Arial", FontSize: "20"},
		},
		AttackerWins: graphviz.NodeProperties{
			Color: graphviz.ColorSet{FontColor: "red", FillColor: "yellow", BorderColor: "red"},
			Font:  graphviz.TextProperties{FontName: "Arial", FontSize: "20"},
		},

		Subgraph: graphviz.TextProperties{FontName: "Arial", FontSize: "24"},
	}
}
