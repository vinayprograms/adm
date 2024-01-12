module github.com/vinayprograms/adm/libadm

go 1.18

replace libadm/loaders => ./loaders

replace libadm/model => ./model

replace libadm/graph => ./graph

replace libadm/graphviz => ./graphviz

require (
	github.com/cucumber/gherkin-go/v19 v19.0.3
	github.com/cucumber/messages-go/v16 v16.0.1
)

require github.com/gofrs/uuid v4.3.0+incompatible // indirect
