package parser

import (
	"fmt"

	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/errors"
	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/parser/ast"
	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/scanner"
)

type Parser struct {
	tokens  []scanner.Token
	current int

	errReporter errors.Reporter
}

func NewParser(tokens []scanner.Token, errReporter errors.Reporter) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,

		errReporter: errReporter,
	}
}

func (p *Parser) Parse() (astTree []ast.Stmt) {
	defer func() {
		if recovered := recover(); recovered != nil {
			astTree = nil
		}
	}()
	var statements []ast.Stmt
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	return statements
}

func (p *Parser) expression() ast.Expr {
	return p.assignment()
}

func (p *Parser) declaration() (astTree ast.Stmt) {
	defer func() {
		if recovered := recover(); recovered != nil {
			astTree = nil
			p.synchronize()
		}
	}()
	if p.match(scanner.VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) statement() ast.Stmt {
	if p.match(scanner.PRINT) {
		return p.printStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) printStatement() ast.Stmt {
	value := p.expression()
	p.consume(scanner.SEMICOLON, "Expect ';' after value.")
	return ast.NewPrint(value)
}

func (p *Parser) varDeclaration() ast.Stmt {
	name := p.consume(scanner.IDENTIFIER, "Expect variable name.")
	var initializer ast.Expr
	if p.match(scanner.EQUAL) {
		initializer = p.expression()
	}
	p.consume(scanner.SEMICOLON, "Expect ';' after variable declaration.")
	return ast.NewVar(name, initializer)
}

func (p *Parser) expressionStatement() ast.Stmt {
	expr := p.expression()
	p.consume(scanner.SEMICOLON, "Expect ';' after expression.")
	return ast.NewExpression(expr)
}

func (p *Parser) assignment() ast.Expr {
	expr := p.equality()
	if p.match(scanner.EQUAL) {
		equals := p.previous()
		value := p.assignment()

		variable, ok := expr.(*ast.Variable)
		if !ok {
			panic(p.erro(equals, "Invalid assignment target."))
		}
		name := variable.Name
		return ast.NewAssign(name, value)
	}
	return expr
}

// Binary expressions.

func (p *Parser) equality() ast.Expr {
	expr := p.comparison()
	for p.match(scanner.BANGEQUAL, scanner.EQUALEQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = ast.NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) comparison() ast.Expr {
	expr := p.term()
	for p.match(scanner.GREATER, scanner.GREATEREQUAL, scanner.LESS, scanner.LESSEQUAL) {
		operator := p.previous()
		right := p.term()
		expr = ast.NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) term() ast.Expr {
	expr := p.factor()
	for p.match(scanner.MINUS, scanner.PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = ast.NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) factor() ast.Expr {
	expr := p.unary()
	for p.match(scanner.SLASH, scanner.STAR) {
		operator := p.previous()
		right := p.unary()
		expr = ast.NewBinary(expr, operator, right)
	}
	return expr
}

// Unary expression.

func (p *Parser) unary() ast.Expr {
	if p.match(scanner.BANG, scanner.MINUS) {
		operator := p.previous()
		right := p.unary()
		return ast.NewUnary(operator, right)
	}
	return p.primary()
}

// Primary expression.

func (p *Parser) primary() ast.Expr {
	if p.match(scanner.FALSE) {
		return ast.NewLiteral(false)
	}
	if p.match(scanner.TRUE) {
		return ast.NewLiteral(true)
	}
	if p.match(scanner.NIL) {
		return ast.NewLiteral(nil)
	}
	if p.match(scanner.NUMBER, scanner.STRING) {
		return ast.NewLiteral(p.previous().Literal())
	}
	if p.match(scanner.IDENTIFIER) {
		return ast.NewVariable(p.previous())
	}
	if p.match(scanner.LEFTPAREN) {
		expr := p.expression()
		p.consume(scanner.RIGHTPAREN, "Expect ')' after expression.")
		return ast.NewGrouping(expr)
	}
	panic(p.erro(p.peek(), "Expect expression."))
}

// helpers.

func (p *Parser) match(types ...scanner.TokenType) bool {
	for _, kind := range types {
		if p.check(kind) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(kind scanner.TokenType, errMsg string) scanner.Token {
	if p.check(kind) {
		return p.advance()
	}
	panic(p.erro(p.peek(), errMsg))
}

func (p *Parser) check(kind scanner.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Kind() == kind
}

func (p *Parser) advance() scanner.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Kind() == scanner.EOF
}

func (p *Parser) peek() scanner.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() scanner.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) erro(token scanner.Token, message string) error {
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
