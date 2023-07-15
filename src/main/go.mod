module main

go 1.18

replace args => ../args

replace sources => ../sources

replace libadm => ../libadm

require args v0.0.0-00010101000000-000000000000

require (
	github.com/cucumber/gherkin-go/v19 v19.0.3 // indirect
	github.com/cucumber/messages-go/v16 v16.0.1 // indirect
	github.com/gofrs/uuid v4.3.0+incompatible // indirect
	libadm v0.0.0-00010101000000-000000000000 // indirect
	sources v0.0.0-00010101000000-000000000000 // indirect
)
