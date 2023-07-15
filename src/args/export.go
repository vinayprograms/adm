package args

import (
	"fmt"
	"libadm/export"
	"libadm/model"
	"os"
	"sources"
	"strings"
)

type gauntltCommand struct {
	models map[string]*model.Model // map file names to models.
	contentSource sources.Source
	path string
}

type gherkinCommand struct {
	models map[string]*model.Model // map file names to models.
	contentSource sources.Source
	path string
}

type deciduousCommand struct {
	models map[string]*model.Model // map file names to models.
	contentSource sources.Source
	includeAttacks bool
	includeDefenses bool
	path string
}

func (g gauntltCommand) execute() error {
	g.path = checkAndCreateDirectory(g.path)
	g.path += "export/attacks"
	g.path = checkAndCreateDirectory(g.path)

	for filename, model := range g.models {
		if model == nil { // no model to export
			fmt.Println("Failed loading model for " + filename + ". Please check file for errors.")
			continue 
		} 
		if len(model.Attacks) == 0 { // no attacks to export
			fmt.Println("No attacks found in " + filename)
			continue 
		}
		lines := export.ExportAttacks(model)
		write(g.contentSource, g.path, filename + ".attack", lines)
	}

	return nil
}

func (g gherkinCommand) execute() error {
	g.path = checkAndCreateDirectory(g.path)
	g.path += "export/defenses"
	g.path = checkAndCreateDirectory(g.path)

	for filename, model := range g.models {
		if model == nil { // no model to export
			fmt.Println("Failed loading model for " + filename + ". Please check file for errors.")
			continue 
		} 
		if len(model.Defenses) == 0 { // no attacks to export
			fmt.Println("No defenses found in " + filename)
			continue 
		}
		lines := export.ExportDefenses(model)
		write(g.contentSource, g.path, filename + ".feature", lines)
	}

	return nil
}

func (d deciduousCommand) execute() error {
	panic("not implemented")
}

////////////////////////////////////////
// Helper functions

func checkAndCreateDirectory(directory string) string {
	if directory[len(directory) - 1] != '/' { // append a "/" if directory string doesn't contain it.
		directory = directory + "/"
	}
	var path string
	for _, dir := range strings.Split(directory, "/") {
		if dir != "" {
			path += dir
			if _, err := os.Stat(path); err != nil {
				err = os.Mkdir(path, 0700) // create directory with RW rights to owner only.
				if err != nil {
					panic(err)
				}
			}
			path += "/"
		}
	}

	return directory
}

// Write lines to target filename
func write(target sources.Source, directory string, filename string, lines []string) {
	output := strings.Join(lines, "\n")
	target.WriteContent(directory + filename, output)
}