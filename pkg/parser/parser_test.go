package parser

import (
	"reflect"
	"testing"

	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/scanner"
)

func TestParser_Parse(t *testing.T) {
	pprinter := AstPrinter{}

	t.Run("Success all expressions", func(t *testing.T) {
		scannr := scanner.NewScanner("3 / 2 + 2 * 4 - 1 > 5 == true != 4 <= 5 + (1 - 2)", nil)
		p := NewParser(scannr.ScanTokens(), nil)
		want := "(!= (== (> (- (+ (/ 3 2) (* 2 4)) 1) 5) true) (<= 4 (+ 5 (group (- 1 2)))))"
		if got := pprinter.Sprint(p.Parse()); !reflect.DeepEqual(got, want) {
			t.Errorf("Parse() = %v, want %v", got, want)
		}
	})

	t.Run("Success basic ungrouped math expression", func(t *testing.T) {
		scannr := scanner.NewScanner("3 / 2 + 2 * 4", nil)
		p := NewParser(scannr.ScanTokens(), nil)
		want := "(+ (/ 3 2) (* 2 4))"
		if got := pprinter.Sprint(p.Parse()); !reflect.DeepEqual(got, want) {
			t.Errorf("Parse() = %v, want %v", got, want)
		}
	})

	t.Run("Success basic grouped math expression", func(t *testing.T) {
		scannr := scanner.NewScanner("3 / (2 + 2) * 4", nil)
		p := NewParser(scannr.ScanTokens(), nil)
		want := "(* (/ 3 (group (+ 2 2))) 4)"
		if got := pprinter.Sprint(p.Parse()); !reflect.DeepEqual(got, want) {
			t.Errorf("Parse() = %v, want %v", got, want)
		}
	})

	t.Run("Success basic logical expression", func(t *testing.T) {
		scannr := scanner.NewScanner("false != true", nil)
		p := NewParser(scannr.ScanTokens(), nil)
		want := "(!= false true)"
		if got := pprinter.Sprint(p.Parse()); !reflect.DeepEqual(got, want) {
			t.Errorf("Parse() = %v, want %v", got, want)
		}
	})

	t.Run("Success basic logical expression with unary not", func(t *testing.T) {
		scannr := scanner.NewScanner("false == !true", nil)
		p := NewParser(scannr.ScanTokens(), nil)
		want := "(== false (! true))"
		if got := pprinter.Sprint(p.Parse()); !reflect.DeepEqual(got, want) {
			t.Errorf("Parse() = %v, want %v", got, want)
		}
	})

	t.Run("Success basic math expression with unary minus", func(t *testing.T) {
		scannr := scanner.NewScanner("5 - -4 == 9", nil)
		p := NewParser(scannr.ScanTokens(), nil)
		want := "(== (- 5 (- 4)) 9)"
		if got := pprinter.Sprint(p.Parse()); !reflect.DeepEqual(got, want) {
			t.Errorf("Parse() = %v, want %v", got, want)
		}
	})

	t.Run("Success basic logical expression with comparisons", func(t *testing.T) {
		scannr := scanner.NewScanner(`5 > 4 != "foo" < "bar"`, nil)
		p := NewParser(scannr.ScanTokens(), nil)
		want := "(!= (> 5 4) (< foo bar))"
		if got := pprinter.Sprint(p.Parse()); !reflect.DeepEqual(got, want) {
			t.Errorf("Parse() = %v, want %v", got, want)
		}
	})

	t.Run("Success basic logical expression with equal comparisons", func(t *testing.T) {
		scannr := scanner.NewScanner(`4 <= 5 == "foo" >= "bar"`, nil)
		p := NewParser(scannr.ScanTokens(), nil)
		want := "(== (<= 4 5) (>= foo bar))"
		if got := pprinter.Sprint(p.Parse()); !reflect.DeepEqual(got, want) {
			t.Errorf("Parse() = %v, want %v", got, want)
		}
	})

	t.Run("Success nil expression", func(t *testing.T) {
		scannr := scanner.NewScanner("nil != 5", nil)
		p := NewParser(scannr.ScanTokens(), nil)
		want := "(!= nil 5)"
		if got := pprinter.Sprint(p.Parse()); !reflect.DeepEqual(got, want) {
			t.Errorf("Parse() = %v, want %v", got, want)
		}
	})

	t.Run("Success math expression associativity on sum", func(t *testing.T) {
		scannr := scanner.NewScanner("3 + 2 + 5 + 4", nil)
		p := NewParser(scannr.ScanTokens(), nil)
		want := "(+ (+ (+ 3 2) 5) 4)"
		if got := pprinter.Sprint(p.Parse()); !reflect.DeepEqual(got, want) {
			t.Errorf("Parse() = %v, want %v", got, want)
		}
	})

	t.Run("Success math expression associativity on diff", func(t *testing.T) {
		scannr := scanner.NewScanner("3 - 2 - (5 - 4)", nil)
		p := NewParser(scannr.ScanTokens(), nil)
		want := "(- (- 3 2) (group (- 5 4)))"
		if got := pprinter.Sprint(p.Parse()); !reflect.DeepEqual(got, want) {
			t.Errorf("Parse() = %v, want %v", got, want)
		}
	})

	t.Run("Success math expression associativity on multiplication", func(t *testing.T) {
		scannr := scanner.NewScanner("3 * 2 * 5 * 4", nil)
		p := NewParser(scannr.ScanTokens(), nil)
		want := "(* (* (* 3 2) 5) 4)"
		if got := pprinter.Sprint(p.Parse()); !reflect.DeepEqual(got, want) {
			t.Errorf("Parse() = %v, want %v", got, want)
		}
	})

	t.Run("Success math expression associativity on division", func(t *testing.T) {
		scannr := scanner.NewScanner("3 / (2 / 5) / 4", nil)
		p := NewParser(scannr.ScanTokens(), nil)
		want := "(/ (/ 3 (group (/ 2 5))) 4)"
		if got := pprinter.Sprint(p.Parse()); !reflect.DeepEqual(got, want) {
			t.Errorf("Parse() = %v, want %v", got, want)
		}
	})
}
