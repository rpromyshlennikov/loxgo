package interpreter

import (
	"reflect"
	"testing"

	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/parser"
	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/plugins"
	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/scanner"
)

func TestInterpreter_Interpret(t *testing.T) {
	pprinter := plugins.AstPrinter{}

	t.Run("Success all expressions", func(t *testing.T) {
		// arrange
		scnr := scanner.NewScanner("!(3 / 2 + 2 * 4 - 1 > 5 == true != 4 <= 5 + (1 - -2))", nil)
		prsr := parser.NewParser(scnr.ScanTokens(), nil)
		expr := prsr.Parse()
		interp := Interpreter{}

		// act
		got, err := interp.Interpret(expr)

		// assert
		if err != nil {
			t.Errorf("Interpret() return error: %s, but shouldn't", err)
		}
		want := "true"
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Interpret() = %v, want %v, ast %s", got, want, pprinter.Sprint(expr))
		}
	})

	t.Run("Logical not operator works as expected", func(t *testing.T) {
		// arrange
		scnr := scanner.NewScanner("!false == !!true == !nil == !!4", nil)
		prsr := parser.NewParser(scnr.ScanTokens(), nil)
		expr := prsr.Parse()
		interp := Interpreter{}

		// act
		got, err := interp.Interpret(expr)

		// assert
		if err != nil {
			t.Errorf("Interpret() return error: %s, but shouldn't", err)
		}
		want := "true"
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Interpret() = %v, want %v, ast %s", got, want, pprinter.Sprint(expr))
		}
	})

	t.Run("Arithmetic operations works fine", func(t *testing.T) {
		// arrange
		scnr := scanner.NewScanner("((2 + 2 * 3) - -1) / -4", nil)
		prsr := parser.NewParser(scnr.ScanTokens(), nil)
		expr := prsr.Parse()
		interp := Interpreter{}

		// act
		got, err := interp.Interpret(expr)

		// assert
		if err != nil {
			t.Errorf("Interpret() return error: %s, but shouldn't", err)
		}
		want := "-2.25"
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Interpret() = %v, want %v, ast %s", got, want, pprinter.Sprint(expr))
		}
	})

	t.Run("String concatenation works just fine", func(t *testing.T) {
		// arrange
		scnr := scanner.NewScanner(`"foo" + "bar"`, nil)
		prsr := parser.NewParser(scnr.ScanTokens(), nil)
		expr := prsr.Parse()
		interp := Interpreter{}

		// act
		got, err := interp.Interpret(expr)

		// assert
		if err != nil {
			t.Errorf("Interpret() return error: %s, but shouldn't", err)
		}
		want := "foobar"
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Interpret() = %v, want %v, ast %s", got, want, pprinter.Sprint(expr))
		}
	})

	t.Run("Comparisons works just fine", func(t *testing.T) {
		// arrange
		scnr := scanner.NewScanner("2 < 4 == -1 <= 0 == 5 > 3 == 9 >= 9", nil)
		prsr := parser.NewParser(scnr.ScanTokens(), nil)
		expr := prsr.Parse()
		interp := Interpreter{}

		// act
		got, err := interp.Interpret(expr)

		// assert
		if err != nil {
			t.Errorf("Interpret() return error: %s, but shouldn't", err)
		}
		want := "true"
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Interpret() = %v, want %v, ast %s", got, want, pprinter.Sprint(expr))
		}
	})

	t.Run("Comparisons between nils works also fine", func(t *testing.T) {
		// arrange
		scnr := scanner.NewScanner("nil == nil", nil)
		prsr := parser.NewParser(scnr.ScanTokens(), nil)
		expr := prsr.Parse()
		interp := Interpreter{}

		// act
		got, err := interp.Interpret(expr)

		// assert
		if err != nil {
			t.Errorf("Interpret() return error: %s, but shouldn't", err)
		}
		want := "true"
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Interpret() = %v, want %v, ast %s", got, want, pprinter.Sprint(expr))
		}
	})

	t.Run("Cannot interpret plus operator between string and number", func(t *testing.T) {
		// arrange
		scnr := scanner.NewScanner(`2 + "foo"`, nil)
		prsr := parser.NewParser(scnr.ScanTokens(), nil)
		expr := prsr.Parse()
		interp := Interpreter{}

		// act
		_, err := interp.Interpret(expr)

		// assert
		wantErr := `Runtime error: "invalid type for operator '+' given, must be numbers or strings." at token: {PLUS + <nil> 1}`
		if err == nil {
			t.Errorf("Interpret() did not return error: %s, but should", wantErr)
		}
		if !reflect.DeepEqual(err.Error(), wantErr) {
			t.Errorf("Interpret() error = %s, want error %s, ast %s", err.Error(), wantErr, pprinter.Sprint(expr))
		}
	})

	t.Run("Cannot interpret multiply operator between string and number", func(t *testing.T) {
		// arrange
		scnr := scanner.NewScanner(`2 * "foo"`, nil)
		prsr := parser.NewParser(scnr.ScanTokens(), nil)
		expr := prsr.Parse()
		interp := Interpreter{}

		// act
		_, err := interp.Interpret(expr)

		// assert
		wantErr := `Runtime error: "invalid type for operator STAR given, must be number." at token: {STAR * <nil> 1}`
		if err == nil {
			t.Errorf("Interpret() did not return error: %s, but should", wantErr)
		}
		if !reflect.DeepEqual(err.Error(), wantErr) {
			t.Errorf("Interpret() error = %s, want error %s, ast %s", err.Error(), wantErr, pprinter.Sprint(expr))
		}
	})

	t.Run("Cannot interpret nil expression (empty sources)", func(t *testing.T) {
		// arrange
		scnr := scanner.NewScanner("", nil)
		prsr := parser.NewParser(scnr.ScanTokens(), nil)
		expr := prsr.Parse()
		interp := Interpreter{}

		// act
		_, err := interp.Interpret(expr)

		// assert
		wantErr := `Runtime error: "no expression given" at token: {EOF  <nil> 0}`
		if err == nil {
			t.Errorf("Interpret() did not return error: %s, but should", wantErr)
		}
		if !reflect.DeepEqual(err.Error(), wantErr) {
			t.Errorf("Interpret() error = %s, want error %s, ast %s", err.Error(), wantErr, pprinter.Sprint(expr))
		}
	})
}
