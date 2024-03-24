package plugins

import (
	"fmt"
	"strings"

	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/parser/ast"
)

type AstPrinter struct {
	results *[]string
}

func NewAstPrinter() AstPrinter {
	return AstPrinter{
		results: new([]string),
	}
}

func (p AstPrinter) Sprint(stmts []ast.Stmt) string {
	if len(stmts) == 0 {
		return ""
	}
	for _, stmt := range stmts {
		stmt.Accept(p)
	}
	return strings.Join(*p.results, ";\n")
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

func (p AstPrinter) VisitExpression(stmt *ast.Expression) {
	value := stmt.Expression.Accept(p)
	p.addResult(value.(string) + ";")
}

func (p AstPrinter) VisitPrint(stmt *ast.Print) {
	value := stmt.Expression.Accept(p)
	result := "print " + value.(string) + ";"
	p.addResult(result)
}

func (p AstPrinter) parenthesize(name string, exprs ...ast.Expr) any {
	str := "(" + name
	for _, expr := range exprs {
		str += " " + expr.Accept(p).(string)
	}
	str += ")"
	return str
}

func (p AstPrinter) addResult(result string) {
	*p.results = append(*p.results, result)
}
