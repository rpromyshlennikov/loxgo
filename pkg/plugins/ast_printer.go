package plugins

import (
	"fmt"

	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/parser/ast"
)

type AstPrinter struct{}

func (p AstPrinter) Sprint(expr ast.Expr) string {
	if expr == nil {
		return ""
	}
	return expr.Accept(p).(string)
}

func (p AstPrinter) VisitUnary(unary *ast.Unary) any {
	return p.parenthesize(unary.Operator.Lexeme(), unary.Right)
}

func (p AstPrinter) VisitBinary(binary *ast.Binary) any {
	return p.parenthesize(binary.Operator.Lexeme(), binary.Left, binary.Right)
}

func (p AstPrinter) VisitLiteral(literal *ast.Literal) any {
	if literal.Value == nil {
		return "nil"
	}
	return fmt.Sprint(literal.Value)
}

func (p AstPrinter) VisitGrouping(grouping *ast.Grouping) any {
	return p.parenthesize("group", grouping.Expression)
}

func (p AstPrinter) parenthesize(name string, exprs ...ast.Expr) any {
	str := "(" + name
	for _, expr := range exprs {
		str += " " + expr.Accept(p).(string)
	}
	str += ")"
	return str
}
