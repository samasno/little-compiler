package vm

import (
	"github.com/samasno/little-compiler/pkg/code"
	"github.com/samasno/little-compiler/pkg/frontend/object"
)

type Frame struct {
	fn *object.CompiledFunction
	ip int
}

func NewFrame(fn *object.CompiledFunction) *Frame {
	return &Frame{fn: fn, ip: -1}
}

func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}
