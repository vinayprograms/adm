package test

import (
	"fmt"
	"libadm/loaders"
	"libadm/model"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

/////////////////////////////////////////////////
// Integration tests - Assumption + Step

func TestAssumption_NullInitialization(t *testing.T) {
	var c model.Assumption
	err := c.Init(nil)
	assert.Equal(t, err.Error(), "expected 'Assumption' spec. Got 'nil'")
}

func TestAssumptionWithTitleOnly(t *testing.T) {
	input := `
	Model: Model with just one rule
		Assumption: Some common assumption
	`

	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var c model.Assumption
	c.Init(gherkinModel.Feature.Children[0].Background)
	assert.Equal(t, c.Title, "Some common assumption")
}

func TestAssumptionWithWrongStatementType(t *testing.T) {
	input := `
	Model: Model with just one rule
		Assumption: Some common assumption
			When a wrong statement type is used in assumption
	`

	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var c model.Assumption
	err = c.Init(gherkinModel.Feature.Children[0].Background)
	assert.Equal(t, err.Error(), "Unexpected keyword - 'When' in Assumption specification")
}

func TestAssumptionSimpleAssumption(t *testing.T) {
	input := `
	Model: Model with just one assumption
		Assumption: Some common assumption
			Given a specific assumption information
	`

	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var c model.Assumption
	c.Init(gherkinModel.Feature.Children[0].Background)
	assert.Equal(t, c.Title, "Some common assumption")
	assert.NotNil(t, c.PreConditions["a specific assumption information"])
}

func TestAssumption_DuplicateStatements(t *testing.T) {
	input := `
	Model: Model with just one assumption
		Assumption: Some common assumption
			Given Part-1 of the assumption
			And Part-1 of the assumption
	`

	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var c model.Assumption
	err = c.Init(gherkinModel.Feature.Children[0].Background)
	assert.Equal(t, "precondition - 'Part-1 of the assumption' is already part of this assumption", err.Error())
}

func TestAssumption_Full(t *testing.T) {
	input := `
	Model: Model with just one assumption
		Assumption: Some common assumption
			Given Part-1 of the assumption
			And Part-2 of the assumption
	`

	gherkinModel, err := loaders.LoadGherkinContent(input)
	if err != nil {
		t.Error(err)
	}

	var c model.Assumption
	c.Init(gherkinModel.Feature.Children[0].Background)
	assert.Equal(t, c.Title, "Some common assumption")
	assert.NotNil(t, c.PreConditions["Part-1 of the assumption"])
	assert.NotNil(t, c.PreConditions["Part-2 of the assumption"])
}

func TestParseArgs(t *testing.T) {
	testVectors := map[string][]string{ // The last item in args list is the expected value
		"EmptyFile":                   {"stat", "./examples/basic/emptyfile.adm", "Found 1 file(s)"},
		"OneFile":                     {"stat", "./examples/oauth/access-tokens.adm", "Found 1 file(s)"},
		"MultiFile":                   {"stat", "./examples/oauth", "Found 7 file(s)"},
		"PathHierarchy":               {"stat", "./examples/subdirs", "Found 4 file(s)"},
		"PathToADMFile":               {"stat", "./examples/oauth/access-tokens.adm", "Found 1 file(s)"},
		"NonADMFile":                  {"stat", "./examples/noADMFile", ""},
		"PathAndAttackFlag":           {"stat", "-a", "./examples/oauth/access-tokens.adm", "Found 1 file(s)"},
		"PathAndDefenseFlag":          {"stat", "-d", "./examples/oauth/access-tokens.adm", "Found 1 file(s)"},
		"PathAndMultipleAttackFlags":  {"stat", "-a", "-a", "./examples/oauth/access-tokens.adm", "Found 1 file(s)"},
	}

	for name, args := range testVectors {
		t.Run(name, func(t *testing.T) {

			expected := args[len(args)-1]
			params := args[:len(args)-1]

			harness := output_interceptor{}
			harness.Hook()
			err := sendToParseArgs(params)
			out, _ := harness.ReadAndRelease()

			// only the last line contains the required output
			out = strings.Split(out, "\n")[1]
			fmt.Println(out)

			if err != nil {	// We don't expect explicit errors in these set of tests
				t.Error("ERROR:", err)
			}

			assert.Equal(t, expected, out)
		})
	}
}

/////////////////////////////////////////////////
// Helper functions
/*func sendToParseArgs(params []string) (err error) {
	parser := args.Args{}
	parser.InitArgs()

	fmt.Println("INPUT: ", os.Args)

	return parser.ParseArgs(params)
}*/

//*/*/*/*/*/*/*/*/*/*/*/*/*/*/*/*/*/*/*/
// TEST HARNESS
// Interceptor captures STDOUT and STDERR
// so that we can see fmt.Println's output and
// output to os.Stderr.
// - Your code must run after call to `Hook()`.
// - You should call ReadAndRelease() to restore
//   OS defaults.
/*type output_interceptor struct {
	originalOut *os.File
	originalErr *os.File

	interceptRead_Out  *os.File
	interceptWrite_Out *os.File

	interceptRead_Err  *os.File
	interceptWrite_Err *os.File
}

func (h *output_interceptor) Hook() {
	h.originalOut = os.Stdout
	h.originalErr = os.Stderr

	h.interceptRead_Out, h.interceptWrite_Out, _ = os.Pipe()
	h.interceptRead_Err, h.interceptWrite_Err, _ = os.Pipe()

	os.Stdout = h.interceptWrite_Out
	os.Stderr = h.interceptWrite_Out
}

func (h *output_interceptor) ReadAndRelease() (string, string) {
	h.interceptWrite_Out.Close()
	h.interceptWrite_Err.Close()

	out, _ := ioutil.ReadAll(h.interceptRead_Out)
	err, _ := ioutil.ReadAll(h.interceptRead_Err)

	os.Stdout = h.originalOut
	os.Stderr = h.originalErr

	h.interceptRead_Out.Close()
	h.interceptRead_Err.Close()
	return string(out), string(err)
}
*/