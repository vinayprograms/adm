package test

import (
	"args"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

/////////////////////////////////////////////////
// Unit tests


func TestShowHelp(t *testing.T) {

	parser := args.Args{}
	parser.InitArgs()

	harness := output_interceptor{}
	harness.Hook()

	parser.PrintHelpToStdout()

	out, _ := harness.ReadAndRelease()

	expected :=
		"Usage: adm [OPTIONS] [PATH]\n" +
			"\n[PATH]: The path to a directory or a single ADM file\n" +
			"\nSub-command: stat\n" +
			"List titles of all model components\n" +
			"  -a	List attacks\n" +
			"  -d	List all defenses\n" +
			"  -i	List defenses (incident-response only). Ignored if -d is present.\n" +
			"  -p	List defenses (pre-emtive only). Ignored if -d is present.\n" + 
			"\n" +
			"Sub-command: graph\n" +
			"Generate a decision graph from one or more ADM files\n" +
			"  -o string\n" +
			"    	Output path for graphviz(dot) file (default \"./graphviz.dot\")\n" +
			"\n" +
			"Sub-command: export\n" +
			"Export ADM files to other formats\n" +
			"  -a	Export attacks only\n" +
			"  -d	Export defenses only\n" +
			"  -o string\n" +
			"    	Export path (default \"./\")\n"

	assert.Equal(t, out, expected)
}

func TestUninitializedArgsStruct_ShowHelp(t *testing.T) {
	parser := args.Args{}

	harness := output_interceptor{}
	harness.Hook()
	parser.ParseArgs([]string{})
	out, _ := harness.ReadAndRelease()
	expected :=
	"Usage: adm [OPTIONS] [PATH]\n" +
	"\n[PATH]: The path to a directory or a single ADM file\n" +
	"\nSub-command: stat\n" +
	"List titles of all model components\n" +
	"  -a	List attacks\n" +
	"  -d	List all defenses\n" +
	"  -i	List defenses (incident-response only). Ignored if -d is present.\n" +
	"  -p	List defenses (pre-emtive only). Ignored if -d is present.\n" + 
	"\n" +
	"Sub-command: graph\n" +
	"Generate a decision graph from one or more ADM files\n" +
	"  -o string\n" +
	"    	Output path for graphviz(dot) file (default \"./graphviz.dot\")\n" +
	"\n" +
	"Sub-command: export\n" +
	"Export ADM files to other formats\n" +
	"  -a	Export attacks only\n" +
	"  -d	Export defenses only\n" +
	"  -o string\n" +
	"    	Export path (default \"./\")\n"

	assert.Equal(t, out, expected)
}

func TestParseArgsWithoutSubCommand(t *testing.T) {
	args := []string{"./some_dir_that_doesn't_exist"}
	err := sendToParseArgs(args)
	assert.Equal(t, "require atleast two parameters - 'sub-command' and 'path'", err.Error())
}

func TestParseArgs_MissingParams(t *testing.T) {
	testVectors := map[string][]string{ // The last item in args list is the expected value
		"StatsInvoke":					{"stat", "require atleast two parameters - 'sub-command' and 'path'"},
		"graphInvoke":					{"graph", "require atleast two parameters - 'sub-command' and 'path'"},
		"exportInvoke":					{"export", "require atleast two parameters - 'sub-command' and 'path'"},
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

			assert.Equal(t, expected, err.Error())
		})
	}
}

// In the flags package, when a flag expects a string but doesn't find it, leads to os.Exit(2).
// Hence The test below is commented for that reason.
// TODO: Find a better way to test os.Exit scenarios.
/*
func TestParseArgsWithoutPathAndOnlyFlags_graph(t *testing.T) {
	
	testVectors := map[string][]string{ // The last item in args list is the expected value
		"graphInvoke":								{"graph", "-f", "dot", "-p", "cannot identify the source for path '-p'"},
		"graphInvokeWithoutOptions":	{"graph", "-f", "-p", "cannot identify the source for path '-p'"},
		"exportInvoke":								{"export", "-a", "-d", "-f", "-p", "cannot identify the source for path '-p'", "require atleast two parameters - 'sub-command' and 'path'"},
	}
	if os.Getenv("SUBPROC") == "1" {
		sendToParseArgs([]string{"graph", "-f", "-p"})
    return
  }

	// Run the test in a subprocess
  cmd := exec.Command(os.Args[0], "-test.run=TestGetConfig")
  cmd.Env = append(os.Environ(), "SUBPROC=1")
  err := cmd.Run()

	e, ok := err.(*exec.ExitError)
	expectedErrorString := "exit status 1"
  assert.Equal(t, true, ok)
  assert.Equal(t, expectedErrorString, e.Error())
}*/

func TestParseArgsWithoutPathAndOnlyFlags_Stats(t *testing.T) {
	err := sendToParseArgs([]string{"stat", "-a", "-d"})
	assert.Equal(t, "cannot identify the source for path '-d'", err.Error())
}

func TestParseArgsWithPathAndWrongFlag(t *testing.T) {
	args := []string{"-z", "./examples/basic"}
	err := sendToParseArgs(args)
	assert.Equal(t, "INVALID ARGUMENT - \"-z\"", err.Error())
}

func TestParseArgsWithUnsupportedPath(t *testing.T) {
	testVectors := map[string][]string{ // The last item in args list is the expected value
		"Stats":					{"stat", "dummy://dummy.dummy/test.adm", "cannot identify the source for path 'dummy://dummy.dummy/test.adm'"},
		"Graph":					{"graph", "dummy://dummy.dummy/test.adm", "cannot identify the source for path 'dummy://dummy.dummy/test.adm'"},
		"Export":					{"export", "dummy://dummy.dummy/test.adm", "cannot identify the source for path 'dummy://dummy.dummy/test.adm'"},
	}

	for name, args := range testVectors {
		t.Run(name, func(t *testing.T) {

			expected := args[len(args)-1]
			params := args[:len(args)-1]
			err := sendToParseArgs(params)
			assert.Equal(t, expected, err.Error())
		})
	}
}

func TestParseArgsWithValidPath(t *testing.T) {
	testVectors := map[string][]string{ // The last item in args list is the expected value
		"Stats":					{"stat", "examples/friends.adm"},
		"Graph":					{"graph", "examples/friends.adm"},
		"Export":					{"export", "examples/friends.adm"},
		"ExportAttacks":	{"export", "-a", "examples/friends.adm"},
		"ExportDefenses":	{"export", "-d", "examples/friends.adm"},
		//"ExportDeciduous":{"export", "-f", "deci", "examples/friends.adm"},
	}

	for name, args := range testVectors {
		t.Run(name, func(t *testing.T) {
			err := sendToParseArgs(args)
			assert.Nil(t, err)
		})
	}
}
