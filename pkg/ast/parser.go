package ast

import (
	"log"
	"strconv"

	"github.com/samasno/little-compiler/pkg/lexer"
	"github.com/samasno/little-compiler/pkg/tokens"
)

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:               l,
		prefixParseFn:   make(map[string]prefixParseFn),
		infixParseFn:    make(map[string]infixParseFn),
		infixPrecedence: make(map[string]int),
	}

	p.registerPrefix(p.parseGroupedExpression, tokens.LPAREN)

	p.registerPrefix(p.parseIdentifier,
		tokens.IDENTIFIER,
	)

	p.registerPrefix(p.parseInteger,
		tokens.INTEGER,
	)

	p.registerPrefix(p.parseBoolean,
		tokens.TRUE,
		tokens.FALSE,
	)

	p.registerPrefix(p.parseIfExpression,
		tokens.IF,
	)

  p.registerPrefix(p.parseFnLiteral,
    tokens.FUNCTION,
  )

	p.registerInfix(p.parseInfixExpression,
		tokens.ASSIGN,
		tokens.PLUS,
		tokens.MINUS,
		tokens.GT,
		tokens.LT,
		tokens.MULTIPLY,
		tokens.DIVIDE,
		tokens.EQUALTO,
		tokens.GTE,
		tokens.LTE,
	)

	p.registerPrefix(p.parsePrefixExpression,
		tokens.NOT,
		tokens.MINUS,
		tokens.INC,
		tokens.DEC,
	)

	p.registerInfixPrecedence(SUM,
		tokens.PLUS,
		tokens.MINUS,
	)

	p.registerInfixPrecedence(PRODUCT,
		tokens.MULTIPLY,
		tokens.DIVIDE,
	)

	p.registerInfixPrecedence(EQUALS,
		tokens.ASSIGN,
		tokens.EQUALTO,
		tokens.LTE,
		tokens.GTE,
	)

	p.registerInfixPrecedence(LESSGREATER,
		tokens.LTE,
		tokens.GTE,
		tokens.LT,
		tokens.GT,
	)

	p.registerInfixPrecedence(CALL,
		tokens.FUNCTION,
	)

	p.nextToken()

	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *Program {
	pr := &Program{}

	for p.currentToken.Type != tokens.EOF {
		stmt := p.parseStatement()

		if stmt != nil {
			pr.Statements = append(pr.Statements, stmt)
		}

		p.nextToken()
	}

	return pr
}

func (p *Parser) parseStatement() Statement {
	var stmt Statement

	switch p.currentToken.Type {

	case tokens.EOF, tokens.SEMICOLON:
		break

	case tokens.LBRACE:
		stmt = p.parseBlockStatement()

	case tokens.RETURN:
		stmt = p.parseReturnStatement()

	case tokens.LET:
		stmt = p.parseLet()

	default:
		stmt = p.parseExpressionStatement()
	}

	return stmt
}

func (p *Parser) parseLet() Statement {
	stmt := &LetStatement{Token: p.currentToken}

	p.expectPeek(tokens.IDENTIFIER)

	stmt.Name = &Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	p.expectPeek(tokens.ASSIGN)

	p.nextToken()

	switch {
	case p.currentIs(tokens.INTEGER):
		stmt.Value = p.parseInteger()
	case p.currentIs(tokens.TRUE), p.currentIs(tokens.FALSE):
		stmt.Value = p.parseBoolean()
	case p.currentIs(tokens.FUNCTION):
		stmt.Value = p.parseFnLiteral()
	default:
		log.Fatalf("unexpected token %s\n", p.currentToken.Type)
	}

	return stmt
}

func (p *Parser) parseReturnStatement() Statement {

	r := &ReturnStatement{
		Token: p.currentToken,
	}

	p.nextToken()

	r.Value = p.parseExpression(LOWEST)

	return r
}

