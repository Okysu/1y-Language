package evaluator

import (
	"1ylang/object"
	"fmt"
	"math/big"
	"strings"
)

// newBuiltin is a helper function to create a new builtin function object.
func newBuiltin(fn func(args ...object.Object) object.Object) *object.Builtin {
	return &object.Builtin{Fn: fn}
}

var builtins = map[string]*object.Builtin{
	"len": newBuiltin(func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1", len(args))
		}

		switch arg := args[0].(type) {
		case *object.String:
			return &object.Integer{Value: big.NewInt(int64(len(arg.Value)))}
		case *object.Array:
			return &object.Integer{Value: big.NewInt(int64(len(arg.Elements)))}
		default:
			return newError("argument to `len` not supported, got %s", args[0].Type())
		}
	}),
	"puts": newBuiltin(func(args ...object.Object) object.Object {
		for _, arg := range args {
			fmt.Println(arg.Inspect())
		}
		return NULL
	}),
	"first": newBuiltin(func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1", len(args))
		}
		if args[0].Type() != object.ARRAY_OBJ {
			return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
		}

		arr := args[0].(*object.Array)
		if len(arr.Elements) > 0 {
			return arr.Elements[0]
		}

		return NULL
	}),
	"last": newBuiltin(func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1", len(args))
		}
		if args[0].Type() != object.ARRAY_OBJ {
			return newError("argument to `last` must be ARRAY, got %s", args[0].Type())
		}

		arr := args[0].(*object.Array)
		length := len(arr.Elements)
		if length > 0 {
			return arr.Elements[length-1]
		}

		return NULL
	}),
	"rest": newBuiltin(func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1", len(args))
		}
		if args[0].Type() != object.ARRAY_OBJ {
			return newError("argument to `rest` must be ARRAY, got %s", args[0].Type())
		}

		arr := args[0].(*object.Array)
		length := len(arr.Elements)
		if length > 0 {
			newElements := make([]object.Object, length-1)
			copy(newElements, arr.Elements[1:length])
			return &object.Array{Elements: newElements}
		}

		return NULL
	}),
	"push": newBuiltin(func(args ...object.Object) object.Object {
		if len(args) != 2 {
			return newError("wrong number of arguments. got=%d, want=2", len(args))
		}
		if args[0].Type() != object.ARRAY_OBJ {
			return newError("argument to `push` must be ARRAY, got %s", args[0].Type())
		}

		arr := args[0].(*object.Array)

		// Modify the original array by appending the new element
		arr.Elements = append(arr.Elements, args[1])

		// Return the modified array
		return arr
	}),
	"pop": newBuiltin(func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1", len(args))
		}
		if args[0].Type() != object.ARRAY_OBJ {
			return newError("argument to `pop` must be ARRAY, got %s", args[0].Type())
		}

		arr := args[0].(*object.Array)
		length := len(arr.Elements)
		if length > 0 {
			popped := arr.Elements[length-1]
			arr.Elements = arr.Elements[:length-1]
			return popped
		}

		return NULL
	}),
	"concat": newBuiltin(func(args ...object.Object) object.Object {
		if len(args) < 2 {
			return newError("wrong number of arguments. got=%d, want=2+", len(args))
		}

		for _, arg := range args {
			if arg.Type() != object.ARRAY_OBJ {
				return newError("argument to `concat` must be ARRAY, got %s", arg.Type())
			}
		}

		newElements := []object.Object{}
		for _, arg := range args {
			newElements = append(newElements, arg.(*object.Array).Elements...)
		}

		return &object.Array{Elements: newElements}
	}),
	"input": newBuiltin(func(args ...object.Object) object.Object {
		if len(args) == 1 {
			prompt := args[0].Inspect()
			// Define a map of escape sequences to their actual characters
			escapeSequences := map[string]string{
				"\\n":  "\n",
				"\\t":  "\t",
				"\\\\": "\\",
				"\\\"": "\"",
				"\\'":  "'",
				"\\r":  "\r",
				"\\b":  "\b",
				"\\f":  "\f",
				"\\v":  "\v",
			}

			// Replace escape sequences with actual characters
			for seq, char := range escapeSequences {
				prompt = strings.ReplaceAll(prompt, seq, char)
			}

			// Print the prompt
			fmt.Print(prompt)
		} else if len(args) > 1 {
			return newError("wrong number of arguments. got=%d, want=0 or 1", len(args))
		}

		var input string
		fmt.Scanln(&input)
		return &object.String{Value: input}
	}),
	"int": newBuiltin(func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1", len(args))
		}

		switch arg := args[0].(type) {
		case *object.String:
			value, ok := new(big.Int).SetString(arg.Value, 10)
			if !ok {
				return newError("cannot convert %s to int", arg.Value)
			}
			return &object.Integer{Value: value}
		case *object.Float:
			value, ok := new(big.Int).SetString(arg.Value.String(), 10)
			if !ok {
				return newError("cannot convert %s to int", arg.Value.String())
			}
			return &object.Integer{Value: value}
		default:
			return newError("argument to `int` must be STRING or FLOAT, got %s", args[0].Type())
		}
	}),
	"float": newBuiltin(func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1", len(args))
		}

		switch arg := args[0].(type) {
		case *object.String:
			value, ok := new(big.Float).SetString(arg.Value)
			if !ok {
				return newError("cannot convert %s to float", arg.Value)
			}
			return &object.Float{Value: value}
		case *object.Integer:
			return &object.Float{Value: new(big.Float).SetInt(arg.Value)}
		default:
			return newError("argument to `float` must be STRING or INTEGER, got %s", args[0].Type())
		}
	}),
	"str": newBuiltin(func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1", len(args))
		}

		switch arg := args[0].(type) {
		case *object.Integer:
			return &object.String{Value: arg.Value.String()}
		case *object.Float:
			return &object.String{Value: arg.Value.Text('f', -1)}
		default:
			return newError("argument to `str` must be INTEGER or FLOAT, got %s", args[0].Type())
		}
	}),
	"type": newBuiltin(func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1", len(args))
		}

		return &object.String{Value: string(args[0].Type())}
	}),
}
