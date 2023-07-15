package loaders

import "github.com/cucumber/gherkin-go/v19"

////////////////////////////////////////
// Hack: Custom dialect spec to accomodate ADM

type admDialectProvider struct{}

func (adp admDialectProvider) GetDialect(language string) *gherkin.GherkinDialect {
	
	// We don't case about 'language' parameter passed to this function
	
	return &(gherkin.GherkinDialect{
		Language: "adm", 
		Name: "attack-defense modeling", 
		Native: "attack-defense modeling", 
		Keywords: map[string][]string{
			// Customized for ADM
			"feature":    {"Model"},
			"scenario":   {"Attack", "Defense"},
			"rule":       {"Policy"},
			"background": {"Assumption"},

			// Default (from Gherkin)
			"examples":   {"Examples"},
			"given":      {"Given"},
			"when":       {"When"},
			"then":       {"Then"},
			"and":        {"And"},
			"but":        {"But"},
		},
	})
}