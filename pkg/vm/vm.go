package vm

import (
	"github.com/samasno/little-compiler/pkg/code"
	"github.com/samasno/little-compiler/pkg/compiler"
	"github.com/samasno/little-compiler/pkg/frontend/object"
)


const StackSize = 2048

type VM struct {
  constants []object.Object
  byteCode code.Instructions
  stack []object.Object
  sp int
}

func New(bytecode *compiler.Bytecode) *VM {
  return &VM{
    byteCode:bytecode.Instructions,
    constants:bytecode.Constants,
    stack:make([]object.Object, StackSize),
    sp:0,
  }
}

func (vm *VM) StackTop() object.Object {
  if vm.sp == 0 {
    return nil
  }

  return vm.stack[vm.sp-1]
}
