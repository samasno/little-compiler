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
		{tokens.PLUS, tokens.PLUS},
		{tokens.LPAREN, tokens.LPAREN},
		{tokens.RPAREN, tokens.RPAREN},
		{tokens.LBRACE, tokens.LBRACE},
		{tokens.RBRACE, tokens.RBRACE},
		{tokens.COLON, tokens.COLON},
		{tokens.SEMICOLON, tokens.SEMICOLON},
		{tokens.COMMA, tokens.COMMA},
		{tokens.IDENTIFIER, "test"},
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

func TestSequencesOfTokens(t *testing.T) {
	test := `
let a = 100
let b=3
let c = fn(x,y){
	return x * y;
}
"string works";
"one";
`
	ex := []tokens.Token{
		{tokens.LET, tokens.LET},
		{tokens.IDENTIFIER, `a`},
		{tokens.ASSIGN, tokens.ASSIGN},
		{tokens.INTEGER, "100"},
		{tokens.LET, tokens.LET},
		{tokens.IDENTIFIER, `b`},
		{tokens.ASSIGN, tokens.ASSIGN},
		{tokens.INTEGER, "3"},
		{tokens.LET, tokens.LET},
		{tokens.IDENTIFIER, `c`},
		{tokens.ASSIGN, tokens.ASSIGN},
		{tokens.FUNCTION, tokens.FUNCTION},
		{tokens.LPAREN, tokens.LPAREN},
		{tokens.IDENTIFIER, `x`},
		{tokens.COMMA, tokens.COMMA},
		{tokens.IDENTIFIER, `y`},
		{tokens.RPAREN, tokens.RPAREN},
		{tokens.LBRACE, tokens.LBRACE},
		{tokens.RETURN, tokens.RETURN},
		{tokens.IDENTIFIER, `x`},
		{tokens.MULTIPLY, tokens.MULTIPLY},
		{tokens.IDENTIFIER, `y`},
		{tokens.SEMICOLON, tokens.SEMICOLON},
		{tokens.RBRACE, tokens.RBRACE},
		{tokens.STRING, "string works"},
		{tokens.SEMICOLON, tokens.SEMICOLON},
		{tokens.STRING, "one"},
		{tokens.SEMICOLON, tokens.SEMICOLON},
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

func TestTokenize(t *testing.T) {
	test := `
let a = 100`
	l := NewLexer(test)

	ts := l.Tokenize()
	for _, t := range ts {
		fmt.Printf("\n%v", t)
	}
}
