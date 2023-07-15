module args

go 1.18

replace sources => ../sources
replace libadm/loaders => ../libadm/loaders
replace libadm/model => ../libadm/model
replace libadm/graph => ../libadm/graph
replace libadm/graphviz => ../libadm/graphviz
replace libadm/export => ../libadm/export


require (
	libadm v0.0.0-00010101000000-000000000000
	sources v0.0.0-00010101000000-000000000000
)

require (
	github.com/cucumber/gherkin-go/v19 v19.0.3 // indirect
	github.com/cucumber/messages-go/v16 v16.0.1 // indirect
	github.com/gofrs/uuid v4.3.0+incompatible // indirect
)
