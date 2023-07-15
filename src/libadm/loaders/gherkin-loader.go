package loaders

import (
	"strings"

	gherkin "github.com/cucumber/gherkin-go/v19"
	messages "github.com/cucumber/messages-go/v16"
)

func LoadGherkinContent(content string) (*messages.GherkinDocument, error) {
	builder := gherkin.NewAstBuilder((&messages.Incrementing{}).NewId)
	parser := gherkin.NewParser(builder)
	parser.StopAtFirstError(false)
	matcher := gherkin.NewLanguageMatcher(&admDialectProvider{}, "adm")
	err := parser.Parse(gherkin.NewScanner(strings.NewReader(content)), matcher)
	if err != nil {
		return nil, err
	}

	return builder.GetGherkinDocument(), nil
}