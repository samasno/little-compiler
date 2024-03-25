package compiler

import (
  "github.com/samasno/little-compiler/pkg/frontend/ast"
  "github.com/samasno/little-compiler/pkg/frontend/object"
  "github.com/samasno/little-compiler/pkg/code"
)

type Compiler struct {
  instructions code.Instructions
  constants []object.Object
}

func New() *Compiler {
  return &Compiler {
    instructions:code.Instructions{},
    constants:[]object.Object{},
  }
}

func (c *Compiler) Compile(node ast.Node) error {
  return nil
}

func (c *Compiler) Bytecode() *Bytecode {
  return &Bytecode{
    Instructions: c.instructions,
    Constants: c.constants,
  } 
}

type Bytecode struct {
  Instructions code.Instructions
  Constants []object.Object
}
