package object

import (
	"bytes"
	"fmt"
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

func NeweEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
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

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
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
