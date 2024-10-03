package lexer

import "github.com/Muto1907/interpreterInGo/token"

type Lexer struct {
	input            string
	currCharPosition int
	readPosition     int // 1 char after currCharPosition
	char             byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.readPosition]
	}
	l.currCharPosition = l.readPosition
	l.readPosition++
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	switch l.char {
	case '=':
		tok = newToken(token.ASSIGN, l.char)
	case ';':
		tok = newToken(token.SEMICOLON, l.char)
	case ',':
		tok = newToken(token.COMMA, l.char)
	case '+':
		tok = newToken(token.PLUS, l.char)
	case '-':
		tok = newToken(token.MINUS, l.char)
	case '*':
		tok = newToken(token.MULT, l.char)
	case '/':
		tok = newToken(token.DIV, l.char)
	case '(':
		tok = newToken(token.PARENL, l.char)
	case ')':
		tok = newToken(token.PARENR, l.char)
	case '{':
		tok = newToken(token.BRACEL, l.char)
	case '}':
		tok = newToken(token.BRACER, l.char)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.char) {
			tok.Literal = l.readIdentifier()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.char)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.currCharPosition
	for isLetter((l.char)) {
		l.readChar()
	}
	return l.input[position:l.currCharPosition]
}

func isLetter(char byte) bool {
	return ('A' <= char && char <= 'Z') || ('a' <= char && char <= 'z') || char == '_'
}

func newToken(tokenType token.TokenType, literal byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(literal),
	}
}
