package lexer

import (
	"testing"

	"github.com/Muto1907/interpreterInGo/token"
)

func TestNextToken(t *testing.T) {
	input := `=+-*/(){},;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.MULT, "*"},
		{token.DIV, "/"},
		{token.PARENL, "("},
		{token.PARENR, ")"},
		{token.BRACEL, "{"},
		{token.BRACER, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - wrong literal. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - wrong literal. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
