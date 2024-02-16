package ast

import (
	"github.com/samasno/little-compiler/pkg/lexer"
	"github.com/samasno/little-compiler/pkg/tokens"
)

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

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
	var statement Statement
	for p.currentToken.Type != tokens.EOF { // bases follow on actions based on current token
		switch p.currentToken.Type {
		case tokens.FUNCTION:
			println("parse function tokens")
			statement = p.parseFunction()
		case tokens.LET:
			println("parse let statements")
			statement = p.parseLet()
		case tokens.IDENTIFIER:
			println("check identifier type")
			statement = p.parseIdentifier()
		case "if":
			println("parsing if statement")
			statement = p.parseIf()
		}

		if statement != nil {
			pr.Statements = append(pr.Statements, statement)
		}
	}

	return pr
}

func (p *Parser) parseIf() Statement {
	return nil
}

func (p *Parser) parseFunction() Statement {
	return nil
}

func (p *Parser) parseIdentifier() Statement {
	return nil
}

func (p *Parser) parseExpression() Statement {
	return nil
}

func (p *Parser) parseGroupedExpression() Statement {
	return nil
}

func (p *Parser) parseLet() Statement {
	// shift forward to next token that is not space (need to remove spaces tokens)
	// parse identifier name, must be identifier token following let
	// next token should be assign, if not it is invalid
	// shift forward
	// next may be integer, function, or expression
	// if integer, lookahead to see if it an expression, use recursive parse

	return nil
}
