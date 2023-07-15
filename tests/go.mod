module tests

go 1.18

replace args => ../src/args

replace sources => ../src/sources

replace libadm => ../src/libadm

require (
	args v0.0.0-00010101000000-000000000000
	github.com/cucumber/messages-go/v16 v16.0.1
	github.com/stretchr/testify v1.8.0
	libadm v0.0.0-00010101000000-000000000000
	sources v0.0.0-00010101000000-000000000000
)

require (
	github.com/cucumber/gherkin-go/v19 v19.0.3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gofrs/uuid v4.3.0+incompatible // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
