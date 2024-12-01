package lexer

import (
	"github.com/Muto1907/interpreterInGo/token"
)

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

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.nomWhitespace()

	switch l.char {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: "=="}
		} else {
			tok = newToken(token.ASSIGN, l.char)
		}
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
	case '<':
		tok = newToken(token.LT, l.char)
	case '>':
		tok = newToken(token.GT, l.char)
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: "!="}
			l.readChar()
		} else {
			tok = newToken(token.NOT, l.char)
		}
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case '[':
		tok = newToken(token.BRACKETL, l.char)
	case ']':
		tok = newToken(token.BRACKETR, l.char)
	case ':':
		tok = newToken(token.COLON, l.char)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.char) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.FindKeywordOrIdent(tok.Literal)
			return tok
		} else if isDigit(l.char) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
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

func (l *Lexer) readNumber() string {
	startPosition := l.currCharPosition
	for isDigit(l.char) {
		l.readChar()
	}
	return l.input[startPosition:l.currCharPosition]
}

func (l *Lexer) readString() string {
	startPosition := l.currCharPosition + 1
	for {
		l.readChar()
		if l.char == '"' || l.char == 0 {
			break
		}
	}
	return l.input[startPosition:l.currCharPosition]

}

func isLetter(char byte) bool {
	return ('A' <= char && char <= 'Z') || ('a' <= char && char <= 'z') || char == '_'
}

func isDigit(char byte) bool {
	return ('0' <= char && char <= '9')
}

func (l *Lexer) nomWhitespace() {
	for l.char == ' ' || l.char == '\n' || l.char == '\t' || l.char == '\r' {
		l.readChar()
	}
}

func newToken(tokenType token.TokenType, literal byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(literal),
	}
}
