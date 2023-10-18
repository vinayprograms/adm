package sources

type Source interface {
	CanHandle(path string) bool                  // Check if the specific source can handle this path
	GetFiles(path string) ([]string, error)      // Translate path to a list of files
	GetContent(file string) (string, error)      // Extract content from file path
	WriteContent(file string, data string) error // Write content to file
}

// Factory function to return the right Source handler
func GetSource(path string) Source {
	var sources []Source
	sources = append(sources, LocalSource{})
	sources = append(sources, GitSource{})
	sources = append(sources, UrlSource{})

	for _, source := range sources {
		if source.CanHandle(path) {
			return source
		}
	}

	return nil
}
