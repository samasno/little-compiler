package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/samasno/little-compiler/pkg/lexer"
	"github.com/samasno/little-compiler/pkg/tokens"
)

type Parser struct {
	l *lexer.Lexer

	currentToken tokens.Token
	peekToken    tokens.Token

	errors []string

	prefixParseFn    map[string]prefixParseFn
	infixParseFn     map[string]infixParseFn
	prefixPrecedence map[string]int
	infixPrecedence  map[string]int
}

type (
	prefixParseFn func() Expression
	infixParseFn  func(Expression) Expression
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement `json:"statements"`
}

type LetStatement struct {
	Token tokens.Token `json:"token"`
	Name  *Identifier  `json:"name"`
	Value Expression   `json:"value"`
}

type IfExpression struct {
	Token       tokens.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

type BlockStatement struct {
	Token      tokens.Token
	Statements []Statement
}

type FnLiteral struct {
	Token  tokens.Token
	Params []*Identifier
	Body   *BlockStatement
}

type Identifier struct {
	Token tokens.Token
	Value string
}

type StaticValue struct {
	Token tokens.Token
	Value string
}

type ExpressionValue struct {
	Token tokens.Token
	Value string
}

type ReturnStatement struct {
	Token tokens.Token
	Value Expression
}

type ExpressionStatement struct {
	Token      tokens.Token
	Expression Expression
}

type IntegerLiteral struct {
	Token tokens.Token
	Value int64
}

type BoolLiteral struct {
	Token tokens.Token
	Value bool
}

type PrefixExpression struct {
	Token    tokens.Token
	Operator string
	Right    Expression
}

type InfixExpression struct {
	Token    tokens.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (fl *FnLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Params {
		params = append(params, p.String())
	}
	out.WriteString(fl.Token.Literal)
	out.WriteString(fmt.Sprintf("(%s) ", strings.Join(params, ",")))
	out.WriteString(fl.Body.String())
	return out.String()
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.Value != nil {
		out.WriteString(rs.Value.String())
	}
	out.WriteString(")")
	return out.String()
}

func (e *ExpressionStatement) String() string {
	var out bytes.Buffer
	if e.Expression != nil {
		out.WriteString(e.Expression.String())
	}
	return out.String()
}

func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String() + " = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	return out.String()
}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

func (is *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString("if")
	out.WriteString(is.Condition.String())
	out.WriteString(" ")
	out.WriteString(is.Consequence.String())
	out.WriteString(")")
	return out.String()
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (px *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString(fmt.Sprintf("(%s%s)", px.Operator, px.Right.String()))
	return out.String()
}

func (is *IfExpression) TokenLiteral() string { return is.Token.Literal }

func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }

func (ix *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString(fmt.Sprintf("(%s %s %s)", ix.Left.String(), ix.Operator, ix.Right.String()))
	return out.String()
}

func (fl *FnLiteral) expressionNode()               {}
func (fl *FnLiteral) TokenLiteral() string          { return fl.Token.Literal }
func (bs *BlockStatement) statementNode()           {}
func (is *IfExpression) expressionNode()            {}
func (bl *BoolLiteral) String() string              { return bl.Token.Literal }
func (bl *BoolLiteral) TokenLiteral() string        { return bl.Token.Literal }
func (bl *BoolLiteral) expressionNode()             {}
func (ix *InfixExpression) TokenLiteral() string    { return ix.Token.Literal }
func (ix *InfixExpression) expressionNode()         {}
func (px *PrefixExpression) TokenLiteral() string   { return px.Token.Literal }
func (px *PrefixExpression) expressionNode()        {}
func (ls *LetStatement) TokenLiteral() string       { return ls.Token.Literal }
func (ls *LetStatement) statementNode()             {}
func (i *Identifier) TokenLiteral() string          { return i.Token.Literal }
func (i *Identifier) String() string                { return i.Value }
func (i *Identifier) expressionNode()               {}
func (s *StaticValue) expressionNode()              {}
func (s *StaticValue) TokenLiteral() string         { return s.Token.Literal }
func (e *ExpressionValue) expressionNode()          {}
func (e *ExpressionValue) TokenLiteral() string     { return e.Token.Literal }
func (r *ReturnStatement) TokenLiteral() string     { return r.Token.Literal }
func (r *ReturnStatement) statementNode()           {}
func (e *ExpressionStatement) TokenLiteral() string { return e.Token.Literal }
func (e *ExpressionStatement) statementNode()       {}
func (i *IntegerLiteral) TokenLiteral() string      { return i.Token.Literal }
func (i *IntegerLiteral) String() string            { return i.Token.Literal }
func (i *IntegerLiteral) expressionNode()           {}
