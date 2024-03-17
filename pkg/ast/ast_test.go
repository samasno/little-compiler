package ast

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/samasno/little-compiler/pkg/lexer"
	"github.com/samasno/little-compiler/pkg/tokens"
)

func TestParseLetStatement(t *testing.T) {
	input := `
let a = 100
let b = 231;
let c = 25252525;
let d = 2 + 2;
let e = (5 * 5) + 2;
`

	l := lexer.NewLexer(input)
	p := New(l)

	program := p.ParseProgram()

	if program == nil {
		t.Fatal("parse program nil")
	}

	ids := []struct {
		expectedIdentifier string
	}{
		{"a"},
		{"b"},
		{"c"},
		{"d"},
		{"e"},
	}

	for i, s := range ids {
		raw := program.Statements[i]
		statement, ok := raw.(*LetStatement)
		println(raw.String())

		if !ok {
			t.Errorf("statement %d not a let statement\n", i+1)
		}

		lit := statement.TokenLiteral()

		if lit != tokens.LET {
			t.Errorf("mismatched token literal, expected %s got %s\n", tokens.LET, lit)
		}

		identifier := statement.Name.TokenLiteral()

		if identifier != s.expectedIdentifier {
			t.Errorf("mismatched identifiers, expected %s got %s\n", s.expectedIdentifier, identifier)
		} else {
			fmt.Printf("got identifier %s\n", s.expectedIdentifier)
		}
	}

}

func TestParseLetFn(t *testing.T) {
	input := `
let addTwo = fn(a){return fn(){return a + 2}}
`
	l := lexer.NewLexer(input)
	p := New(l)
	program := p.ParseProgram()

	if len(program.Statements) != 1 {
		t.Errorf("expected 1 statement got %d\n", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*LetStatement)
	if !ok {
		t.Errorf("expected exp stmt got %s\n", reflect.TypeOf(program.Statements[0]))
	}

	fn, ok := stmt.Value.(*FnLiteral)
	if !ok {
		t.Errorf("expected fn lit got %s\n", reflect.TypeOf(stmt.Value).String())
	}

	if len(fn.Params) != 1 {
		t.Errorf("expected 1 param got %d\n", len(fn.Params))
	}

	ident := fn.Params[0]

	if ident.Value != "a" {
		t.Errorf("expected identifier 'a' got %v\n", ident.Value)
	}

	if len(fn.Body.Statements) != 1 {
		t.Errorf("expected on statement in fn lit body")
	}

	rs, ok := fn.Body.Statements[0].(*ReturnStatement)
	if !ok {
		t.Errorf("expected return statement got %s\n", reflect.TypeOf(fn.Body.Statements[0]).String())
	}

	rfn, ok := rs.Value.(*FnLiteral)
	if !ok {
		t.Errorf("expected return value to be a fn lit got %s\n", reflect.TypeOf(rs.Value))
	}

	if len(rfn.Params) != 0 {
		t.Errorf("expected returned fn to have 0 params got %d\n", len(rfn.Params))
	}

	if len(rfn.Body.Statements) != 1 {
		t.Errorf("expected on statement in return fn got %d\n", rfn.Body.Statements)
	}

	drs, ok := rfn.Body.Statements[0].(*ReturnStatement)
	if !ok {
		t.Errorf("expected return return statemetn got %s\n", reflect.TypeOf(rfn.Body.Statements[0]))
	}

	difx, ok := drs.Value.(*InfixExpression)
	if !ok {
		t.Errorf("final infix exp got %s\n", reflect.TypeOf(drs.Value).String())
	}

	println(difx.String())

}

func TestParseReturnStatement(t *testing.T) {
	input := `
return 100;
return 200;
return 1;
`
	l := lexer.NewLexer(input)
	p := New(l)

	program := p.ParseProgram()

	ex := []struct {
		Value string
	}{
		{"100"},
		{"200"},
		{"1"},
	}

	if len(program.Statements) != len(ex) {
		t.Errorf("expected %d statements got %d", len(ex), len(program.Statements))
	}

	for i, s := range program.Statements {
		rs, ok := s.(*ReturnStatement)
		if !ok {
			t.Errorf("statement %d not a return statement\n", i)
		}

		if rs.Token.Literal != tokens.RETURN {
			t.Errorf("expected token literal %s but got %s\n", tokens.RETURN, rs.TokenLiteral())
		}
		println(i, s.String())
	}

	println(program.String())
}

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: tokens.Token{
					Type: tokens.LET, Literal: "let",
				},
				Name: &Identifier{
					Token: tokens.Token{Type: tokens.IDENTIFIER, Literal: "a"},
					Value: "a",
				},
				Value: &Identifier{
					Token: tokens.Token{Type: tokens.IDENTIFIER, Literal: "b"},
					Value: "b",
				},
			},
		},
	}

	str := `let a = b`
	pstr := program.String()

	if pstr != str {
		t.Fatalf("expected '%s' got '%s'", str, pstr)
	}
}

