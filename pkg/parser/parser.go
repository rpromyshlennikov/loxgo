package parser

import (
	"fmt"

	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/errors"
	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/scanner"
)

type Parser struct {
	tokens  []Token
	current int

	errReporter errors.Reporter
}

func NewParser(tokens []Token, errReporter errors.Reporter) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,

		errReporter: errReporter,
	}
}

func (p *Parser) Parse() (astTree Expr) {
	defer func() {
		if recovered := recover(); recovered != nil {
			astTree = nil
		}
	}()
	return p.expression()
}

func (p *Parser) expression() Expr {
	return p.equality()
}

// Binary expressions.

func (p *Parser) equality() Expr {
	expr := p.comparison()
	for p.match(scanner.BANGEQUAL, scanner.EQUALEQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()
	for p.match(scanner.GREATER, scanner.GREATEREQUAL, scanner.LESS, scanner.LESSEQUAL) {
		operator := p.previous()
		right := p.term()
		expr = NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()
	for p.match(scanner.MINUS, scanner.PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()
	for p.match(scanner.SLASH, scanner.STAR) {
		operator := p.previous()
		right := p.unary()
		expr = NewBinary(expr, operator, right)
	}
	return expr
}

// Unary expression.

func (p *Parser) unary() Expr {
	if p.match(scanner.BANG, scanner.MINUS) {
		operator := p.previous()
		right := p.unary()
		return NewUnary(operator, right)
	}
	return p.primary()
}

// Primary expression.

func (p *Parser) primary() Expr {
	if p.match(scanner.FALSE) {
		return NewLiteral(false)
	}
	if p.match(scanner.TRUE) {
		return NewLiteral(true)
	}
	if p.match(scanner.NIL) {
		return NewLiteral(nil)
	}
	if p.match(scanner.NUMBER, scanner.STRING) {
		return NewLiteral(p.previous().Literal())
	}
	if p.match(scanner.LEFTPAREN) {
		expr := p.expression()
		p.consume(scanner.RIGHTPAREN, "Expect ')' after expression.")
		return NewGrouping(expr)
	}
	panic(p.erro(p.peek(), "Expect expression."))
}

// helpers.

func (p *Parser) match(types ...TokenType) bool {
	for _, kind := range types {
		if p.check(kind) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(kind TokenType, errMsg string) Token {
	if p.check(kind) {
		return p.advance()
	}
	panic(p.erro(p.peek(), errMsg))
}

func (p *Parser) check(kind TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Kind() == kind
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Kind() == scanner.EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) erro(token Token, message string) error {
	if token.Kind() == scanner.EOF {
		p.errReporter(token.Line(), " at end. Message: "+message)
	} else {
		p.errReporter(token.Line(), " at '"+token.Lexeme()+"'. Message: "+message)
	}

	return fmt.Errorf("parse error")
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Kind() == scanner.SEMICOLON {
			return
		}

		switch p.peek().Kind() {
		case scanner.CLASS, scanner.FUN, scanner.VAR, scanner.FOR, scanner.IF:
			return
		}
		p.advance()
	}
}
