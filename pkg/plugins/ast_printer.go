package plugins

import (
	"fmt"
	"strings"

	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/parser/ast"
)

type AstPrinter struct {
	results   *[]string
	currLevel *uint
}

func NewAstPrinter() AstPrinter {
	return AstPrinter{
		results:   new([]string),
		currLevel: new(uint),
	}
}

func (p AstPrinter) Sprint(stmts []ast.Stmt) string {
	if len(stmts) == 0 {
		return ""
	}
	for _, stmt := range stmts {
		stmt.Accept(p)
	}
	return strings.Join(*p.results, "\n")
}

func (p AstPrinter) VisitUnary(unary *ast.Unary) any {
	return p.parenthesize(unary.Operator.Lexeme(), unary.Right)
}

func (p AstPrinter) VisitVariable(variable *ast.Variable) any {
	if variable == nil {
		return "nil"
	}
	return fmt.Sprint(variable.Name.Lexeme())
}

func (p AstPrinter) VisitBinary(binary *ast.Binary) any {
	return p.parenthesize(binary.Operator.Lexeme(), binary.Left, binary.Right)
}

func (p AstPrinter) VisitLogical(logical *ast.Logical) any {
	return p.parenthesize(logical.Operator.Lexeme(), logical.Left, logical.Right)
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

func (p AstPrinter) VisitBlock(stmt *ast.Block) {
	p.addResult("{")
	*p.currLevel += 1
	for i := range stmt.Statements {
		stmt.Statements[i].Accept(p)
	}
	*p.currLevel -= 1
	p.addResult("}")
}

func (p AstPrinter) VisitExpression(stmt *ast.Expression) {
	value := stmt.Expression.Accept(p)
	p.addResult(value.(string) + ";")
}

func (p AstPrinter) VisitIf(stmt *ast.If) {
	value := stmt.Condition.Accept(p)
	result := "if (" + value.(string) + ") then"
	p.addResult(result)
	stmt.ThenBranch.Accept(p)
	if stmt.ElseBranch != nil {
		p.addResult("else")
		stmt.ThenBranch.Accept(p)
		p.addResult("")
	}
}

func (p AstPrinter) VisitPrint(stmt *ast.Print) {
	value := stmt.Expression.Accept(p)
	result := "print " + value.(string) + ";"
	p.addResult(result)
}

func (p AstPrinter) VisitVar(stmt *ast.Var) {
	value := stmt.Initializer.Accept(p)
	result := "var " + stmt.Name.Lexeme() + " = " + value.(string) + ";"
	p.addResult(result)
}

func (p AstPrinter) VisitAssign(expr *ast.Assign) any {
	value := expr.Value.Accept(p)
	result := expr.Name.Lexeme() + " = " + value.(string) + ";"

	return result
}

func (p AstPrinter) VisitWhile(stmt *ast.While) {
	value := stmt.Condition.Accept(p)
	result := "while (" + value.(string) + ") "
	p.addResult(result)
	stmt.Body.Accept(p)
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
	tabs := strings.Repeat("\t", int(*p.currLevel))
	*p.results = append(*p.results, tabs+result)
}
