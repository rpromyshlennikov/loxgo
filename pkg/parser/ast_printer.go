package parser

import (
	"fmt"
)

type AstPrinter struct{}

func (p AstPrinter) Sprint(expr Expr) string {
	return expr.accept(p).(string)
}

func (p AstPrinter) visitUnary(unary *Unary) any {
	return p.parenthesize(unary.operator.Lexeme(), unary.right)
}

func (p AstPrinter) visitBinary(binary *Binary) any {
	return p.parenthesize(binary.operator.Lexeme(), binary.left, binary.right)
}

func (p AstPrinter) visitLiteral(literal *Literal) any {
	if literal.value == nil {
		return "nil"
	}
	return fmt.Sprint(literal.value)
}

func (p AstPrinter) visitGrouping(grouping *Grouping) any {
	return p.parenthesize("group", grouping.expression)
}

func (p AstPrinter) parenthesize(name string, exprs ...Expr) any {
	str := "(" + name
	for _, expr := range exprs {
		str += " " + expr.accept(p).(string)
	}
	str += ")"
	return str
}
