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

func TestConditionals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
				if (true) { 10 }; 3333;
			`,
			expectedConstants: []interface{}{10, 3333},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpJumpNotTruthy, 10),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpJump, 11),
				code.Make(code.OpNull),
				code.Make(code.OpPop),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
				if (true) { 10 } else { 20 }; 3333;
			`,
			expectedConstants: []interface{}{10, 20, 3333},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpJumpNotTruthy, 10),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpJump, 13),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
      let one = 1;
      let two = 2;
      `,
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSetGlobal, 1),
			},
		},
		{
			input: `
      let one = 1;
      one;
      `,
			expectedConstants: []interface{}{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
      let one = 1;
      let two = one;
      two;
      `,
			expectedConstants: []interface{}{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpSetGlobal, 1),
				code.Make(code.OpGetGlobal, 1),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestDefine(t *testing.T) {
	expected := map[string]Symbol{
		"a": {Name: "a", Scope: GlobalScope, Index: 0},
		"b": {Name: "b", Scope: GlobalScope, Index: 1},
	}

	global := NewSymbolTable()

	a := global.Define("a")
	if a != expected["a"] {
		t.Errorf("a expected %v got %v", expected["a"], a)
	}

	b := global.Define("b")
	if b != expected["b"] {
		t.Errorf("b expected %v got %v", expected["b"], b)
	}

}

func TestResolveGlobal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	expected := []Symbol{
		{Name: "a", Scope: GlobalScope, Index: 0},
		{Name: "b", Scope: GlobalScope, Index: 1},
	}

	for _, sym := range expected {
		result, ok := global.Resolve(sym.Name)
		if !ok {
			t.Errorf("could not resolve name %s", sym.Name)
			continue
		}

		if result != sym {
			t.Errorf("%s expected %+v got %+v", sym.Name, sym, result)
		}
	}
}

func TestStringExpression(t *testing.T) {
  tests := []compilerTestCase{
      {
        input: `"monkey"`,
        expectedConstants: []interface{}{"monkey"},
        expectedInstructions: []code.Instructions {
          code.Make(code.OpConstant, 0),
          code.Make(code.OpPop),
      },
    },
    {
      input: `"mon" + "key"`,
      expectedConstants: []interface{}{"mon", "key"},
      expectedInstructions: []code.Instructions {
        code.Make(code.OpConstant, 0),
        code.Make(code.OpConstant, 1),
        code.Make(code.OpAdd),
        code.Make(code.OpPop),
      },
    },
  }
  
  runCompilerTests(t, tests)
}

func TestArrayLiterals(t *testing.T) {
  tests := []compilerTestCase{
    {
      input: `[]`,
      expectedConstants: []interface{}{},
      expectedInstructions: []code.Instructions{
        code.Make(code.OpArray, 0),
        code.Make(code.OpPop),
      },
    },
    {
      input: `[1,2,3]`,
      expectedConstants: []interface{}{1,2,3},
      expectedInstructions: []code.Instructions{
        code.Make(code.OpConstant, 0),
        code.Make(code.OpConstant, 1),
        code.Make(code.OpConstant, 2),
        code.Make(code.OpArray, 3),
        code.Make(code.OpPop),
      },
    },
    {
      input: `[1+2, "test", 10 - 10]`,
      expectedConstants: []interface{}{1,2, "test",10,10},
      expectedInstructions: []code.Instructions{
        code.Make(code.OpConstant, 0),
        code.Make(code.OpConstant, 1),
        code.Make(code.OpAdd),
        code.Make(code.OpConstant, 2),
        code.Make(code.OpConstant, 3),
        code.Make(code.OpConstant, 4),
        code.Make(code.OpSub),
        code.Make(code.OpArray, 3),
        code.Make(code.OpPop),
      },
    },
  }

  runCompilerTests(t, tests[2:])
}

func TestHashLiterals(t *testing.T) {
  tests := []compilerTestCase {
    {
      input: `{}`,
      expectedConstants: []interface{}{},
      expectedInstructions: []code.Instructions{
        code.Make(code.OpHash, 0),
        code.Make(code.OpPop),
      },
    },
    {
      input: `{1:2, 3:4, 5:6}`,
      expectedConstants: []interface{}{1,2,3,4,5,6},
      expectedInstructions: []code.Instructions{
        code.Make(code.OpConstant, 0),
        code.Make(code.OpConstant, 1),
        code.Make(code.OpConstant, 2),
        code.Make(code.OpConstant, 3),
        code.Make(code.OpConstant, 4),
        code.Make(code.OpConstant, 5),
        code.Make(code.OpHash, 6),
        code.Make(code.OpPop),
      },
    },
    {
      input: `{1:2+3,4:5*6}`,
      expectedConstants: []interface{}{1,2,3,4,5,6},
      expectedInstructions: []code.Instructions{
        code.Make(code.OpConstant, 0),
        code.Make(code.OpConstant, 1),
        code.Make(code.OpConstant, 2),
        code.Make(code.OpAdd),
        code.Make(code.OpConstant, 3),
        code.Make(code.OpConstant, 4),
        code.Make(code.OpConstant, 5),
        code.Make(code.OpMul),
        code.Make(code.OpHash, 4),
        code.Make(code.OpPop),
      },    
    },
  }

  runCompilerTests(t, tests)
}

func TestIndexExpression(t *testing.T) {
  tests := []compilerTestCase{
    {
      input: `[1,2,3][0+1]`,
      expectedConstants: []interface{}{1,2,3,0,1},
      expectedInstructions:[]code.Instructions{
        code.Make(code.OpConstant, 0),
        code.Make(code.OpConstant, 1),
        code.Make(code.OpConstant, 2),
        code.Make(code.OpArray, 3),
        code.Make(code.OpConstant, 3),
        code.Make(code.OpConstant, 4),
        code.Make(code.OpAdd),
        code.Make(code.OpIndex),
        code.Make(code.OpPop),
      },
    },
    {
      input: `{1:2}[2-1]`,
      expectedConstants: []interface{}{1,2,2,1},
      expectedInstructions: []code.Instructions{
        code.Make(code.OpConstant, 0),
        code.Make(code.OpConstant, 1),
        code.Make(code.OpHash, 2),
        code.Make(code.OpConstant, 2),
        code.Make(code.OpConstant, 3),
        code.Make(code.OpSub),
        code.Make(code.OpIndex),
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
    case string:
      err := testStringObject(constant, actual[i])
      if err != nil {
        return fmt.Errorf("constant %v - testStringObject failed: %s",i, err)
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

func testStringObject(expected string, actual object.Object) error {
  result, ok := actual.(*object.String)
  if !ok {
    return fmt.Errorf("object is not String. got %T (%v)", actual, actual)
  }

  if result.Value != expected {
    return fmt.Errorf("object has wrong value. want %s got %s", expected, result.Value)
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
