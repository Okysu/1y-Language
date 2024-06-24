package main

import (
	"1ylang/repl"
	"flag"
	"fmt"
	"os"
)

const (
	VERSION = "0.1.0 (alpha-20240607)"
	HELP    = `Type "quit" or "exit" to exit.`
)

func main() {
	// Define command line flags
	filePath := flag.String("f", "", "Path to file to execute")
	timed := flag.Bool("t", false, "Enable timing of REPL commands")
	flag.Parse()

	if *filePath != "" {
		// If a file is provided with -f, run the script
		content, err := os.ReadFile(*filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", *filePath, err)
			os.Exit(1)
		}
		repl.StartWithString(os.Stdout, string(content), *timed)
	} else {
		// Otherwise, start the REPL
		fmt.Printf("1y Language %s -- %s\n", VERSION, "A programming language written in Go")
		fmt.Println(HELP)
		repl.Start(os.Stdin, os.Stdout, *timed)
	}
}
