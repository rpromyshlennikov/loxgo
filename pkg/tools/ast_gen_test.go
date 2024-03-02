package tools

import (
	"strings"
	"testing"
)

func Test_defineType(t *testing.T) {
	t.Run("Binary Expr", func(t *testing.T) {
		builder := &strings.Builder{}
		kind := "Binary   : left Expr, operator Token, right Expr"

		defineType(builder, kind)

		want := `
type Binary struct {
	// left field.
	left Expr
	// operator field.
	operator Token
	// right field.
	right Expr
}

func NewBinary(left Expr, operator Token, right Expr) *Binary {
	this := Binary{}
	this.left = left
	this.operator = operator
	this.right = right
	return &this
}

func (b *Binary) accept(visitor Visitor) any {
	return visitor.visitBinary(b)
}
`
		if got := builder.String(); got != want {
			t.Errorf("String() = %v, want %v", got, want)
		}
	})
}
