package object

import (
	"fmt"
	"math/big"
	"reflect"
)

type EnvValue struct {
	Value    Object
	ReadOnly bool // If true, the value cannot be changed
}

type Environment struct {
	store map[string]EnvValue
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]EnvValue)
	return &Environment{store: s, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Store() map[string]EnvValue {
	return e.store
}

func (e *Environment) Get(name string) (Object, bool, bool) {
	env, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}
	return env.Value, ok, env.ReadOnly
}

func (e *Environment) NewVar(name string, val Object) Object {
	if !isValidName(name) {
		return newError("invalid variable name '%s'", name)
	}
	if e.isExist(name) {
		return newError("cannot redeclare variable '%s'", name)
	}
	e.store[name] = EnvValue{Value: val, ReadOnly: false}
	return val
}

func (e *Environment) Set(name string, val Object) Object {
	env, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Set(name, val)
	}
	if env.ReadOnly {
		return newError("cannot assign to constant '%s'", name)
	}
	e.store[name] = EnvValue{Value: val, ReadOnly: false}
	return val
}

func (e *Environment) isExist(name string) bool {
	_, ok := e.store[name]
	return ok
}

func (e *Environment) NewConst(name string, val Object) Object {
	if !isValidName(name) {
		return newError("invalid variable name '%s'", name)
	}
	if e.isExist(name) {
		return newError("cannot redeclare constant '%s'", name)
	}
	e.store[name] = EnvValue{Value: val, ReadOnly: true}
	return val
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func isValidName(name string) bool {
	if len(name) == 0 {
		return false
	}
	if name[0] >= '0' && name[0] <= '9' {
		return false
	}
	return true
}

func RegisterFunctions(env *Environment, namespace string, funcs map[string]interface{}) *Hash {
	pairs := make(map[HashKey]HashPair)

	for name, fn := range funcs {
		builtin := &Builtin{
			Fn: createBuiltinFunction(fn),
		}
		key := &String{Value: name}
		pairs[key.HashKey()] = HashPair{Key: key, Value: builtin}
	}

	hash := &Hash{Pairs: pairs}

	if namespace != "" {
		env.NewConst(namespace, hash)
	}

	return hash
}

func createBuiltinFunction(fn interface{}) BuiltinFunction {
	return func(args ...Object) Object {
		fnValue := reflect.ValueOf(fn)
		if fnValue.Kind() != reflect.Func {
			return newError("provided value is not a function")
		}

		fnType := fnValue.Type()
		if len(args) != fnType.NumIn() {
			return newError("wrong number of arguments: expected %d, got %d", fnType.NumIn(), len(args))
		}

		in := make([]reflect.Value, len(args))
		for i, arg := range args {
			in[i] = convertToReflectValue(arg, fnType.In(i))
			if !in[i].IsValid() {
				return newError("unsupported argument type: %s", arg.Type())
			}
		}

		out := fnValue.Call(in)
		if len(out) != 1 {
			return newError("unexpected number of return values")
		}

		return convertFromReflectValue(out[0])
	}
}

func convertToReflectValue(arg Object, targetType reflect.Type) reflect.Value {
	switch v := arg.(type) {
	case *Integer:
		if targetType.Kind() == reflect.Float64 {
			return reflect.ValueOf(float64(v.Value.Int64()))
		}
		return reflect.ValueOf(v.Value)
	case *Float:
		return reflect.ValueOf(v.Value)
	case *String:
		return reflect.ValueOf(v.Value)
	case *Array:
		if targetType.Kind() == reflect.Slice {
			elements := make([]interface{}, len(v.Elements))
			for i, elem := range v.Elements {
				elements[i] = convertToReflectValue(elem, targetType.Elem()).Interface()
			}
			return reflect.ValueOf(elements)
		}
	default:
		return reflect.Value{}
	}

	return reflect.Value{}
}

func convertFromReflectValue(val reflect.Value) Object {
	switch val.Kind() {
	case reflect.Float64:
		return &Float{Value: big.NewFloat(val.Float())}
	case reflect.String:
		return &String{Value: val.String()}
	case reflect.Slice:
		elements := make([]Object, val.Len())
		for i := 0; i < val.Len(); i++ {
			elements[i] = convertFromReflectValue(val.Index(i))
		}
		return &Array{Elements: elements}
	default:
		return newError("unsupported return type")
	}
}
