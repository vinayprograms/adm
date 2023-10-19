package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/////////////////////////////////////////////////
// Unit tests

func TestExportWithFilePath(t *testing.T) {
	err := sendToParseArgs([]string{"export", "-o", "./examples", "examples/friends.adm"})
	assert.Nil(t, err)
}

func TestExportWithMissingFile(t *testing.T) {
	err := sendToParseArgs([]string{"export", "-o", "./examples/others/", "./examples/others/missing-file.adm"})
	assert.Equal(t, "cannot identify the source for path './examples/others/missing-file.adm'", err.Error())
}

func TestExportWithFaultyADMSpec(t *testing.T) { // This should print errors to STDOUT.
	harness := output_interceptor{}
	harness.Hook()
	err := sendToParseArgs([]string{"export", "-o", "./examples/", "./examples/basic/faulty-admspec.adm"})
	_, _, log := harness.ReadAndRelease()

	assert.Nil(t, err)
	assert.Contains(t, log, "Skipping processing of faulty-admspec.adm")
	assert.Contains(t, log, "Failed loading model for faulty-admspec.adm")
}

func TestCorrectStatementSequenceAfterExport(t *testing.T) { // Tests if primary statements (Given/When/Then) and their "and" statements are sequenced properly.
	err := sendToParseArgs([]string{"export", "-o", "./examples/others/", "./examples/others/lengthy.adm"})
	assert.Nil(t, err)
}

func TestExportWithDirectoryPath(t *testing.T) {
	err := sendToParseArgs([]string{"export", "-o", "./examples/oauth", "examples/oauth/"})
	assert.Nil(t, err)
}
