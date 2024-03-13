package parser

import (
	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/scanner"
)

type Token = scanner.Token
type TokenType = scanner.TokenType

type Expr interface {
	Accept(visitor Visitor) any
}

type Visitor interface {
	VisitBinary(*Binary) any
	VisitGrouping(*Grouping) any
	VisitLiteral(*Literal) any
	VisitUnary(*Unary) any
}

type Binary struct {
	// Left field.
	Left Expr
	// Operator field.
	Operator Token
	// Right field.
	Right Expr
}

func NewBinary(left Expr, operator Token, right Expr) *Binary {
	this := Binary{}
	this.Left = left
	this.Operator = operator
	this.Right = right
	return &this
}

func (b *Binary) Accept(visitor Visitor) any {
	return visitor.VisitBinary(b)
}

type Grouping struct {
	// Expression field.
	Expression Expr
}

func NewGrouping(expression Expr) *Grouping {
	this := Grouping{}
	this.Expression = expression
	return &this
}

func (g *Grouping) Accept(visitor Visitor) any {
	return visitor.VisitGrouping(g)
}

type Literal struct {
	// Value field.
	Value any
}

func NewLiteral(value any) *Literal {
	this := Literal{}
	this.Value = value
	return &this
}

func (l *Literal) Accept(visitor Visitor) any {
	return visitor.VisitLiteral(l)
}

type Unary struct {
	// Operator field.
	Operator Token
	// Right field.
	Right Expr
}

func NewUnary(operator Token, right Expr) *Unary {
	this := Unary{}
	this.Operator = operator
	this.Right = right
	return &this
}

func (u *Unary) Accept(visitor Visitor) any {
	return visitor.VisitUnary(u)
}
