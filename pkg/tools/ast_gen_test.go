package tools

import (
	"strings"
	"testing"
)

func Test_defineType(t *testing.T) {
	t.Run("Binary Expr", func(t *testing.T) {
		builder := &strings.Builder{}
		kind := "Binary   : Left Expr, Operator scanner.Token, Right Expr"

		defineType(builder, []string{"Expr", "any"}, kind)

		want := `
type Binary struct {
	// Left field.
	Left Expr
	// Operator field.
	Operator scanner.Token
	// Right field.
	Right Expr
}

func NewBinary(left Expr, operator scanner.Token, right Expr) *Binary {
	this := Binary{}
	this.Left = left
	this.Operator = operator
	this.Right = right
	return &this
}

func (b *Binary) Accept(visitor VisitorExpr) any {
	return visitor.VisitBinary(b)
}
`
		if got := builder.String(); got != want {
			t.Errorf("String() = %v, want %v", got, want)
		}
	})

	t.Run("Non-returnable visitor", func(t *testing.T) {
		builder := &strings.Builder{}
		kind := "SomeOp   : Operator scanner.Token, Right Literal"

		defineType(builder, []string{"Foo"}, kind)

		want := `
type SomeOp struct {
	// Operator field.
	Operator scanner.Token
	// Right field.
	Right Literal
}

func NewSomeOp(operator scanner.Token, right Literal) *SomeOp {
	this := SomeOp{}
	this.Operator = operator
	this.Right = right
	return &this
}

func (s *SomeOp) Accept(visitor VisitorFoo) {
	visitor.VisitSomeOp(s)
}
`
		if got := builder.String(); got != want {
			t.Errorf("String() = %v, want %v", got, want)
		}
	})
}

func Test_defineVisitor(t *testing.T) {
	t.Run("Expr interfaces", func(t *testing.T) {
		builder := &strings.Builder{}
		types := []string{
			"Binary   : Left Expr, Operator scanner.Token, Right Expr",
			"Grouping : Expression Expr",
			"Literal  : Value any",
			"Unary    : Operator scanner.Token, Right Expr",
		}

		defineVisitor(builder, []string{"Expr", "any"}, types)

		want := `
type Expr interface {
	Accept(visitor VisitorExpr) any
}

type VisitorExpr interface {
	VisitBinary(*Binary) any
	VisitGrouping(*Grouping) any
	VisitLiteral(*Literal) any
	VisitUnary(*Unary) any
}
`
		if got := builder.String(); got != want {
			t.Errorf("String() = %v, want %v", got, want)
		}
	})
}

func Test_defineImports(t *testing.T) {
	t.Run("Imports with scanner package", func(t *testing.T) {
		builder := &strings.Builder{}
		types := []string{
			"Binary   : Left Expr, Operator scanner.Token, Right Expr",
			"Grouping : Expression Expr",
			"Literal  : Value any",
			"Unary    : Operator scanner.Token, Right Expr",
		}

		defineImports(builder, []string{"Expr", "any"}, types)

		want := `
import (
	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/scanner"
)
`
		if got := builder.String(); got != want {
			t.Errorf("String() = %v, want %v", got, want)
		}
	})

	t.Run("Imports without scanner package", func(t *testing.T) {
		builder := &strings.Builder{}
		types := []string{
			"Grouping : Expression Expr",
			"Literal  : Value any",
		}

		defineImports(builder, []string{"Expr", "any"}, types)

		want := ""
		if got := builder.String(); got != want {
			t.Errorf("String() = %v, want %v", got, want)
		}
	})
}