func TestParseStrings(t *testing.T) {
	input := `
let a = "one hundred"
let b = "test"
`
	expected := []string{
		`one hundred`,
		`test`,
	}

	l := lexer.NewLexer(input)
	p := New(l)
	program := p.ParseProgram()

	for i, ll := range program.Statements {
		ls, ok := ll.(*LetStatement)
		if !ok {
			t.Fatalf("expected let statement got %s\n", reflect.TypeOf(ll))
		}

		str, ok := ls.Value.(*StringLiteral)
		if !ok {
			t.Fatalf("expected str literal got %s\n", reflect.TypeOf(ls.Value))
		}

		if str.Value != expected[i] {
			t.Fatalf("expected '%s' got '%s'\n", expected[i], str.Value)
		}

	}
}

func TestParseStringsInfix(t *testing.T) {
	input := `
let a = "string" + " works"
`

	l := lexer.NewLexer(input)
	p := New(l)
	pr := p.ParseProgram()

	if len(pr.Statements) != 1 {
		t.Fatalf("exected 1 statement got %d\n", len(pr.Statements))

	}

	ls, ok := pr.Statements[0].(*LetStatement)
	if !ok {
		t.Fatalf("didn't get let statement got %s\n", reflect.TypeOf(pr.Statements[0]))
	}

	ifx, ok := ls.Value.(*InfixExpression)
	if !ok {
		t.Fatalf("expected infix got %s\n", reflect.TypeOf(ls.Value))
	}

	left, ok := ifx.Left.(*StringLiteral)
	if !ok {
		t.Fatalf("expected left string got %s\n", reflect.TypeOf(ifx.Left))
	}

	if left.Value != "string" {
		t.Fatalf("expected left to be 'string' got %s\n", left.Value)
	}

	right, ok := ifx.Right.(*StringLiteral)

	if !ok {
		t.Fatalf("expected right to be str lit got %s\n", reflect.TypeOf(ifx.Right))
	}

	if right.Value != " works" {
		t.Fatalf("expected right ' works' but got %s\n", right.Value)
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.NewLexer(input)

	p := New(l)

	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d\n", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatal("expected expression statement")
	}

	ident, ok := stmt.Expression.(*Identifier)
	if !ok {
		t.Fatal("expected identifie literal")
	}

	if ident.TokenLiteral() != "foobar" {
		t.Fatalf("identifier token expected %s got %s", "foobar", ident.TokenLiteral())
	}
}

func TestIntegerExpression(t *testing.T) {
	input := "5;"

	l := lexer.NewLexer(input)

	p := New(l)

	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d\n", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatal("expected expression statement")
	}

	lit, ok := stmt.Expression.(*IntegerLiteral)
	if !ok {
		t.Fatal("expected literal")
	}

	if lit.TokenLiteral() != "5" {
		t.Fatalf("identifier token expected %s got %s\n", "foobar", lit.TokenLiteral())
	}

	if lit.Value != 5 {
		t.Fatalf("expected literal %d got %d\n", 5, lit.Value)
	}
}

func TestParsePrefixOperator(t *testing.T) {
	input := `
!5;
--111;
`

	tests := []struct {
		Operator string
		input    string
		value    int64
	}{
		{"!", "!5", 5},
		{"--", "--111", 111},
	}

	l := lexer.NewLexer(input)
	p := New(l)
	program := p.ParseProgram()

	if len(program.Statements) != 2 {
		t.Errorf("expected 1 statement got %d\n", len(program.Statements))
	}

	for i, s := range program.Statements {
		stmt, ok := s.(*ExpressionStatement)

		if !ok {
			t.Errorf("received invalid expression statements for %d got %s\n", i+1, reflect.TypeOf(stmt).String())
		}

		x, ok := stmt.Expression.(*PrefixExpression)
		if !ok {
			t.Errorf("expected a prefix expression")
		}

		if x.Operator != tests[i].Operator {
			t.Errorf("expected operator %s but got %s\n", x.Operator, tests[i].Operator)
		}

		right, ok := x.Right.(*IntegerLiteral)
		if !ok {
			t.Error("right not an int literal")
		}

		if right.Value != tests[i].value {
			t.Errorf("expected value %d but got %d\n", x.Right, tests[i].value)
		}
	}

}

