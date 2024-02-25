package scanner

import (
	"reflect"
	"testing"

	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/errors"
)

type interprtrErr struct {
	line    int
	message string
}

func getErrorReporterStub() (*[]interprtrErr, errors.Reporter) {
	errs := &[]interprtrErr{}
	fn := func(line int, message string) {
		*errs = append(
			*errs,
			interprtrErr{
				line:    line,
				message: message,
			},
		)
	}
	return errs, fn
}

func TestScanner_ScanTokens(t *testing.T) {
	t.Run("test special symbols", func(t *testing.T) {
		sources :=
			`(
)
{
}
,
.
-
+
;
*
`
		s := NewScanner(sources, nil)
		want := NewScanner("", nil)
		want.tokens = []Token{
			NewToken(LEFTPAREN, "(", nil, 1),
			NewToken(RIGHTPAREN, ")", nil, 2),
			NewToken(LEFTBRACE, "{", nil, 3),
			NewToken(RIGHTBRACE, "}", nil, 4),
			NewToken(COMMA, ",", nil, 5),
			NewToken(DOT, ".", nil, 6),
			NewToken(MINUS, "-", nil, 7),
			NewToken(PLUS, "+", nil, 8),
			NewToken(SEMICOLON, ";", nil, 9),
			NewToken(STAR, "*", nil, 10),
			NewToken(EOF, "", nil, 11),
		}
		if got := s.ScanTokens(); !reflect.DeepEqual(got, want.tokens) {
			t.Errorf("ScanTokens() = %v, want %v", got, want.tokens)
		}
	})

	t.Run("test operators symbols", func(t *testing.T) {
		sources :=
			`!	=	<	>
!= == <= >=
`
		s := NewScanner(sources, nil)
		want := NewScanner("", nil)
		want.tokens = []Token{
			NewToken(BANG, "!", nil, 1),
			NewToken(EQUAL, "=", nil, 1),
			NewToken(LESS, "<", nil, 1),
			NewToken(GREATER, ">", nil, 1),
			NewToken(BANGEQUAL, "!=", nil, 2),
			NewToken(EQUALEQUAL, "==", nil, 2),
			NewToken(LESSEQUAL, "<=", nil, 2),
			NewToken(GREATEREQUAL, ">=", nil, 2),
			NewToken(EOF, "", nil, 3),
		}
		if got := s.ScanTokens(); !reflect.DeepEqual(got, want.tokens) {
			t.Errorf("ScanTokens() = %v, want %v", got, want.tokens)
		}
	})

	t.Run("test comments and whitespace symbols", func(t *testing.T) {
		sources :=
			`// this is a comment
	// tab

// there is \r in previous line
     // 5 spaces
`
		s := NewScanner(sources, nil)
		want := NewScanner("", nil)
		want.tokens = []Token{
			NewToken(EOF, "", nil, 6),
		}
		if got := s.ScanTokens(); !reflect.DeepEqual(got, want.tokens) {
			t.Errorf("ScanTokens() = %v, want %v", got, want.tokens)
		}
	})

	t.Run("test string literals", func(t *testing.T) {
		sources := `"some string"
"some other
multi-line string"
`
		s := NewScanner(sources, nil)
		want := NewScanner("", nil)
		want.tokens = []Token{
			NewToken(STRING, `"some string"`, "some string", 1),
			NewToken(STRING, `"some other
multi-line string"`, "some other\nmulti-line string", 3),
			NewToken(EOF, "", nil, 4),
		}
		if got := s.ScanTokens(); !reflect.DeepEqual(got, want.tokens) {
			t.Errorf("ScanTokens() = %v, want %v", got, want.tokens)
		}
		// TODO: add edge cases: non-terminated string
	})

	t.Run("test number literals", func(t *testing.T) {
		sources :=
			`1234567890
1234.056789
`
		s := NewScanner(sources, nil)
		want := NewScanner("", nil)
		want.tokens = []Token{
			NewToken(NUMBER, "1234567890", float64(1234567890), 1),
			NewToken(NUMBER, "1234.056789", 1234.056789, 2),
			NewToken(EOF, "", nil, 3),
		}
		if got := s.ScanTokens(); !reflect.DeepEqual(got, want.tokens) {
			t.Errorf("ScanTokens() = %v, want %v", got, want.tokens)
		}

		// TODO: add cases: .03434, 3434., 34234..3434 and so on
	})

	t.Run("test reserved words and identifiers symbols", func(t *testing.T) {
		sources :=
			`and
class
else
false
for
fun
if
nil
or
print
return
super
this
true
var
while
_some_user_defined_identifier
super_puper_var
`
		s := NewScanner(sources, nil)
		want := NewScanner("", nil)
		want.tokens = []Token{
			NewToken(AND, "and", nil, 1),
			NewToken(CLASS, "class", nil, 2),
			NewToken(ELSE, "else", nil, 3),
			NewToken(FALSE, "false", nil, 4),
			NewToken(FOR, "for", nil, 5),
			NewToken(FUN, "fun", nil, 6),
			NewToken(IF, "if", nil, 7),
			NewToken(NIL, "nil", nil, 8),
			NewToken(OR, "or", nil, 9),
			NewToken(PRINT, "print", nil, 10),
			NewToken(RETURN, "return", nil, 11),
			NewToken(SUPER, "super", nil, 12),
			NewToken(THIS, "this", nil, 13),
			NewToken(TRUE, "true", nil, 14),
			NewToken(VAR, "var", nil, 15),
			NewToken(WHILE, "while", nil, 16),
			NewToken(IDENTIFIER, "_some_user_defined_identifier", nil, 17),
			NewToken(IDENTIFIER, "super_puper_var", nil, 18),
			NewToken(EOF, "", nil, 19),
		}
		if got := s.ScanTokens(); !reflect.DeepEqual(got, want.tokens) {
			t.Errorf("ScanTokens() = %v, want %v", got, want.tokens)
		}

		// TODO: add
	})

	t.Run("test unknown symbols", func(t *testing.T) {
		sources := `@	#
 ^`
		savedErrors, reporter := getErrorReporterStub()
		s := NewScanner(sources, reporter)
		wantErrors := []interprtrErr{
			{
				line:    1,
				message: "Unexpected character: @.",
			},
			{
				line:    1,
				message: "Unexpected character: #.",
			},
			{
				line:    2,
				message: "Unexpected character: ^.",
			},
		}
		s.ScanTokens()
		if !reflect.DeepEqual(*savedErrors, wantErrors) {
			t.Errorf("ScanTokens() = %v, want %v", *savedErrors, wantErrors)
		}
	})
}
