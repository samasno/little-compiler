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
		{"-10", -10},
	}

	runVmTests(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"true", true},
		{"false", false},
		{"true == true", true},
		{"false == false", true},
		{"false != true", true},
		{"true != false", true},
		{"10 < 100", true},
		{"100 < 10", false},
		{"100 > 10", true},
		{"10 > 100", false},
		{"10 != 100", true}, //break
		{"10 == 10", true},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"!true", false},
		{"!false", true},
		{"!(1 < 3)", false},
		{"!!true", true},
		{"!5", false},
    {"!(if(false){ 5;})", true},
	}

	runVmTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []vmTestCase{
		{"if (true) { 10} ", 10},
		{"if (true) { 10} else { 20 }", 10},
		{"if (false) { 10} else { 20 }", 20},
		{"if (1) { 10} ", 10},
		{"if (1 < 2) { 10} ", 10},
		{"if (1 < 2) { 10} else { 20 } ", 10},
		{"if (1 > 2) { 10} else { 20 }", 20},
		{"if (false) {1}", Null},
		{"if (1 > 2) {1}", Null},
	  {"if((if(true){10})) {10} else {20}", 10},
    {"if((if(false){10})) {10} else {20}", 20},
  }

	runVmTests(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
  tests := []vmTestCase{
    {`let one = 1; one;`, 1},
    {`let one = 1; let two = 2; one + two`, 3},
    {`let one = 1; let two = one + one; one + two;`, 3},
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

	case bool:
		err := testBooleanObject(expected, actual)
		if err != nil {
			t.Errorf("testBooleanObject failed: %s", err.Error())
		}
	case *object.Null:
		if actual != Null {
			t.Errorf("testBooleanObject failed: expected Null got %T (%+v)", actual, actual)
		}
	}
}

func testBooleanObject(expected bool, obj object.Object) error {
	v := obj.(*object.Boolean).Value
	if v != expected {
		return fmt.Errorf("expected %v got %v", expected, v)
	}
	return nil
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
		return fmt.Errorf("object has wrong value. want %d got %d", expected, result.Value)
	}

	return nil
}

type vmTestCase struct {
	input    string
	expected interface{}
}
