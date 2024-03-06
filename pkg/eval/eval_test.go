package eval

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/samasno/little-compiler/pkg/ast"
	"github.com/samasno/little-compiler/pkg/lexer"
	"github.com/samasno/little-compiler/pkg/object"
)

func TestEvalIntegerObject(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"4", 4},
		{"100", 100},
		{"-100", -100},
		{"-33", -33},
		{"5 + 5 + 5 - 10", 5},
		{"7 - 2 + 3", 8},
		{"30 / 10 * 2", 6},
		{"3 + 3 * 2", 9},
		{"2 * ( 2 + 3 )", 10},
		{"(5 - ( 2 * 1 )) + 5", 8},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		ok := testIntegerObject(t, evaluated, tt.expected)
		if !ok {
			fmt.Printf("failed test case for '%s'\n", tt.input)
		}
	}

}

func testEval(input string) object.Object {
	l := lexer.NewLexer(input)
	p := ast.New(l)
	program := p.ParseProgram()
	return Eval(program)
}

func TestEvalBoolObject(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"3 < 5", true},
		{"1 != 2", true},
		{"2 + 2 == 4", true},
		{"2 - 2 == 4", false},
		{"3 * 3 == 3", false},
		{"10 * 10 != 100", false},
		{"2 > 100", false},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		testBoolObject(t, obj, tt.expected)
	}
}

func TestEvalNotOperand(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`!true`, false},
		{`!!true`, true},
		{`!false`, true},
		{`!!false`, false},
		{`!1`, false},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		ok := testBoolObject(t, obj, tt.expected)
		if !ok {
			fmt.Printf("Failed test case '%s'\n", tt.input)
		}

	}
}

func testIntegerObject(t *testing.T, obj object.Object, exp int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("expected integer obj got %s\n", reflect.TypeOf(obj))
	}

	if result.Value != exp {
		t.Errorf("expected %d got %d\n", exp, result.Value)
		return false
	}

	return true
}

func testBoolObject(t *testing.T, obj object.Object, exp bool) bool {
	b, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("expected bool obj got %s\n", reflect.TypeOf(obj))
	}
	if b.Value != exp {
		t.Errorf("expected %t got %t\n", exp, b.Value)
		return false
	}

	return true
}