func (p *Parser) parseExpressionStatement() Statement {
	stmt := &ExpressionStatement{Token: p.currentToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekIs(tokens.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) Expression {
	prefix, ok := p.prefixParseFn[p.currentToken.Type]

	if !ok {
		log.Fatalf("no parsing function found for %s\n", p.currentToken.Type)
	}

	lexp := prefix()

	for !p.peekIs(tokens.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFn[p.peekToken.Type]

		if infix == nil {
			return lexp
		}

		p.nextToken()

		lexp = infix(lexp)

	}

	return lexp
}

func (p *Parser) expectPeek(tokenType string) bool {
	if !p.peekIs(tokenType) {
		log.Fatalf("unexpected token. expected %s got %s\n", tokenType, p.peekToken.Type)
	}

	p.nextToken()
	return true
}

func (p *Parser) peekPrecedence() int {
	precedence, ok := p.infixPrecedence[p.peekToken.Type]

	if ok {
		return precedence
	}

	return LOWEST
}

func (p *Parser) currentPrecedence() int {
	precedence, ok := p.infixPrecedence[p.currentToken.Type]

	if ok {
		return precedence
	}

	return LOWEST
}

func (p *Parser) peekIs(tokenType string) bool {
	if p.peekToken.Type != tokenType {
		return false
	}

	return true
}

func (p *Parser) currentIs(tokenType string) bool {
	if p.currentToken.Type != tokenType {
		return false
	}

	return true
}

func (p *Parser) registerPrefix(fn prefixParseFn, tokenTypes ...string) {
	for _, t := range tokenTypes {
		p.prefixParseFn[t] = fn
	}
}

func (p *Parser) registerInfix(fn infixParseFn, tokenTypes ...string) {
	for _, t := range tokenTypes {
		p.infixParseFn[t] = fn
	}
}

func (p *Parser) registerInfixPrecedence(precedence int, tokenTypes ...string) {
	for _, tt := range tokenTypes {
		p.infixPrecedence[tt] = precedence
	}
}

func (p *Parser) registerPrefixPrecedence(precedence int, tokenTypes ...string) {
	for _, tt := range tokenTypes {
		p.prefixPrecedence[tt] = precedence
	}
}

func (p *Parser) parseIdentifier() Expression {
	return &Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parsePrefixExpression() Expression {
	px := &PrefixExpression{Token: p.currentToken, Operator: p.currentToken.Literal}

	p.nextToken()

	px.Right = p.parseExpression(PREFIX)

	return px
}

func (p *Parser) parseInfixExpression(left Expression) Expression {
	ix := &InfixExpression{Token: p.currentToken, Operator: p.currentToken.Literal, Left: left}
	prec, ok := p.infixPrecedence[p.currentToken.Literal]

	if !ok {
		log.Fatalf("no precedence found for %s \n", p.currentToken.Literal)
	}

	p.nextToken()

	ix.Right = p.parseExpression(prec)

	return ix
}

func (p *Parser) parseInteger() Expression {
	v, err := strconv.Atoi(p.currentToken.Literal)

	if err != nil {
		return nil
	}

	return &IntegerLiteral{Token: p.currentToken, Value: int64(v)}
}

func (p *Parser) parseBoolean() Expression {
	v, err := strconv.ParseBool(p.currentToken.Literal)
	if err != nil {
		log.Fatal("invalid bool literal")
	}
	return &BoolLiteral{Token: p.currentToken, Value: v}
}

func (p *Parser) parseGroupedExpression() Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(tokens.RPAREN) {
		log.Fatal("unexpected token, expected RPAREN")
	}

	return exp
}

func (p *Parser) parseIfExpression() Expression {
	exp := &IfExpression{Token: p.currentToken}

	p.expectPeek(tokens.LPAREN)

	p.nextToken()

	exp.Condition = p.parseExpression(LOWEST)

	p.expectPeek(tokens.RPAREN)

	p.expectPeek(tokens.LBRACE)

	exp.Consequence = p.parseBlockStatement()

	if p.peekIs(tokens.ELSE) {

		p.nextToken()

		exp.Alternative = p.parseBlockStatement()
		println(exp.Alternative.String())
	}

	return exp
}

func (p *Parser) parseFnLiteral() Expression {
	fn := &FnLiteral{Token: p.currentToken}
  
  p.expectPeek(tokens.LPAREN)

  fn.Params = p.parseFnParams()
 
  if p.peekIs(tokens.LBRACE) {
    p.expectPeek(tokens.LBRACE)
  
    fn.Body = p.parseBlockStatement()
  }
  
	return fn
}

func (p *Parser) parseFnParams() []*Identifier {
  params := []*Identifier{}
  
  for !p.peekIs(tokens.RPAREN) && !p.peekIs(tokens.EOF) {
    p.nextToken()

    if !p.currentIs(tokens.COMMA) && !p.currentIs(tokens.IDENTIFIER) {
      log.Fatalf("expected identifier or comma got %s\n", p.currentToken.Type)
    }
    
    if p.currentIs(tokens.IDENTIFIER) {
      params = append(params, &Identifier{Token: p.currentToken, Value: p.currentToken.Literal})
    }

  }
   
  p.nextToken()

  return params
}

func (p *Parser) parseBlockStatement() *BlockStatement {
	bs := &BlockStatement{Token: p.currentToken}

	bs.Statements = []Statement{}

	p.nextToken()

	for !p.currentIs(tokens.RBRACE) && !p.currentIs(tokens.EOF) {
		stmt := p.parseStatement()

		if stmt != nil {
			bs.Statements = append(bs.Statements, stmt)
		}

		p.nextToken()
  }

	return bs
}

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)
