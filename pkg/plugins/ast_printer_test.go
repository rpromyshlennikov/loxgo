package plugins

import (
	"testing"

	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/parser"
	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/parser/ast"
	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/scanner"
)

func TestAstPrinter_Sprint(t *testing.T) {
	t.Run("Check (* (- 123) (group 45.67)) expression pretty printing", func(t *testing.T) {
		expr := ast.NewBinary(
			ast.NewUnary(scanner.NewToken(scanner.MINUS, "-", nil, 1),
				ast.NewLiteral(123),
			),
			scanner.NewToken(scanner.STAR, "*", nil, 1),
			ast.NewGrouping(
				ast.NewLiteral(45.67)))

		p := AstPrinter{}
		want := "(* (- 123) (group 45.67))"
		if got := p.Sprint(expr); got != want {
			t.Errorf("Sprint() = %v, want %v", got, want)
		}
	})

	t.Run("Check -123 * (45.67)) as source pretty printing", func(t *testing.T) {
		expr := parser.NewParser(
			scanner.NewScanner(
				"-123 * (45.67)",
				nil,
			).ScanTokens(),
			nil,
		).Parse()

		p := AstPrinter{}
		want := "(* (- 123) (group 45.67))"
		if got := p.Sprint(expr); got != want {
			t.Errorf("Sprint() = %v, want %v", got, want)
		}
	})
}
