package sources

type GitSource struct {
	//branch string
	//commit string
}

func (g GitSource) CanHandle(path string) bool {
	return false
}

func (g GitSource) GetFiles(path string) ([]string, error) {
	panic("not implemented")
}

func (g GitSource) GetContent(path string) (string, error) {
	panic("not implemented")
}

func (g GitSource) WriteContent(file string, data string) error {
	panic("not implemented")
}