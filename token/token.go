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
	NOT    = "!"
	LT     = "<"
	GT     = ">"

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
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

var keywords = map[string]TokenType{
	"func":   FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func FindKeywordOrIdent(keyword string) TokenType {
	if word, ok := keywords[keyword]; ok {
		return word
	}
	return IDENT
}
