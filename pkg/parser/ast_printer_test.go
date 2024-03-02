package parser

import (
	"testing"

	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/scanner"
)

func TestAstPrinter_Sprint(t *testing.T) {
	t.Run("Check (* (- 123) (group 45.67)) expression pretty printing", func(t *testing.T) {
		expr := NewBinary(
			NewUnary(scanner.NewToken(scanner.MINUS, "-", nil, 1),
				NewLiteral(123),
			),
			scanner.NewToken(scanner.STAR, "*", nil, 1),
			NewGrouping(
				NewLiteral(45.67)))

		p := AstPrinter{}
		want := "(* (- 123) (group 45.67))"
		if got := p.Sprint(expr); got != want {
			t.Errorf("Sprint() = %v, want %v", got, want)
		}
	})
}
