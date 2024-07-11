package repl

import (
	"1ylang/evaluator"
	"1ylang/lexer"
	"1ylang/lib"
	"1ylang/object"
	"1ylang/parser"
	"bufio"
	"fmt"
	"io"
	"time"
)

func initEnv() *object.Environment {
	env := object.NewEnvironment()

	lib.RegisterStringFuncs(env)
	lib.RegisterArrayFuncs(env)
	lib.RegisterMathFuncs(env)

	return env
}

const PROMPT = ">> "

// Start starts the REPL
func Start(in io.Reader, out io.Writer, timed bool) {
	scanner := bufio.NewScanner(in)
	env := initEnv()

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()

		// Check for exit or quit commands
		if line == "exit" || line == "quit" {
			return
		}

		executeLine(out, line, env, timed)
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

// StartWithString executes a given input string
func StartWithString(out io.Writer, input string, timed bool) {
	env := initEnv()
	executeLine(out, input, env, timed)
}

// executeLine executes a single line of input and optionally times it
func executeLine(out io.Writer, line string, env *object.Environment, timed bool) {
	var startTime time.Time
	if timed {
		startTime = time.Now()
	}

	l := lexer.New(line)
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(out, p.Errors())
		return
	}

	evaluated := evaluator.Eval(program, env)
	if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}

	if timed {
		duration := time.Since(startTime)
		fmt.Fprintf(out, "Execution time: %v\n", duration)
	}
}
