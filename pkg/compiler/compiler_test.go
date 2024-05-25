package compiler

import (
	"fmt"
	"testing"

	"github.com/samasno/little-compiler/pkg/code"
	"github.com/samasno/little-compiler/pkg/frontend/ast"
	"github.com/samasno/little-compiler/pkg/frontend/lexer"
	"github.com/samasno/little-compiler/pkg/frontend/object"
	"github.com/samasno/little-compiler/pkg/frontend/parser"
)

type compilerTestCase struct {
	input                string
	expectedConstants    []interface{}
	expectedInstructions []code.Instructions
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             `1 + 2`,
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
		{
			input:             `1; 2`,
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
			},
		},
		{
			input:             `1 - 3`,
			expectedConstants: []interface{}{1, 3},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSub),
				code.Make(code.OpPop),
			},
		},
		{
			input:             `2 * 2`,
			expectedConstants: []interface{}{2, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpMul),
				code.Make(code.OpPop),
			},
		},
		{
			input:             `3 / 9`,
			expectedConstants: []interface{}{3, 9},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpDiv),
				code.Make(code.OpPop),
			},
		},
		{
			input:             `true;false;`,
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpPop),
				code.Make(code.OpFalse),
				code.Make(code.OpPop),
			},
		},
		{
			input:             `2 < 5`,
			expectedConstants: []interface{}{5, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGreaterThan),
				code.Make(code.OpPop),
			},
		},
		{
			input:             `1 != 100`,
			expectedConstants: []interface{}{1, 100},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpNotEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:             `2 == 2`,
			expectedConstants: []interface{}{2, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:             `true == true`,
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpTrue),
				code.Make(code.OpEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:             `true != false`,
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpFalse),
				code.Make(code.OpNotEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "-1",
			expectedConstants: []interface{}{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpMinus),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()
	for _, tt := range tests {
		program := parse(tt.input)
		compiler := New()

		err := compiler.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s\n", err)
		}

		bytecode := compiler.Bytecode()

		err = testInstructions(tt.expectedInstructions, bytecode.Instructions)
		if err != nil {
			t.Fatalf("testInstructions failed: %s\n", err)
		}

		err = testConstants(t, tt.expectedConstants, bytecode.Constants)
		if err != nil {
			t.Fatalf("testConstants failed: %s\n", err)
		}
	}
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func testInstructions(expected []code.Instructions, actual code.Instructions) error {
	concatted := concatInstructions(expected)
	if len(concatted) != len(actual) {
		return fmt.Errorf("wrong instructions length. \nwant %q\ngot %q", concatted.String(), actual.String())
	}

	for i, ins := range concatted {
		if actual[i] != ins {
			return fmt.Errorf("wrong instructions at %d\nwant %q\ngot %q\n", i, concatted.String(), actual.String())
		}
	}

	return nil
}

func testConstants(t *testing.T, expected []interface{}, actual []object.Object) error {
	t.Helper()
	if len(expected) != len(actual) {
		return fmt.Errorf("wrong number of constants\nwant %q\ngot %q", len(expected), len(actual))
	}

	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			err := testIntegerObject(int64(constant), actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed: %s", i, err)
			}
		}
	}

	return nil
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

func concatInstructions(s []code.Instructions) code.Instructions {
	out := code.Instructions{}

	for _, ins := range s {
		out = append(out, ins...)
	}

	return out
}
