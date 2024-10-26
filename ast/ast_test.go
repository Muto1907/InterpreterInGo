package ast

import (
	"testing"

	"github.com/Muto1907/interpreterInGo/token"
)

func TestString(t *testing.T) {
	Program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{
					Type:    token.LET,
					Literal: "let",
				},
				Name: &Identifier{
					Token: token.Token{
						Type:    token.IDENT,
						Literal: "myVar",
					},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{
						Type:    token.IDENT,
						Literal: "anotherVar",
					},
					Value: "anotherVar",
				},
			},
		},
	}
	if Program.String() != "let myVar = anotherVar;" {
		t.Errorf("program.String() error, got=%q", Program.String())
	}
}
