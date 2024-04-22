package ast

import (
	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/scanner"
)

type Stmt interface {
	Accept(visitor VisitorStmt)
}

type VisitorStmt interface {
	VisitBlock(*Block)
	VisitExpression(*Expression)
	VisitPrint(*Print)
	VisitVar(*Var)
}

type Block struct {
	// Statements field.
	Statements []Stmt
}

func NewBlock(statements []Stmt) *Block {
	this := Block{}
	this.Statements = statements
	return &this
}

func (b *Block) Accept(visitor VisitorStmt) {
	visitor.VisitBlock(b)
}

type Expression struct {
	// Expression field.
	Expression Expr
}

func NewExpression(expression Expr) *Expression {
	this := Expression{}
	this.Expression = expression
	return &this
}

func (e *Expression) Accept(visitor VisitorStmt) {
	visitor.VisitExpression(e)
}

type Print struct {
	// Expression field.
	Expression Expr
}

func NewPrint(expression Expr) *Print {
	this := Print{}
	this.Expression = expression
	return &this
}

func (p *Print) Accept(visitor VisitorStmt) {
	visitor.VisitPrint(p)
}

type Var struct {
	// Name field.
	Name scanner.Token
	// Initializer field.
	Initializer Expr
}

func NewVar(name scanner.Token, initializer Expr) *Var {
	this := Var{}
	this.Name = name
	this.Initializer = initializer
	return &this
}

func (v *Var) Accept(visitor VisitorStmt) {
	visitor.VisitVar(v)
}
