package object

import "fmt"

type EnvValue struct {
	Value Object
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

func (e *Environment) Get(name string) (Object, bool, bool) {
	env, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}
	return env.Value, ok, env.ReadOnly
}

func (e *Environment) NewVar(name string, val Object) Object {
	if e.isExist(name) {
		return newError("cannot redeclare variable '%s'", name)
	}
	e.store[name] = EnvValue{Value: val, ReadOnly: false}
	return val
}

func (e *Environment) Set(name string, val Object) Object {
	env := e.store[name]
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
	if e.isExist(name) {
		return newError("cannot redeclare constant '%s'", name)
	}
	e.store[name] = EnvValue{Value: val, ReadOnly: true}
	return val
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}
