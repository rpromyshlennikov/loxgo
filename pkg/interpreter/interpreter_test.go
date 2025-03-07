package interpreter

import (
	"reflect"
	"testing"

	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/parser"
	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/plugins"
	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/scanner"
)

func TestInterpreter_Interpret(t *testing.T) {
	t.Run("Success all expressions", func(t *testing.T) {
		// arrange
		pprinter := plugins.NewAstPrinter()
		scnr := scanner.NewScanner("print !(3 / 2 + 2 * 4 - 1 > 5 == true != 4 <= 5 + (1 - -2));", nil)
		prsr := parser.NewParser(scnr.ScanTokens(), nil)
		parsed := prsr.Parse()
		interp := NewInterpreter()

		// act
		err := interp.Interpret(parsed)

		// assert
		if err != nil {
			t.Errorf("Interpret() return error: %s, but shouldn't", err)
		}
		want := "true"
		got := *interp.lastPrintedValue
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Interpret() = %v, want %v, ast %s", got, want, pprinter.Sprint(parsed))
		}
	})

	t.Run("Logical not operator works as expected", func(t *testing.T) {
		// arrange
		pprinter := plugins.NewAstPrinter()
		scnr := scanner.NewScanner("print !false == !!true == !nil == !!4;", nil)
		prsr := parser.NewParser(scnr.ScanTokens(), nil)
		parsed := prsr.Parse()
		interp := NewInterpreter()

		// act
		err := interp.Interpret(parsed)

		// assert
		if err != nil {
			t.Errorf("Interpret() return error: %s, but shouldn't", err)
		}
		want := "true"
		got := *interp.lastPrintedValue
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Interpret() = %v, want %v, ast %s", got, want, pprinter.Sprint(parsed))
		}
	})

	t.Run("Logical OR operator works just fine", func(t *testing.T) {
		tests := []struct {
			name    string
			sources string
			want    string
		}{
			{
				name: "Basic success non-lazy execution on booleans",
				sources: `
					var a = false or true;
					print a;
				`,
				want: "true",
			},
			{
				name: "Basic success lazy execution",
				sources: `
					var x = 1;
					var y = 0;
					x = x or (y=2);
					print x;
					// ensure that y = 0, but not 2; it will be 2 in case of lack of laziness. 
					if (y == 2) print y;
				`,
				want: "1",
			},
			{
				name: "Second expression evaluates when first is not truthy: nil value",
				sources: `
					var x = nil;
					var y = 0;
					x = x or (y=2);
					if (y == 2) print y;
				`,
				want: "2",
			},
			{
				name: "Second expression evaluates when first is not truthy: false value",
				sources: `
					var x = false;
					var y = 0;
					x = x or (y=2);
					if (y == 2) print y;
				`,
				want: "2",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// arrange
				pprinter := plugins.NewAstPrinter()
				scnr := scanner.NewScanner(tt.sources, nil)
				prsr := parser.NewParser(scnr.ScanTokens(), nil)
				parsed := prsr.Parse()
				interp := NewInterpreter()

				// act
				err := interp.Interpret(parsed)

				// assert
				if err != nil {
					t.Errorf("Interpret() return error: %s, but shouldn't", err)
				}

				got := *interp.lastPrintedValue
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Interpret() = %v, want %v, ast %s", got, tt.want, pprinter.Sprint(parsed))
				}
			})
		}
	})

	t.Run("Logical AND operator works just fine", func(t *testing.T) {
		tests := []struct {
			name    string
			sources string
			want    string
		}{
			{
				name: "Basic success non-lazy execution on booleans",
				sources: `
					var a = true and false;
					print a;
				`,
				want: "false",
			},
			{
				name: "Basic success lazy execution",
				sources: `
					var x = nil;
					var y = 0;
					x = x and (y=2);
					print x;
					// ensure that y = 0, but not 2; it will be 2 in case of lack of laziness. 
					if (y == 2) print y;
				`,
				want: "nil",
			},
			{
				name: "Second expression evaluates when first is truthy: int value",
				sources: `
					var x = 1;
					var y = 0;
					x = x and (y=2);
					if (y == 2) print y;
				`,
				want: "2",
			},
			{
				name: "Second expression evaluates when first is truthy: true value",
				sources: `
					var x = true;
					var y = 0;
					x = x and (y=2);
					if (y == 2) print y;
				`,
				want: "2",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// arrange
				pprinter := plugins.NewAstPrinter()
				scnr := scanner.NewScanner(tt.sources, nil)
				prsr := parser.NewParser(scnr.ScanTokens(), nil)
				parsed := prsr.Parse()
				interp := NewInterpreter()

				// act
				err := interp.Interpret(parsed)

				// assert
				if err != nil {
					t.Errorf("Interpret() return error: %s, but shouldn't", err)
				}

				got := *interp.lastPrintedValue
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Interpret() = %v, want %v, ast %s", got, tt.want, pprinter.Sprint(parsed))
				}
			})
		}
	})

	t.Run("Arithmetic operations works fine", func(t *testing.T) {
		// arrange
		pprinter := plugins.NewAstPrinter()
		scnr := scanner.NewScanner("print ((2 + 2 * 3) - -1) / -4;", nil)
		prsr := parser.NewParser(scnr.ScanTokens(), nil)
		parsed := prsr.Parse()
		interp := NewInterpreter()

		// act
		err := interp.Interpret(parsed)

		// assert
		if err != nil {
			t.Errorf("Interpret() return error: %s, but shouldn't", err)
		}
		want := "-2.25"
		got := *interp.lastPrintedValue
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Interpret() = %v, want %v, ast %s", got, want, pprinter.Sprint(parsed))
		}
	})

	t.Run("String concatenation works just fine", func(t *testing.T) {
		// arrange
		pprinter := plugins.NewAstPrinter()
		scnr := scanner.NewScanner(`print "foo" + "bar";`, nil)
		prsr := parser.NewParser(scnr.ScanTokens(), nil)
		parsed := prsr.Parse()
		interp := NewInterpreter()

		// act
		err := interp.Interpret(parsed)

		// assert
		if err != nil {
			t.Errorf("Interpret() return error: %s, but shouldn't", err)
		}
		want := "foobar"
		got := *interp.lastPrintedValue
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Interpret() = %v, want %v, ast %s", got, want, pprinter.Sprint(parsed))
		}
	})

	t.Run("Comparisons works just fine", func(t *testing.T) {
		// arrange
		pprinter := plugins.NewAstPrinter()
		scnr := scanner.NewScanner("print 2 < 4 == -1 <= 0 == 5 > 3 == 9 >= 9;", nil)
		prsr := parser.NewParser(scnr.ScanTokens(), nil)
		parsed := prsr.Parse()
		interp := NewInterpreter()

		// act
		err := interp.Interpret(parsed)

		// assert
		if err != nil {
			t.Errorf("Interpret() return error: %s, but shouldn't", err)
		}
		want := "true"
		got := *interp.lastPrintedValue
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Interpret() = %v, want %v, ast %s", got, want, pprinter.Sprint(parsed))
		}
	})

	t.Run("Comparisons between nils works also fine", func(t *testing.T) {
		// arrange
		pprinter := plugins.NewAstPrinter()
		scnr := scanner.NewScanner("print nil == nil;", nil)
		prsr := parser.NewParser(scnr.ScanTokens(), nil)
		parsed := prsr.Parse()
		interp := NewInterpreter()

		// act
		err := interp.Interpret(parsed)

		// assert
		if err != nil {
			t.Errorf("Interpret() return error: %s, but shouldn't", err)
		}
		want := "true"
		got := *interp.lastPrintedValue
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Interpret() = %v, want %v, ast %s", got, want, pprinter.Sprint(parsed))
		}
	})

	t.Run("Var declaration with initialization and access to it works fine", func(t *testing.T) {
		// arrange
		pprinter := plugins.NewAstPrinter()
		scnr := scanner.NewScanner("var x = 1+2; print x; x = x+5; print x;", nil)
		prsr := parser.NewParser(scnr.ScanTokens(), nil)
		parsed := prsr.Parse()
		interp := NewInterpreter()

		// act
		err := interp.Interpret(parsed)

		// assert
		if err != nil {
			t.Errorf("Interpret() return error: %s, but shouldn't", err)
		}
		want := "8"
		got := *interp.lastPrintedValue
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Interpret() = %v, want %v, ast %s", got, want, pprinter.Sprint(parsed))
		}
	})

	t.Run("Scoping and blocks works just fine", func(t *testing.T) {
		tests := []struct {
			name    string
			sources string
			want    string
		}{
			{
				name: "Three nested blocks shadowing works fine",
				sources: `
					var a = "global a";
					{
						var a = "outer a";
						{
							var a = "inner a";
							print a;
						}
					}
				`,
				want: "inner a",
			},
			{
				name: "The nested shadowing variables does not rewrite most global value",
				sources: `
					var a = "global a";
					{
						var a = "outer a";
						{
							var a = "inner a";
						}
					}
					print a;
				`,
				want: "global a",
			},
			{
				name: "Outer value can be set in nested blocks",
				sources: `
					var a = "global a";
					{
						var b = "middle b";
						{
							a = "inner a";
						}
					}
					print a;
				`,
				want: "inner a",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// arrange
				pprinter := plugins.NewAstPrinter()
				scnr := scanner.NewScanner(tt.sources, nil)
				prsr := parser.NewParser(scnr.ScanTokens(), nil)
				parsed := prsr.Parse()
				interp := NewInterpreter()

				// act
				err := interp.Interpret(parsed)

				// assert
				if err != nil {
					t.Errorf("Interpret() return error: %s, but shouldn't", err)
				}

				got := *interp.lastPrintedValue
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Interpret() = %v, want %v, ast %s", got, tt.want, pprinter.Sprint(parsed))
				}
			})
		}
	})

	t.Run("If statements works fine", func(t *testing.T) {
		tests := []struct {
			name    string
			sources string
			want    string
		}{
			{
				name: "truthy with else statement, going to then statement, ignoring else",
				sources: `
					var five = 5;
					var six = 6;
					if (five < six) {
						print "five < six"; // <- going here
					} else {
						print "five >= six";
					}`,
				want: "five < six",
			},
			{
				name: "not truthy with else statement, going to else statement, ignoring then",
				sources: `
					var five = 5;
					var six = 6;
					if (five > six) {
						print "five > six";
					} else {
						print "five <= six"; // <- going here
					}`,
				want: "five <= six",
			},
			{
				name: "truthy without else so going to then",
				sources: `
					var five = 5;
					var six = 6;
					if (five < six) {
						print "five < six";
					}`,
				want: "five < six",
			},
			{
				name: "not truthy without else so no execution at all",
				sources: `
					var five = 5;
					var six = 6;
					if (five > six) {
						print "five > six";
					}`,
				want: "",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// arrange
				pprinter := plugins.NewAstPrinter()
				scnr := scanner.NewScanner(
					tt.sources,
					nil,
				)
				prsr := parser.NewParser(scnr.ScanTokens(), nil)
				parsed := prsr.Parse()
				interp := NewInterpreter()

				// act
				err := interp.Interpret(parsed)

				// assert
				if err != nil {
					t.Errorf("Interpret() return error: %s, but shouldn't", err)
				}

				got := *interp.lastPrintedValue
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Interpret() = %v, want %v, ast %s", got, tt.want, pprinter.Sprint(parsed))
				}
			})
		}
	})

	t.Run("While statements works fine", func(t *testing.T) {
		tests := []struct {
			name    string
			sources string
			want    string
		}{
			{
				name: "success repeating while condition is truthy",
				sources: `
					var x = 0;
					while (x < 5) {
						x = x + 1;
					}
					print x;
					`,
				want: "5",
			},
			{
				name: "should not execute while condition is NOT truthy",
				sources: `
					var x = 1;
					while (x < 0) {
						x = x + 1;
					}
					print x;
					`,
				want: "1",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// arrange
				pprinter := plugins.NewAstPrinter()
				scnr := scanner.NewScanner(
					tt.sources,
					nil,
				)
				prsr := parser.NewParser(scnr.ScanTokens(), nil)
				parsed := prsr.Parse()
				interp := NewInterpreter()

				// act
				err := interp.Interpret(parsed)

				// assert
				if err != nil {
					t.Errorf("Interpret() return error: %s, but shouldn't", err)
				}

				got := *interp.lastPrintedValue
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Interpret() = %v, want %v, ast %s", got, tt.want, pprinter.Sprint(parsed))
				}
			})
		}
	})

	t.Run("For statements works fine", func(t *testing.T) {
		tests := []struct {
			name    string
			sources string
			want    string
		}{
			{
				name: "success repeating for condition is truthy",
				sources: `
					var x = 0;
					for (var i = 0; i < 5; i = i + 1) {
						x = i;
					}
					print x;
					`,
				want: "4",
			},
			{
				name: "success repeating for condition is truthy and there is no initialization",
				sources: `
					var x = 0;
					for (; x < 5; x = x + 1) {
						// do nothing in cycle
					}
					print x;
					`,
				want: "5",
			},
			{
				name: "success repeating for condition is truthy and there is expression initialization",
				sources: `
					var x = 0;
					var i = 0;
					for (i = 3; i < 5; i = i + 1) {
						x = x + i;
					}
					print x;
					`,
				want: "7",
			},
			{
				name: "success repeating for condition is truthy and there is no increment",
				sources: `
					var x = 0;
					for (var i = 0; i < 5;) {
						x = i;
						i = i + 1;
					}
					print x;
					`,
				want: "4",
			},
			{
				name: "success repeating for condition is truthy and there is no initialization and no increment",
				sources: `
					var x = 0;
					for (; x < 5;) {
						x = x + 1;
					}
					print x;
					`,
				want: "5",
			},
			{
				name: "should not execute if for condition is NOT truthy",
				sources: `
					var x = 1;
					for (var i = 0; i < 0; i = i + 1) {
						x = i;
					}
					print x;
					`,
				want: "1",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// arrange
				pprinter := plugins.NewAstPrinter()
				scnr := scanner.NewScanner(
					tt.sources,
					nil,
				)
				prsr := parser.NewParser(scnr.ScanTokens(), nil)
				parsed := prsr.Parse()
				interp := NewInterpreter()

				// act
				err := interp.Interpret(parsed)

				// assert
				if err != nil {
					t.Errorf("Interpret() return error: %s, but shouldn't", err)
				}

				got := *interp.lastPrintedValue
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Interpret() = %v, want %v, ast %s", got, tt.want, pprinter.Sprint(parsed))
				}
			})
		}
	})

	t.Run("Cannot interpret plus operator between string and number", func(t *testing.T) {
		// arrange
		pprinter := plugins.NewAstPrinter()
		scnr := scanner.NewScanner(`2 + "foo";`, nil)
		prsr := parser.NewParser(scnr.ScanTokens(), nil)
		parsed := prsr.Parse()
		interp := NewInterpreter()

		// act
		err := interp.Interpret(parsed)

		// assert
		wantErr := `Runtime error: "invalid type for operator '+' given, must be numbers or strings." at token: {PLUS + <nil> 1}`
		if err == nil {
			t.Errorf("Interpret() did not return error: %s, but should", wantErr)
		}
		if !reflect.DeepEqual(err.Error(), wantErr) {
			t.Errorf("Interpret() error = %s, want error %s, ast %s", err.Error(), wantErr, pprinter.Sprint(parsed))
		}
	})

	t.Run("Cannot interpret multiply operator between string and number", func(t *testing.T) {
		// arrange
		pprinter := plugins.NewAstPrinter()
		scnr := scanner.NewScanner(`2 * "foo";`, nil)
		prsr := parser.NewParser(scnr.ScanTokens(), nil)
		parsed := prsr.Parse()
		interp := NewInterpreter()

		// act
		err := interp.Interpret(parsed)

		// assert
		wantErr := `Runtime error: "invalid type for operator STAR given, must be number." at token: {STAR * <nil> 1}`
		if err == nil {
			t.Errorf("Interpret() did not return error: %s, but should", wantErr)
		}
		if !reflect.DeepEqual(err.Error(), wantErr) {
			t.Errorf("Interpret() error = %s, want error %s, ast %s", err.Error(), wantErr, pprinter.Sprint(parsed))
		}
	})

	t.Run("Cannot interpret no statements (empty sources)", func(t *testing.T) {
		// arrange
		pprinter := plugins.NewAstPrinter()
		scnr := scanner.NewScanner("", nil)
		prsr := parser.NewParser(scnr.ScanTokens(), nil)
		parsed := prsr.Parse()
		interp := NewInterpreter()

		// act
		err := interp.Interpret(parsed)

		// assert
		wantErr := `Runtime error: "no statements given" at token: {EOF  <nil> 0}`
		if err == nil {
			t.Errorf("Interpret() did not return error: %s, but should", wantErr)
		}
		if !reflect.DeepEqual(err.Error(), wantErr) {
			t.Errorf("Interpret() error = %s, want error %s, ast %s", err.Error(), wantErr, pprinter.Sprint(parsed))
		}
	})

}
