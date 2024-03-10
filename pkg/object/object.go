package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/samasno/little-compiler/pkg/ast"
)

func NewEnvironment() *Environment {
  return &Environment{
    store: map[string]Object{},
  }
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
  return &Environment{store:map[string]Object{}, outer: outer}
}

type Environment struct {
  store map[string]Object
  outer *Environment
}

func (e *Environment) Get(name string) (Object, bool) {
  obj, ok := e.store[name]

  if !ok {
    obj, ok = e.outer.Get(name)
  }

  return obj, ok
}

func(e *Environment) Set(name string, value Object) Object {
  e.store[name] = value
  return value
}

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

type Boolean struct {
	Value bool
}

type Return struct {
	Value Object
}

type Error struct {
  Message string 
}

type Function struct {
  Params []*ast.Identifier
  Body *ast.BlockStatement
  Env *Environment
}

type Null struct{}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return fmt.Sprintf("%v", nil) }

func (r *Return) Type() ObjectType { return RETURN_OBJ }
func (r *Return) Inspect() string  { return r.Value.Inspect() }

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string { return e.Message }

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
  var out bytes.Buffer
  params := []string{}
  for _, p := range f.Params {
    params = append(params, p.String())

  }

  out.WriteString("fn")
  out.WriteString("()")
  out.WriteString(strings.Join(params, ","))
  out.WriteString(") {\n")
  out.WriteString(f.Body.String())
  out.WriteString("\n}")

  return out.String()
  
}

type ObjectType string

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ    = "NULL"
	RETURN_OBJ  = "RETURN"
  ERROR_OBJ = "ERROR"
  FUNCTION_OBJ = "FUNCTION"
)
