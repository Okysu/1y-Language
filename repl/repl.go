package repl

import (
	"1ylang/evaluator"
	"1ylang/lexer"
	"1ylang/object"
	"1ylang/parser"
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func initEnv() *object.Environment {
	env := object.NewEnvironment()

	workingDir, err := os.Getwd()
	if err == nil {
		env.Set("WORKING_DIR", &object.String{Value: workingDir})
	}

	executablePath, err := os.Executable()
	if err == nil {
		installDir := filepath.Dir(executablePath)
		env.Set("INSTALL_DIR", &object.String{Value: installDir})
	}
	return env
}

const PROMPT = ">> "

// Start starts the REPL
func Start(in io.Reader, out io.Writer) {
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

		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

func StartWithString(out io.Writer, input string) {
	env := initEnv()

	l := lexer.New(input)
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
}