func TestParseInfixExpression(t *testing.T) {
	tests := []struct {
		input    string
		left     int64
		operator string
		right    int64
	}{
		{`5 + 5`, 5, "+", 5},
		{"2 - 1;", 2, "-", 1},
		{"3 * 3;", 3, "*", 3},
		{"1 > 2;", 1, ">", 2},
	}
	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := New(l)

		program := p.ParseProgram()

		if len(program.Statements) != 1 {
			t.Fatalf("expected 1 statement got %d\n", len(program.Statements))
		}
		s := program.Statements[0]

		stmt, ok := s.(*ExpressionStatement)

		if !ok {
			t.Fatalf("expected expression statement got %s\n", reflect.TypeOf(s).String())
		}

		exp, ok := stmt.Expression.(*InfixExpression)
		if !ok {
			t.Fatalf("expected infix expression got %s\n", reflect.TypeOf(stmt).String())
		}

		left, ok := exp.Left.(*IntegerLiteral)
		if !ok {
			t.Fatalf("expected integer literal got %s\n", reflect.TypeOf(exp.Left).String())
		}

		if left.Value != tt.left {
			t.Fatalf("expected left to be %d got %d\n", tt.left, left.Value)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("expected operator %s got %s\n", tt.operator, exp.Operator)
		}

		right, ok := exp.Right.(*IntegerLiteral)
		if !ok {
			t.Fatalf("expected integer literal got %s\n", reflect.TypeOf(exp.Right).String())
		}

		if right.Value != tt.right {
			t.Fatalf("expected value %d got %d\n", tt.right, right.Value)
		}
	}
}

func TestParseBoolean(t *testing.T) {
	tests := []struct {
		input string
		value bool
	}{
		{`let a = false`, false},
		{`let b = true`, true},
	}

	for i, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := New(l)
		program := p.ParseProgram()

		if len(program.Statements) != 1 {
			t.Fatalf("expected 1 statement got %d\n", len(program.Statements))
		}

		ls, ok := program.Statements[0].(*LetStatement)
		if !ok {
			t.Fatalf("expected let statement got %s\n", reflect.TypeOf(program.Statements[0]))
		}

		bl, ok := ls.Value.(*BoolLiteral)
		if !ok {
			t.Fatalf("expected bool literal got %s\n", reflect.TypeOf(ls.Value))
		}

		if bl.Value != tt.value {
			t.Fatalf("expected %t got %t\n", tt.value, bl.Value)
		}

		fmt.Printf("got %t value for statement %d\n", bl.Value, i)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
    {`a * [1,2,3][b * c] * d`, `((a * ([1, 2, 3][(b * c)])) * d)`},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := New(l)
		program := p.ParseProgram()
		for _, stmt := range program.Statements {
			str := stmt.String()
			if tt.expected != str {
				t.Errorf("expected '%s' got '%s'\n", tt.expected, str)
			}
		}
	}

}

func TestIfExpression(t *testing.T) {
	input := `
  if (x < y) { 
    return x 
  } else {
    return y;
  }
  `

	l := lexer.NewLexer(input)
	p := New(l)
	program := p.ParseProgram()

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement got %d\n", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ExpressionStatement)
	println(program.String())
	if !ok {
		t.Fatalf("expected expression statement got %s\n", reflect.TypeOf(stmt))
	}

	exp, ok := stmt.Expression.(*IfExpression)
	if !ok {
		t.Fatalf("expected if expression got %s\n", reflect.TypeOf(stmt.Expression))
	}

	if exp.Condition.String() != "(x < y)" {
		t.Fatalf("got unexpected condition '%s'\n", exp.Condition.String())
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Fatalf("expected 1 consequence got %d\n", len(exp.Consequence.Statements))
	}

	rs, ok := exp.Consequence.Statements[0].(*ReturnStatement)

	if !ok {
		t.Fatalf("expected return statement got %s\n", reflect.TypeOf(exp.Consequence.Statements[0]))
	}

	if rs.String() != `(return x)` {
		t.Fatalf("unexpected return statement got %s\n", rs.String())
	}

}

