package object

import (
	"1ylang/ast"
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
	FLOAT_OBJ        = "FLOAT"

	BREAK_OBJ    = "BREAK"
	CONTINUE_OBJ = "CONTINUE"
)

// Integer represents an integer object
type Integer struct {
	Value   int64
	hashKey HashKey // Cached HashKey
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

func (i *Integer) HashKey() HashKey {
	if i.hashKey == (HashKey{}) {
		i.hashKey = HashKey{Type: i.Type(), Value: uint64(i.Value)}
	}
	return i.hashKey
}

// Boolean represents a boolean object
type Boolean struct {
	Value   bool
	hashKey HashKey // Cached HashKey
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

func (b *Boolean) HashKey() HashKey {
	if b.hashKey == (HashKey{}) {
		var value uint64
		if b.Value {
			value = 1
		} else {
			value = 0
		}
		b.hashKey = HashKey{Type: b.Type(), Value: value}
	}
	return b.hashKey
}

// Null represents a null object
type Null struct{}

func (n *Null) Inspect() string {
	return "null"
}

func (n *Null) Type() ObjectType {
	return NULL_OBJ
}

// ReturnValue represents a return value object
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

func (rv *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}

// Error represents an error object
type Error struct {
	Message string
}

func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}

// Function represents a function object
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}

func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

// String represents a string object
type String struct {
	Value   string
	hashKey HashKey // Cached HashKey
}

func (s *String) Type() ObjectType {
	return STRING_OBJ
}

func (s *String) Inspect() string {
	return s.Value
}

func (s *String) HashKey() HashKey {
	if s.hashKey == (HashKey{}) {
		h := fnv.New64a()
		h.Write([]byte(s.Value))
		s.hashKey = HashKey{Type: s.Type(), Value: h.Sum64()}
	}
	return s.hashKey
}

// BuiltinFunction represents a built-in function object
type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType {
	return BUILTIN_OBJ
}

func (b *Builtin) Inspect() string {
	return "builtin function"
}

// Array represents an array object
type Array struct {
	Elements []Object
}

func (ao *Array) Type() ObjectType {
	return ARRAY_OBJ
}

func (ao *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range ao.Elements {
		elements = append(elements, el.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// HashKey represents a hash key
type HashKey struct {
	Type  ObjectType
	Value uint64
}

// HashPair represents a key-value pair in a hash
type HashPair struct {
	Key   Object
	Value Object
}

// Hash represents a hash object
type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType {
	return HASH_OBJ
}

func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

type Hashable interface {
	HashKey() HashKey
}

type Float struct {
	Value   float64
	hashKey HashKey // Cached HashKey
}

func (f *Float) Inspect() string {
	return fmt.Sprintf("%f", f.Value)
}

func (f *Float) Type() ObjectType {
	return FLOAT_OBJ
}

func (f *Float) HashKey() HashKey {
	if f.hashKey == (HashKey{}) {
		h := fnv.New64a()
		h.Write([]byte(fmt.Sprintf("%f", f.Value)))
		f.hashKey = HashKey{Type: f.Type(), Value: h.Sum64()}
	}
	return f.hashKey
}

type Break struct{}

func (b *Break) Type() ObjectType { return BREAK_OBJ }
func (b *Break) Inspect() string  { return "break" }

type Continue struct{}

func (c *Continue) Type() ObjectType { return CONTINUE_OBJ }
func (c *Continue) Inspect() string  { return "continue" }
