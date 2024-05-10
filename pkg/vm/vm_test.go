package vm

import (
	"fmt"
	"testing"

	"github.com/samasno/little-compiler/pkg/compiler"
	"github.com/samasno/little-compiler/pkg/frontend/ast"
	"github.com/samasno/little-compiler/pkg/frontend/lexer"
	"github.com/samasno/little-compiler/pkg/frontend/object"
	"github.com/samasno/little-compiler/pkg/frontend/parser"
)

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1+2", 3},
		{"7+7", 14},
		{"3*9", 27},
		{"3/1", 3},
		{"100-50", 50},
	}

	runVmTests(t, tests)
}

func testExpectedObject(t *testing.T, expected interface{}, actual object.Object) {
	t.Helper()
	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	}
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		comp := compiler.New()

		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error:%s", err)
		}

		vm := New(comp.Bytecode())

		err = vm.Run()
		if err != nil {
			t.Fatalf("vm error: %s", err)
		}

		stackElem := vm.LastPoppedStackElement()

		testExpectedObject(t, tt.expected, stackElem)

	}
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got %T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got %d want %d", result.Value, expected)
	}

	return nil
}

type vmTestCase struct {
	input    string
	expected interface{}
}
