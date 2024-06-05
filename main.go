package main

import (
	"1ylang/repl"
	"bytes"
	"fmt"
	"os"
)

const (
	VERSION = "0.1.0 (alpha-20240605)"
	HELP    = `Type "quit" or "exit" to exit.`
)

func main() {
	if len(os.Args) > 1 {
		// If a file is provided as an argument, run the script
		filename := os.Args[1]
		content, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", filename, err)
			os.Exit(1)
		}
		repl.StartWithString(os.Stdout, string(bytes.TrimSpace(content)))
	} else {
		// Otherwise, start the REPL
		fmt.Printf("1y Language %s -- %s\n", VERSION, "A programming language written in Go")
		fmt.Println(HELP)
		repl.Start(os.Stdin, os.Stdout)
	}
}
