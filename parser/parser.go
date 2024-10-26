package parser

import (
	"fmt"

	"github.com/Muto1907/interpreterInGo/ast"
	"github.com/Muto1907/interpreterInGo/lexer"
	"github.com/Muto1907/interpreterInGo/token"
)

type (
	prefixParseFnc func() ast.Expression
	infixParseFnc  func(ast.Expression) ast.Expression
)

type Parser struct {
	l               *lexer.Lexer
	currToken       token.Token
	peekToken       token.Token
	errors          []string
	prefixParseFncs map[token.TokenType]prefixParseFnc
	infixParseFncs  map[token.TokenType]infixParseFnc
}

func New(lex *lexer.Lexer) *Parser {
	parser := &Parser{
		l:      lex,
		errors: []string{},
	}
	parser.nextToken()
	parser.nextToken()
	return parser
}

func (parser *Parser) Errors() []string {
	return parser.errors
}

func (parser *Parser) nextToken() {
	parser.currToken = parser.peekToken
	parser.peekToken = parser.l.NextToken()
}

func (parser *Parser) AddPrefixFnc(tokenType token.TokenType, fnc prefixParseFnc) {
	parser.prefixParseFncs[tokenType] = fnc
}

func (parser *Parser) AddInfixFnc(tokenType token.TokenType, fnc infixParseFnc) {
	parser.infixParseFncs[tokenType] = fnc
}

func (parser *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for parser.currToken.Type != token.EOF {
		stmt := parser.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		parser.nextToken()
	}
	return program
}

func (parser *Parser) parseStatement() ast.Statement {
	switch parser.currToken.Type {
	case token.LET:
		return parser.parseLetStatement()
	case token.RETURN:
		return parser.parseReturnStatement()
	default:
		return nil
	}
}

func (parser *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: parser.currToken}
	if !parser.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: parser.currToken, Value: parser.currToken.Literal}
	if !parser.expectPeek(token.ASSIGN) {
		return nil
	}
	for parser.currToken.Type != token.SEMICOLON {
		parser.nextToken()
	}
	return stmt
}
func (parser *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: parser.currToken}
	parser.nextToken()
	for parser.currToken.Type != token.SEMICOLON {
		parser.nextToken()
	}
	return stmt
}
func (parser *Parser) currentTokenIs(t token.TokenType) bool {
	return parser.currToken.Type == t
}

func (parser *Parser) peekTokenIs(t token.TokenType) bool {
	return parser.peekToken.Type == t
}
func (parser *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, parser.peekToken.Type)
	parser.errors = append(parser.errors, msg)
}
func (parser *Parser) expectPeek(t token.TokenType) bool {
	if parser.peekTokenIs(t) {
		parser.nextToken()
		return true
	}
	parser.peekError(t)
	return false
}
