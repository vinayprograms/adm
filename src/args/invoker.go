package args

import (
	"errors"
	"fmt"
	"libadm/graph"
	"libadm/loaders"
	"libadm/model"
	"log"
	"sources"
	"strings"
)

func statsInvoker(attacksOnly bool, defensesOnly bool, preemtiveDefensesOnly bool, incResponsesOnly bool, path string) error {
	src := sources.GetSource(path)
	if src == nil {
		return errors.New("cannot identify the source for path '" + path + "'")
	}

	models, err := getContent(src, path)
	if err != nil {
		return err
	}

	for _, modelText := range models {
		gherkin, err := loaders.LoadGherkinContent(modelText)
		if err != nil {
			return err
		}
		var model model.Model
		model.Init(gherkin.Feature)

		fmt.Println("\nMODEL: " + model.Title)
		fmt.Println(generateSeparator(len(model.Title) + 7))

		if attacksOnly {
			err := attackSummaryCommand{model: &model}.execute()
			if err != nil {
				return err
			}
		}

		if defensesOnly {
			err := defenseSummaryCommand{model: &model}.execute()
			if err != nil {
				return err
			}
		}

		if preemtiveDefensesOnly && !defensesOnly {
			err := preemtiveDefenseSummaryCommand{model: &model}.execute()
			if err != nil {
				return err
			}
		}

		if incResponsesOnly && !defensesOnly {
			err := incidentResponseSummaryCommand{model: &model}.execute()
			if err != nil {
				return err
			}
		}

		fmt.Println(generateSeparator(len(model.Title) + 7))
	}

	return nil
}

func graphInvoker(outputPath string, path string) error {
	src := sources.GetSource(path)
	if src == nil {
		return errors.New("cannot identify the source for path '" + path + "'")
	}

	g, err := getAdmGraph(src, path)
	if err != nil {
		return err
	}

	cmd := graphvizCommand{admGraph: g, outputPath: outputPath, destination: sources.LocalSource{}}

	return cmd.execute()
}

func exportInvoker(attacksOnly bool, defensesOnly bool /*format string,*/, outputPath string, path string) error {
	src := sources.GetSource(path)
	if src == nil {
		return errors.New("cannot identify the source for path '" + path + "'")
	}

	admModels, err := getAdmModels(src, path)
	if err != nil {
		return err
	}

	//var cmd command
	//switch format {
	//case "gherkin":
	if attacksOnly {
		return gauntltCommand{contentSource: src, path: outputPath, models: admModels}.execute()
	} else if defensesOnly {
		return gherkinCommand{contentSource: src, path: outputPath, models: admModels}.execute()
	} else {
		err := gauntltCommand{contentSource: src, path: outputPath, models: admModels}.execute()
		if err != nil {
			return err
		}
		err = gherkinCommand{contentSource: src, path: outputPath, models: admModels}.execute()
		if err != nil {
			return err
		}
	}
	/*case "deci":
		cmd = deciduousCommand{
			contentSource: src,
			includeAttacks: attacksOnly,
			includeDefenses: defensesOnly,
			path: outputPath,
			models: admModels}
		return cmd.execute()
	default:
		return errors.New("unsupported format - " + format)
	}*/

	return nil
}

/////////////////////////////////////////////////
// Helper functions

func generateSeparator(length int) (separator string) {
	for i := 0; i < length; i++ {
		separator += "="
	}
	return
}

func getAdmGraph(source sources.Source, path string) (graph.Graph, error) {
	var graph graph.Graph
	graph.Init()

	models, err := getAdmModels(source, path)
	if err != nil {
		return graph, err
	}
	for _, model := range models {
		err = graph.AddModel(model)
		if err != nil {
			return graph, err
		}
	}

	return graph, err
}

func getAdmModels(source sources.Source, path string) (map[string]*model.Model, error) {
	models := make(map[string]*model.Model)
	modelContents, err := getContent(source, path)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}
	for fileName, modelText := range modelContents {
		gherkin, err := loaders.LoadGherkinContent(modelText)
		if err != nil {
			log.Print(err.Error())
			log.Print("Skipping processing of " + fileName)
			models[fileName] = nil
			continue
		}
		var model model.Model
		err = model.Init(gherkin.Feature)
		if err != nil {
			log.Print(err.Error())
			log.Print("Skipping processing of " + fileName)
			models[fileName] = nil
			continue

		}
		models[fileName] = &model
	}

	return models, nil
}

func getContent(source sources.Source, path string) (map[string]string, error) {
	fileAndContent := make(map[string]string)

	files, err := source.GetFiles(path)
	if err != nil {
		return nil, err
	}
	if len(files) > 0 {
		log.Printf("Found %d file(s)", len(files))
	}
	for _, file := range files {
		newContent, err := source.GetContent(file)
		if err != nil {
			return nil, err
		}
		fileAndContent[getFileName(file)] = newContent
	}

	return fileAndContent, nil
}

func getFileName(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}
