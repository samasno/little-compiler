package lexer

import (
	"fmt"
	"testing"

	"github.com/samasno/little-compiler/pkg/tokens"
)

func TestNextTokenDelimiters(t *testing.T) {
	input := `= +(){}:;,
test
`
	expected := []tokens.Token{
		{tokens.ASSIGN, tokens.ASSIGN},
		{tokens.SPACE, tokens.SPACE},
		{tokens.PLUS, tokens.PLUS},
		{tokens.LPAREN, tokens.LPAREN},
		{tokens.RPAREN, tokens.RPAREN},
		{tokens.LBRACE, tokens.LBRACE},
		{tokens.RBRACE, tokens.RBRACE},
		{tokens.COLON, tokens.COLON},
		{tokens.SEMICOLON, tokens.SEMICOLON},
		{tokens.COMMA, tokens.COMMA},
		{tokens.NEWLINE, tokens.NEWLINE},
		{tokens.IDENTIFIER, "test"},
		{tokens.NEWLINE, tokens.NEWLINE},
		{tokens.EOF, tokens.EOF},
	}

	l := NewLexer(input)

	for _, in := range expected {
		out := l.NextToken()

		if out.Literal != in.Literal {
			fmt.Printf("\nexpected %s got %s", in.Literal, out.Literal)
			t.Fail()
		}
	}
}

func TestLetDeclaration(t *testing.T) {
	test := `
let a = 100
let b=3
let c = fn(x,y){
	return x * y;
}
`
	ex := []tokens.Token{
		{tokens.NEWLINE, tokens.NEWLINE},
		{tokens.LET, tokens.LET},
		{tokens.SPACE, tokens.SPACE},
		{tokens.IDENTIFIER, `a`},
		{tokens.SPACE, tokens.SPACE},
		{tokens.ASSIGN, tokens.ASSIGN},
		{tokens.SPACE, tokens.SPACE},
		{tokens.INTEGER, "100"},
		{tokens.NEWLINE, tokens.NEWLINE},
		{tokens.LET, tokens.LET},
		{tokens.SPACE, tokens.SPACE},
		{tokens.IDENTIFIER, `b`},
		{tokens.ASSIGN, tokens.ASSIGN},
		{tokens.INTEGER, "3"},
		{tokens.NEWLINE, tokens.NEWLINE},
		{tokens.LET, tokens.LET},
		{tokens.SPACE, tokens.SPACE},
		{tokens.IDENTIFIER, `c`},
		{tokens.SPACE, tokens.SPACE},
		{tokens.ASSIGN, tokens.ASSIGN},
		{tokens.SPACE, tokens.SPACE},
		{tokens.FUNCTION, tokens.FUNCTION},
		{tokens.LPAREN, tokens.LPAREN},
		{tokens.IDENTIFIER, `x`},
		{tokens.COMMA, tokens.COMMA},
		{tokens.IDENTIFIER, `y`},
		{tokens.RPAREN, tokens.RPAREN},
		{tokens.LBRACE, tokens.LBRACE},
		{tokens.NEWLINE, tokens.NEWLINE},
		{tokens.TAB, tokens.TAB},
		{tokens.RETURN, tokens.RETURN},
		{tokens.SPACE, tokens.SPACE},
		{tokens.IDENTIFIER, `x`},
		{tokens.SPACE, tokens.SPACE},
		{tokens.MULTIPLY, tokens.MULTIPLY},
		{tokens.SPACE, tokens.SPACE},
		{tokens.IDENTIFIER, `y`},
		{tokens.SEMICOLON, tokens.SEMICOLON},
		{tokens.NEWLINE, tokens.NEWLINE},
		{tokens.RBRACE, tokens.RBRACE},
		{tokens.NEWLINE, tokens.NEWLINE},
		{tokens.EOF, tokens.EOF},
	}

	l := NewLexer(test)

	for i, x := range ex {
		tk := l.NextToken()

		if tk.Literal != x.Literal || tk.Type != x.Type {
			t.Fail()
			fmt.Printf("\nchar %d expected %s:%s got %s:%s", i, x.Type, x.Literal, tk.Type, tk.Literal)
		}
	}
}
