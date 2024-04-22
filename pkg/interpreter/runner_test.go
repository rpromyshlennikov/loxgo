package interpreter

import (
	"testing"
)

func TestLoxGo_Run(t *testing.T) {
	tests := []struct {
		name    string
		sources string
	}{
		{
			name:    "success",
			sources: "var x = 1+2;print x;",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lox := New()
			lox.Run(tt.sources)
		})
	}
}
