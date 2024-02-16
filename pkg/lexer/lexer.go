package lexer

import (
	"strconv"

	"github.com/samasno/little-compiler/pkg/tokens"
)

type Lexer struct {
	source       string
	sourceLength int
	position     int
	readPosition int
}

func NewLexer(source string) *Lexer {
	return &Lexer{
		position:     0,
		readPosition: 0,
		sourceLength: len(source),
		source:       source,
	}
}

func (l *Lexer) Tokenize() []tokens.Token {
	ts := []tokens.Token{}

	for {
		tk := l.NextToken()
		ts = append(ts, tk)
		if tk.Type == tokens.EOF {
			break
		}
	}

	return ts
}

func (l *Lexer) NextToken() tokens.Token {
	// skip whitespace function here
	l.eatWhiteSpace()
	l.position = l.readPosition
	var t tokens.Token

	if l.position >= l.sourceLength {
		t.Type, t.Literal = tokens.EOF, tokens.EOF
		return t
	}

	next := string(l.source[l.position])

	switch {
	case isDelimiter(next),
		isOperator(next):
		l.readPosition++
		return tokens.Token{next, next}
	default:
		return l.readSequence()
	}
}

func (l *Lexer) eatWhiteSpace() {
	for {
		var next string
		if l.readPosition < l.sourceLength {
			next = string(l.source[l.readPosition])
		}

		if isWhiteSpace(next) && l.readPosition < l.sourceLength {
			l.readPosition++
		} else {
			break
		}
	}
}

func (l *Lexer) readSequence() tokens.Token {

loop:
	for {
		var next string
		l.readPosition++

		if l.readPosition < l.sourceLength {
			next = string(l.source[l.readPosition])
		}

		switch {
		case isDelimiter(next),
			isWhiteSpace(next),
			isOperator(next),
			l.readPosition == l.sourceLength:
			break loop
		}
	}

	v := l.source[l.position:l.readPosition]
	switch {
	case isInteger(v):
		return tokens.Token{tokens.INTEGER, v}
	case isKeyword(v):
		return tokens.Token{v, v}
	default:
		return tokens.Token{tokens.IDENTIFIER, v}
	}
}

func isDelimiter(c string) bool {
	switch c {
	case tokens.COMMA,
		tokens.SEMICOLON,
		tokens.COLON,
		tokens.LBRACE,
		tokens.RBRACE,
		tokens.LPAREN,
		tokens.RPAREN,
		tokens.LBRACKET,
		tokens.RBRACKET,
		tokens.EOF:
		return true
	default:
		return false

	}
}

func isWhiteSpace(s string) bool {
	switch s {
	case tokens.SPACE,
		tokens.TAB,
		tokens.NEWLINE:
		return true
	default:
		return false
	}
}

func isInteger(s string) bool {
	_, err := strconv.Atoi(s)
	if err != nil {
		return false
	}
	return true
}

func isOperator(s string) bool {
	switch s {
	case tokens.ASSIGN,
		tokens.PLUS,
		tokens.MINUS,
		tokens.MULTIPLY,
		tokens.DIVIDE:
		return true
	default:
		return false
	}
}

func isKeyword(s string) bool {
	switch s {
	case tokens.LET,
		tokens.INT,
		tokens.FUNCTION,
		tokens.RETURN:
		return true
	default:
		return false
	}

}

func isNewline(s string) bool {
	return s == `\n`
}
