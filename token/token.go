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
	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"

	// Operators
	ASSIGN    = "="
	PLUS      = "+"
	MINUS     = "-"
	MULT      = "*"
	DIV       = "/"
	NOT       = "!"
	LT        = "<"
	GT        = ">"
	EQ        = "=="
	NOT_EQ    = "!="
	AMPERSAND = "&"

	// Delimeters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	PARENL   = "("
	PARENR   = ")"
	BRACEL   = "{"
	BRACER   = "}"
	BRACKETL = "["
	BRACKETR = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	WHILE    = "WHILE"
)

var keywords = map[string]TokenType{
	"fnc":    FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"while":  WHILE,
}

func FindKeywordOrIdent(keyword string) TokenType {
	if word, ok := keywords[keyword]; ok {
		return word
	}
	return IDENT
}