func TestParseFnLiteral(t *testing.T) {
	input := "let a = fn(x,y) { return x + y; }"

	l := lexer.NewLexer(input)
	p := New(l)
	program := p.ParseProgram()
	println(program.String())
	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement got %d\n", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*LetStatement)

	if !ok {
		t.Fatalf("expected let statement got %s\n", reflect.TypeOf(program.Statements[0]))
	}

	fn, ok := stmt.Value.(*FnLiteral)
	if !ok {
		t.Fatalf("expected fn literal got %s\n", reflect.TypeOf(stmt.Value))
	}

	if len(fn.Params) != 2 {
		t.Fatalf("expected %d params got %d\n", 2, len(fn.Params))
	}
	println(fn.String())
}

func TestParseCallExpression(t *testing.T) {
	input := "add(3,5*5+2,9)"

	l := lexer.NewLexer(input)
	p := New(l)
	program := p.ParseProgram()

	if len(program.Statements) != 1 {
		t.Errorf("expected 1 statement got %d\n", len(program.Statements))
	}

	// statement should be  exp statement
	stmt, ok := program.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("exprected call exp got %s\n", reflect.TypeOf(program.Statements[0]))
	}

	// statement exp sould be call exp
	ce, ok := stmt.Expression.(*CallExpression)
	if !ok {
		t.Errorf("expected call exp got %s\n", reflect.TypeOf(stmt))
	}
	// call fn should be add

	fn, ok := ce.Function.(*Identifier)
	if !ok {
		t.Errorf("expected fn to be identifier got %s\n", reflect.TypeOf(ce.Function))
	}

	if fn.Value != "add" {
		t.Errorf("expected fn value to be add got %s\n", fn.Value)
	}

	if len(ce.Arguments) != 3 {
		t.Errorf("expected 3 args got %d\n", len(ce.Arguments))
	}

	println(ce.String()) // call
}

func TestParsesArray(t *testing.T) {
  input := `[1, "two", len("5")]` 
  
  l := lexer.NewLexer(input)
  p := New(l)
  program := p.ParseProgram()

  stmt, ok := program.Statements[0].(*ExpressionStatement)
  if !ok {
    t.Fatalf("expected exp statement got %s\n", reflect.TypeOf(program.Statements[0]).String())
  }

  arr, ok := stmt.Expression.(*ArrayLiteral)
  if !ok {
    t.Errorf("expected array literal got %s\n", reflect.TypeOf(stmt.Expression).String())
  }

  zero, ok := arr.Elements[0].(*IntegerLiteral)
  if !ok {
    t.Errorf("expected array literal got %s\n", reflect.TypeOf(stmt.Expression).String())
  }

  if zero.Value != 1 {
    t.Errorf("input[0] expected one got %d\n", zero.Value)
  }

  one, ok := arr.Elements[1].(*StringLiteral)
  if !ok {
    t.Errorf("input[1] expected string got %s\n", reflect.TypeOf(arr.Elements[1]).String())
  }

  if one.Value != "two" {
    t.Errorf("got wrong str value")
  }

  two, ok := arr.Elements[2].(*CallExpression)
  if !ok {
    t.Errorf("input[2] expected call got %s\n", reflect.TypeOf(arr.Elements[2]).String())
  }


  arg, ok := two.Arguments[0].(*StringLiteral)
  if !ok {
    t.Errorf("expected string call arg got %s\n", reflect.TypeOf(two.Arguments[0]).String())
  }

  if arg.Value != "5" {
    t.Errorf("expected arg value '5' got %s\n", arg.Value)
  }
}

func TestParseIndexExpression(t *testing.T) {
  input := "arr[1]"
  l := lexer.NewLexer(input)
  p := New(l)
  pr := p.ParseProgram()

  stmt, ok := pr.Statements[0].(*ExpressionStatement)
  
  ie, ok := stmt.Expression.(*IndexExpression)
  if !ok {
    t.Errorf("expected index expression got %s\n", reflect.TypeOf(stmt))
  }

  ident, ok := ie.Left.(*Identifier)
  if !ok {
    t.Errorf("didn't get identifier got %s\n", reflect.TypeOf(ie.Left))
  }

  if ident.Value != "arr" {
    t.Errorf("expected name 'arr' got '%s'\n", ident.Value)
  }

  i, ok := ie.Index.(*IntegerLiteral)
  if !ok {
    t.Errorf("expected int index got %s\n", reflect.TypeOf(ie.Index))
  }

  if i.Value != 1 {
    t.Errorf("expected index value of 1 got %d\n", i.Value)
  }
}

func checkParserErrors(t *testing.T, p *Parser) {
	for _, err := range p.errors {
		t.Errorf(err)
	}
}
