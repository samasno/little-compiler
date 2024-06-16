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

func TestStringExpression(t *testing.T) {
	tests := []vmTestCase{
		{`"monkey"`, "monkey"},
		{`"test " + "works"`, "test works"},
	}

	runVmTests(t, tests)
}

func TestArrayLiterals(t *testing.T) {
	tests := []vmTestCase{
		{`["test""]`, []string{"test"}},
		{`[]`, []int{}},
		{`[1+2, 4, 5]`, []int{3, 4, 5}},
	}

	runVmTests(t, tests)
}

func TestHashLiterals(t *testing.T) {
	tests := []vmTestCase{
		{`{}`, map[object.HashKey]int64{}},
		{`{1:2, 3:4}`, map[object.HashKey]int64{
			(&object.Integer{Value: 1}).HashKey(): 2,
			(&object.Integer{Value: 3}).HashKey(): 4},
		},
		{`{2 * 2: 10 / 2,}`, map[object.HashKey]int64{
			(&object.Integer{Value: 4}).HashKey(): 5},
		},
	}

	runVmTests(t, tests)
}

func TestIndexExpression(t *testing.T) {
	tests := []vmTestCase{
		{`[1,2,3][0]`, 1},
		{`[1,2,3][2-1]`, 2},
		{`[[1,2]][0][1]`, 2},
		{`[][0]`, Null},
		{`[1,2][100]`, Null},
		{`[1][-1]`, Null},
		{`{4:1}[4]`, 1},
		{`{1:2, 2:3}[2]`, 3},
		{`{}[0]`, Null},
		{`{1:3}[0]`, Null},
	}

	runVmTests(t, tests)
}

func TestFunctionsWithoutArguments(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `let fivePlusTen = fn(){5+10}();`, expected: 15,
		},
		{
			input: `
				let one = fn() { 1;};
				let two = fn() {2;};
				one() + two()
			`,
			expected: 3,
		},
		{
			input: `
				let a = fn() {1;};
				let b = fn() { a() + 1 };
				let c = fn() { return b() + 1 };
				c();
			`,
			expected: 3,
		},
		{
			input:    `let a = fn(){}; a()`,
			expected: Null,
		},
	}

	runVmTests(t, tests)
}

func TestFirstClassFunctions(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
				let returnsOne = fn() {1;};
				let returnsOneReturner = fn() { returnsOne;};
				returnsOneReturner()();	
			`,
			expected: 1,
		},
	}

	runVmTests(t, tests)
}

func TestFunctionCallWithLocalBindings(t *testing.T) {
	tests := []vmTestCase{
		{
			`let one = fn() { let one = 1; one; }();`,
			1,
		},
		{
			`let oneAndTwo = fn() { let one = 1; let two = 2; one + two; }();`,
			3,
		},
		{
			`
			let oneAndTwo = fn() { let one = 1; let two = 2; one + two };
			let threeAndFour = fn() { let three = 3; let four = 4; three + 4 };
			oneAndTwo() + threeAndFour()
			`,
			10,
		},
		{
			`
			let firstFoobar = fn() { let foobar = 50; foobar; };
			let secondFoobar = fn() { let foobar = 100; foobar;};
			firstFoobar() + secondFoobar();
			`,
			150,
		},
		{
			`
			let globalSeed = 50;

			let minusOne = fn() {
				let num = 1;
				globalSeed - num
			}

			let minusTwo = fn() {
				let num = 2
				globalSeed - num
			}

			minusOne() + minusTwo()
			`,
			97,
		},
	}

	runVmTests(t, tests)
}

func TestCallingFunctionsWithArgumentsAndBindings(t *testing.T) {
	tests := []vmTestCase{
		{
			`
			let identity = fn(a) { a; };
			identity(4)
			`,
			4,
		},
		{
			`
			let sum = fn(a,b) { a + b };
			sum(1,2);
			`,
			3,
		},
	}

	runVmTests(t, tests)
}

func TestCallingFunctionsWithWrongArguments(t *testing.T) {
	tests := []vmTestCase{
		{
			`
			fn(){1'}(1)
			`,
			`wrong number of arguments: want 0 got 1`,
		},
		{
			`fn(a) {a;}();`,
			`wrong number of arguments: want 1 got 0`,
		},
	}

	for _, tt := range tests {
		program := parse(tt.input)

		c := compiler.New()

		err := c.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		vm := New(c.Bytecode())
		err = vm.Run()
		if err == nil {
			t.Errorf("expected vm error but got none")
		}
		if err.Error() != tt.expected {
			t.Fatalf("wrong vm error: want %q got %q", tt.expected, err)
		}
	}
}

func testExpectedObject(t *testing.T, expected interface{}, actual object.Object) {
	t.Helper()
	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	case string:
		err := testStringObject(expected, actual)
		if err != nil {
			t.Errorf("testStringObject failed: %s", err)
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
	case []int:
		actual, ok := actual.(*object.Array)
		if !ok {
			t.Errorf("testArrayObject failed: not array got %T (+%v)", actual, actual)
		}

		if len(expected) != len(actual.Elements) {
			t.Errorf("wrong length of array: want %d got %d", len(expected), len(actual.Elements))
		}

		for i, el := range expected {
			err := testIntegerObject(int64(el), actual.Elements[i])
			if err != nil {
				t.Errorf("testIntegerObject failed: %s", err)
			}
		}
	case map[object.HashKey]int64:
		hash, ok := actual.(*object.Hash)
		if !ok {
			t.Errorf("object not Hash, got %T (%+v)", actual, actual)
			return
		}

		if len(hash.Pairs) != len(expected) {
			t.Errorf("hash has wrong number of pairs. want %d got %d", len(expected), len(hash.Pairs))
			return
		}

		for expectedKey, expectedValue := range expected {
			pair, ok := hash.Pairs[expectedKey]
			if !ok {
				t.Errorf("hash does not contain value for key %v", expectedKey)
			}

			err := testIntegerObject(expectedValue, pair.Value)
			if err != nil {
				t.Errorf("testIntegerObject failed: %s", err)
			}
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

func testStringObject(expected string, actual object.Object) error {
	result, ok := actual.(*object.String)
	if !ok {
		return fmt.Errorf("object is not string, got %T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. want %s got %s", expected, result.Value)
	}

	return nil
}

type vmTestCase struct {
	input    string
	expected interface{}
}
