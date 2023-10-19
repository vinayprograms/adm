package test

import (
	"args"
	"fmt"
	"io"
	"log"
	"os"
)

// ///////////////////////////////////////////////
// Helper functions
func sendToParseArgs(params []string) (err error) {
	parser := args.Args{}
	parser.InitArgs()

	fmt.Println("INPUT: ", os.Args) // Print args so that tests can capture it using the harness (see below for harness).

	return parser.ParseArgs(params)
}

// */*/*/*/*/*/*/*/*/*/*/*/*/*/*/*/*/*/*/
// TEST HARNESS
// Interceptor captures STDOUT and STDERR
// so that we can see fmt.Println's output and
// output to os.Stderr.
//   - Your code must run after call to `Hook()`.
//   - You should call ReadAndRelease() to restore
//     OS defaults.
type output_interceptor struct {
	originalOut *os.File
	originalErr *os.File

	interceptRead_Out  *os.File
	interceptWrite_Out *os.File

	interceptRead_Err  *os.File
	interceptWrite_Err *os.File

	interceptRead_Log  *os.File
	interceptWrite_Log *os.File
}

func (h *output_interceptor) Hook() {
	h.originalOut = os.Stdout
	h.originalErr = os.Stderr

	h.interceptRead_Out, h.interceptWrite_Out, _ = os.Pipe()
	h.interceptRead_Err, h.interceptWrite_Err, _ = os.Pipe()
	h.interceptRead_Log, h.interceptWrite_Log, _ = os.Pipe()

	os.Stdout = h.interceptWrite_Out
	os.Stderr = h.interceptWrite_Out
	log.SetOutput(h.interceptWrite_Log)
}

func (h *output_interceptor) ReadAndRelease() (string, string, string) {
	h.interceptWrite_Out.Close()
	h.interceptWrite_Err.Close()
	h.interceptWrite_Log.Close()

	out, _ := io.ReadAll(h.interceptRead_Out)
	err, _ := io.ReadAll(h.interceptRead_Err)
	log, _ := io.ReadAll(h.interceptRead_Log)

	os.Stdout = h.originalOut
	os.Stderr = h.originalErr
	os.Stderr = h.originalErr

	h.interceptRead_Out.Close()
	h.interceptRead_Err.Close()
	h.interceptRead_Log.Close()
	return string(out), string(err), string(log)
}
