package args

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
)

// This is the client in 'Command' pattern implementation
type Args struct {
	statCmd			*flag.FlagSet
	graphCmd   		*flag.FlagSet
	exportCmd  		*flag.FlagSet
	path      		string
}

func (a *Args) InitArgs() {
	a.statCmd = flag.NewFlagSet("stat", flag.ExitOnError)
	a.statCmd.Bool("a", false, "List attacks")
	a.statCmd.Bool("d", false, "List all defenses")
	a.statCmd.Bool("p", false, "List defenses (pre-emtive only). Ignored if -d is present.")
	a.statCmd.Bool("i", false, "List defenses (incident-response only). Ignored if -d is present.")

	a.graphCmd = flag.NewFlagSet("graph", flag.ExitOnError)
	a.graphCmd.String("o", "./graphviz.dot", "Output path for graphviz(dot) file")

	a.exportCmd = flag.NewFlagSet("export", flag.ExitOnError)
	a.exportCmd.Bool("a", false, "Export attacks only")
	a.exportCmd.Bool("d", false, "Export defenses only")
	//a.exportCmd.String("f", "gherkin", "Output format: gherkin, deci(Deciduous)")
	a.exportCmd.String("o", "./", "Export path")
}

func (a *Args) PrintHelpToStdout() {

	if !a.isInitialized() {
		a.InitArgs()
	}

	fmt.Println("Usage: adm [OPTIONS] [PATH]")
	fmt.Println("\n[PATH]: The path to a directory or a single ADM file")
	fmt.Println("\nSub-command: stat")
	fmt.Println("List titles of all model components")
	a.statCmd.PrintDefaults()
	fmt.Println("\nSub-command: graph")
	fmt.Println("Generate a decision graph from one or more ADM files")
	a.graphCmd.PrintDefaults()
	fmt.Println("\nSub-command: export")
	fmt.Println("Export ADM files to other formats")
	a.exportCmd.PrintDefaults()
}

func (a *Args) ParseArgs(args []string) error {

	if !a.isInitialized() {
		a.InitArgs()
	}

	switch len(args) {
	case 0:
		a.PrintHelpToStdout()
		return nil
	case 1:
		return errors.New("require atleast two parameters - 'sub-command' and 'path'")
	}

	a.path = args[len(args)-1]	// Path must always be the last parameter

	// Process the remaining arguments.
	switch args[0] {
	case "stat":
		err := a.statCmd.Parse(args[1:])
		if err != nil {
			return err
		}
		aFlag,_ := strconv.ParseBool(a.statCmd.Lookup("a").Value.String())
		dFlag,_ := strconv.ParseBool(a.statCmd.Lookup("d").Value.String())
		pdFlag,_ := strconv.ParseBool(a.statCmd.Lookup("p").Value.String())
		irFlag,_ := strconv.ParseBool(a.statCmd.Lookup("i").Value.String())
		if !(aFlag || dFlag || pdFlag || irFlag) { // if all flags are skipped, assume -a and -d are true
			aFlag = true
			dFlag = true
		}
		return statsInvoker(aFlag, dFlag, pdFlag, irFlag, a.path)
	case "graph":
		err := a.graphCmd.Parse(args[1:])
		if err != nil {
			return err
		}
		return graphInvoker(a.graphCmd.Lookup("o").Value.String(), a.path)
		
	case "export":
		err := a.exportCmd.Parse(args[1:])
		if err != nil {
			return err
		}
		aFlag,_ := strconv.ParseBool(a.exportCmd.Lookup("a").Value.String())
		dFlag,_ := strconv.ParseBool(a.exportCmd.Lookup("d").Value.String())
		return exportInvoker(aFlag, dFlag, /*a.exportCmd.Lookup("f").Value.String(),*/ a.exportCmd.Lookup("o").Value.String(), a.path)

	default:
		return errors.New("INVALID ARGUMENT - \"" + args[0] + "\"")
	}
}

func (a *Args) isInitialized() bool {
	return (a.statCmd != nil)
}