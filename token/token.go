package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	//For unknown tokens
	ILLEGAL = "ILLEGAL"
	//End of File
	EOF = "EOF"

	// Identifiers + literals
	IDENT = "IDENT"
	INT   = "INT"

	// Operators
	ASSIGN = "="
	PLUS   = "+"
	MINUS  = "-"
	MULT   = "*"
	DIV    = "/"

	// Delimeters
	COMMA     = ","
	SEMICOLON = ";"

	PARENL = "("
	PARENR = ")"
	BRACEL = "{"
	BRACER = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
)
