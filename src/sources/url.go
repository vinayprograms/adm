package sources

type UrlSource struct {
}

func (u UrlSource) CanHandle(path string) bool {
	return false
}

func (u UrlSource) GetFiles(path string) ([]string, error) {
	panic("not implemented")
}

func (u UrlSource) GetContent(path string) (string, error) {
	panic("not implemented")
}

func (u UrlSource) WriteContent(file string, data string) error {
	panic("not implemented")
}