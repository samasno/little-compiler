package vm

import (
	"fmt"

	"github.com/samasno/little-compiler/pkg/code"
	"github.com/samasno/little-compiler/pkg/compiler"
	"github.com/samasno/little-compiler/pkg/frontend/object"
)


const StackSize = 2048

type VM struct {
  constants []object.Object
  instructions code.Instructions
  stack []object.Object
  sp int
}

func New(bytecode *compiler.Bytecode) *VM {
  return &VM{
    instructions:bytecode.Instructions,
    constants:bytecode.Constants,
    stack:make([]object.Object, StackSize),
    sp:0,
  }
}

func (vm *VM) Run() error {
  for ip:= 0; ip < len(vm.instructions); ip++ {
    op := code.Opcode(vm.instructions[ip])

    switch(op) {
      case code.OpConstant:
        constIndex := code.ReadUint16(vm.instructions[ip+1:])
        ip+=2

        err := vm.push(vm.constants[constIndex])
        if err != nil {
          return err
        }

      case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
        vm.executeBinaryOperation(op)

      case code.OpPop:
        vm.pop()
    }
  }

  return nil
}

func (vm *VM) executeBinaryOperation(op code.Opcode) error {
  right := vm.pop()
  left := vm.pop()

  switch {
    case right.Type() == object.INTEGER_OBJ && left.Type() == object.INTEGER_OBJ:
      return vm.executeBinaryIntegerOperation(op, left, right)
    default:
      return fmt.Errorf("invalid object types for binary operation: %s & %s", left.Type(), right.Type())
  }
}

func (vm *VM) executeBinaryIntegerOperation(op code.Opcode, left, right object.Object) error {
  l,_ := left.(*object.Integer)
  r,_ := right.(*object.Integer)

  var result int64
  switch op {
  case code.OpAdd:
    result = l.Value + r.Value
  case code.OpSub:
    result = l.Value - r.Value
  case code.OpMul:
    result = l.Value * r.Value
  case code.OpDiv:
    result = l.Value / r.Value
  }

  return vm.push(&object.Integer{Value: result})
}

func (vm *VM) push(o object.Object) error {
  if vm.sp >= StackSize {
    return fmt.Errorf("stack overflow")
  }

  vm.stack[vm.sp] = o
  vm.sp++

  return nil
}

func (vm *VM) StackTop() object.Object {
  if vm.sp == 0 {
    return nil
  }

  return vm.stack[vm.sp-1]
}

func (vm *VM) LastPoppedStackElement() object.Object {
  return vm.stack[vm.sp]
}

func (vm *VM) pop() object.Object {
  o := vm.stack[vm.sp-1]
  vm.sp--

  return o
}
