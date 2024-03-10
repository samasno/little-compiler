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
		{"if(1 == 1){ return 2}", 2},
		{"if(false) { return 100} else { return 2+2 }", 4},
		{"9; 100; return 5; 7;6;", 5},
		{"if(true){ 2+2; return 3; 100;}", 3},
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
  env := object.NewEnvironment()
	return Eval(program, env)
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
		{"(2 + 1) == 4", false},
		{"3 * 3 == 3", false},
		{"(10 * 10) == 100", true},
		{"2 > 100", false},
		{"if (2==2){ return true }", true},
		{"if (!2){return false} else { return true }", true},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		testBoolObject(t, obj, tt.expected)
	}
}

func TestEvalBangOperator(t *testing.T) {
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

func TestIfReturnsNull(t *testing.T) {
	inputs := []string{
		"if (false) { return 80 }",
		"if (!3) { return 100 }",
	}

	for _, tt := range inputs {
		obj := testEval(tt)
		ok := testIsNull(t, obj)
		if !ok {
			fmt.Printf("failed test case for '%s'\n", tt)
		}
	}
}

func TestErrorHandling(t *testing.T) {
  tests := []struct {
    input string
    expectedMessage string
  } {
    { "5 + true", "type mismatch: INTEGER + BOOLEAN"},
    { "2 + false; 3;", "type mismatch: INTEGER + BOOLEAN"},
    { "-true", "unknown operator: -BOOLEAN"},
    { "5; true + false; 2;", "unknown operator: BOOLEAN + BOOLEAN"},
    { "if (2 > 1) { true + false }", "unknown operator: BOOLEAN + BOOLEAN"},
    { "if (2 > 1) { if (2 > 1) { return true + true }} else { return 1}", "unknown operator: BOOLEAN + BOOLEAN"},
  }

  for _, tt := range tests {
    obj := testEval(tt.input)
    err, ok := obj.(*object.Error)
    if !ok {
      t.Errorf("failed test case for %s expected ERROR obj got %s\n", tt.input, reflect.TypeOf(obj))
      continue
    }

    if err.Message != tt.expectedMessage {
      t.Errorf("expected message '%s' got '%s'\n", tt.expectedMessage, err.Message)
    }
  }
}

func TestLetStatements(t *testing.T) {
  tests := []struct {
    input string
    expected int64
  }{
    {"let a = 100;a;", 100},
   // { "let b = 2 + 2;b;", 4},
   // {"let c = 10 * 2;c;", 20},
  }

  for _, tt := range tests {
    obj := testEval(tt.input)
    res := testIntegerObject(t, obj, tt.expected)
    if !res {
      fmt.Printf("failed test case for %s\n", tt.input)
    }
  }
  
}

func unwrapReturn(obj object.Object) object.Object {
	r, ok := obj.(*object.Return)
	if !ok {
		return obj
	}

	return r.Value
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

func testIsNull(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("expected NULL got %v\n", reflect.TypeOf(obj))
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


