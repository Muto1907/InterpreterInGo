package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/Muto1907/interpreterInGo/ast"
)

type ObjectType string

const (
	INTEGER_OBJ  = "INTEGER"
	BOOLEAN_OBJ  = "BOOLEAN"
	NULL_OBJ     = "NULL"
	RETURN_OBJ   = "RETURN_VALUE"
	ERROR_OBJ    = "ERROR"
	FUNCTION_OBJ = "FUNCTION"
	STRING_OBJ   = "STRING"
	BUILTIN_OBJ  = "BUILTIN"
	ARRAY_OBJ    = "ARRAY"
	HASH_OBJ     = "HASH"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Environment struct {
	state map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	st := make(map[string]Object)
	return &Environment{state: st, outer: nil}
}

func (env *Environment) Set(ident string, val Object) Object {
	env.state[ident] = val
	return val
}

func (env *Environment) Get(ident string) (Object, bool) {
	val, ok := env.state[ident]
	if !ok && env.outer != nil {
		val, ok = env.outer.Get(ident)
	}
	return val, ok
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

type Hashable interface {
	HashKey() HashKey
}
type HashKey struct {
	Type  ObjectType
	Value uint64
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (hash *Hash) Type() ObjectType {
	return HASH_OBJ
}

func (hash *Hash) Inspect() string {
	var output bytes.Buffer
	pairs := []string{}
	for _, pair := range hash.Pairs {
		pairs = append(pairs, pair.Key.Inspect()+": "+pair.Value.Inspect())
	}
	output.WriteString("{")
	output.WriteString(strings.Join(pairs, ", "))
	output.WriteString("}")
	return output.String()
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

type String struct {
	Value string
}

func (str *String) Type() ObjectType {
	return STRING_OBJ
}

func (str *String) Inspect() string {
	return str.Value
}

func (str *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(str.Value))
	return HashKey{Type: str.Type(), Value: h.Sum64()}
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) HashKey() HashKey {
	var val uint64 = 1
	if !b.Value {
		val = 0
	}
	return HashKey{Type: b.Type(), Value: val}
}

type NULL struct{}

func (n *NULL) Type() ObjectType {
	return NULL_OBJ
}

func (n *NULL) Inspect() string {
	return "null"
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType {
	return RETURN_OBJ
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

type Error struct {
	Message string
}

func (er *Error) Type() ObjectType {
	return ERROR_OBJ
}

func (er *Error) Inspect() string {
	return "ERROR: " + er.Message
}

type Function struct {
	Params []*ast.Identifier
	Body   *ast.BlockStatement
	Env    *Environment
}

func (fn *Function) Type() ObjectType {
	return FUNCTION_OBJ
}

func (fn *Function) Inspect() string {
	var output bytes.Buffer
	parameters := []string{}
	for _, param := range fn.Params {
		parameters = append(parameters, param.String())
	}

	output.WriteString("fnc(")
	output.WriteString(strings.Join(parameters, ", "))
	output.WriteString(") {\n")
	output.WriteString(fn.Body.String())
	output.WriteString("\n}")
	return output.String()
}

type BuiltInFunction func(args ...Object) Object

type BuiltIn struct {
	Fnc BuiltInFunction
}

func (bi *BuiltIn) Type() ObjectType { return BUILTIN_OBJ }
func (bi *BuiltIn) Inspect() string  { return "builtIn Function" }

type Array struct {
	Elements []Object
}

func (arr *Array) Type() ObjectType { return ARRAY_OBJ }
func (arr *Array) Inspect() string {
	var output bytes.Buffer

	elements := []string{}
	for _, element := range arr.Elements {
		elements = append(elements, element.Inspect())
	}
	output.WriteString("[")
	output.WriteString(strings.Join(elements, ", "))
	output.WriteString("]")
	return output.String()
}
