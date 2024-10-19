package parser

import (
	"github.com/Muto1907/interpreterInGo/ast"
	"github.com/Muto1907/interpreterInGo/lexer"
	"github.com/Muto1907/interpreterInGo/token"
)

type Parser struct {
	l         *lexer.Lexer
	currToken token.Token
	peekToken token.Token
}

func New(lex *lexer.Lexer) *Parser {
	parser := &Parser{l: lex}
	parser.nextToken()
	parser.nextToken()
	return parser
}

func (parser *Parser) nextToken() {
	parser.currToken = parser.peekToken
	parser.peekToken = parser.l.NextToken()
}

func (parser *Parser) ParseProgram() *ast.Program {
	return nil
}
