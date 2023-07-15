package main

import (
	"args"
	"os"
)

func main() {
	args := args.Args{}
	err := args.ParseArgs(os.Args[1:])
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\r\n")
		os.Exit(1)
	}
}