package tools

import (
	"strings"
	"testing"
)

func Test_defineType(t *testing.T) {
	t.Run("Binary Expr", func(t *testing.T) {
		builder := &strings.Builder{}
		kind := "Binary   : Left Expr, Operator Token, Right Expr"

		defineType(builder, kind)

		want := `
type Binary struct {
	// Left field.
	Left Expr
	// Operator field.
	Operator Token
	// Right field.
	Right Expr
}

func NewBinary(left Expr, operator Token, right Expr) *Binary {
	this := Binary{}
	this.Left = left
	this.Operator = operator
	this.Right = right
	return &this
}

func (b *Binary) Accept(visitor Visitor) any {
	return visitor.VisitBinary(b)
}
`
		if got := builder.String(); got != want {
			t.Errorf("String() = %v, want %v", got, want)
		}
	})
}
