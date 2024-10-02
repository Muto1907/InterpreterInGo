package lexer

import (
	"testing"

	"../token"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},;`

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
		{token.BRACEL, "}"},
		{token.BRACER, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, "EOF"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - wrong literal. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - wrong literal. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
