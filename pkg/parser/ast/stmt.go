package ast

type Stmt interface {
	Accept(visitor VisitorStmt)
}

type VisitorStmt interface {
	VisitExpression(*Expression)
	VisitPrint(*Print)
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
