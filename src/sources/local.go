package sources

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
)

type LocalSource struct {
}

func (l LocalSource) CanHandle(path string) bool {
	// Can host can handle this path?
	_, err := os.Stat(path)
	return err == nil
}

func (l LocalSource) GetFiles(path string) ([]string, error) {
	if l.CanHandle(path) {
		return traverse(path)
	}

	return nil, nil
}

func (l LocalSource) GetContent(path string) (string, error) {

	var content string

	if l.CanHandle(path) {
		file, err := os.Open(path)
		if err != nil {
			return "", err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			content = content + "\n" + line
		}
	}

	return content, nil
}

func (l LocalSource) WriteContent(file string, data string) error {
	err := os.WriteFile(file, []byte(data), 0777)
	if err != nil {
		return err
	}
	return nil
}

/////////////////////////////////////////////////
// Helper functions

// Prepare a list of file paths
func traverse(path string) ([]string, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	var files []string

	if fileInfo.IsDir() {
		items, _ := ioutil.ReadDir(path)
		for _, item := range items {

			if item.IsDir() { // subdirectories
				newFiles, err := traverse(path + "/" + item.Name())
				if err != nil {
					return files, err
				}
				files = append(files, newFiles...)
			} else {
				// Only pick files with extension .adm
				fileparts := strings.Split(item.Name(), ".")
				fileExtension := fileparts[len(fileparts)-1]
				if fileExtension != "adm" {
					continue
				}
				files = append(files, path+"/"+item.Name())
			}
		}
	} else { // received a single file
		files = append(files, path)
	}

	return files, nil
}
