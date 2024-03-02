package parser

import (
	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/scanner"
)

type Token = scanner.Token
type TokenType = scanner.TokenType

type Expr interface {
	accept(visitor Visitor) any
}

type Visitor interface {
	visitBinary(*Binary) any
	visitGrouping(*Grouping) any
	visitLiteral(*Literal) any
	visitUnary(*Unary) any
}

type Binary struct {
	// left field.
	left Expr
	// operator field.
	operator Token
	// right field.
	right Expr
}

func NewBinary(left Expr, operator Token, right Expr) *Binary {
	this := Binary{}
	this.left = left
	this.operator = operator
	this.right = right
	return &this
}

func (b *Binary) accept(visitor Visitor) any {
	return visitor.visitBinary(b)
}

type Grouping struct {
	// expression field.
	expression Expr
}

func NewGrouping(expression Expr) *Grouping {
	this := Grouping{}
	this.expression = expression
	return &this
}

func (g *Grouping) accept(visitor Visitor) any {
	return visitor.visitGrouping(g)
}

type Literal struct {
	// value field.
	value any
}

func NewLiteral(value any) *Literal {
	this := Literal{}
	this.value = value
	return &this
}

func (l *Literal) accept(visitor Visitor) any {
	return visitor.visitLiteral(l)
}

type Unary struct {
	// operator field.
	operator Token
	// right field.
	right Expr
}

func NewUnary(operator Token, right Expr) *Unary {
	this := Unary{}
	this.operator = operator
	this.right = right
	return &this
}

func (u *Unary) accept(visitor Visitor) any {
	return visitor.visitUnary(u)
}
