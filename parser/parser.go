package parser

import (
	"fmt"
	"strconv"

	"github.com/Muto1907/interpreterInGo/ast"
	"github.com/Muto1907/interpreterInGo/lexer"
	"github.com/Muto1907/interpreterInGo/token"
)

type (
	prefixParseFnc func() ast.Expression
	infixParseFnc  func(ast.Expression) ast.Expression
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precedences = map[token.TokenType]int{
	token.EQ:     EQUALS,
	token.NOT_EQ: EQUALS,
	token.GT:     LESSGREATER,
	token.LT:     LESSGREATER,
	token.MINUS:  SUM,
	token.PLUS:   SUM,
	token.MULT:   PRODUCT,
	token.DIV:    PRODUCT,
}

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
		l:               lex,
		errors:          []string{},
		prefixParseFncs: make(map[token.TokenType]prefixParseFnc),
		infixParseFncs:  make(map[token.TokenType]infixParseFnc),
	}
	parser.addPrefixFnc(token.IDENT, parser.parseIdentifier)
	parser.addPrefixFnc(token.INT, parser.parseIntegerLiteral)
	parser.addPrefixFnc(token.MINUS, parser.parsePrefixExpression)
	parser.addPrefixFnc(token.NOT, parser.parsePrefixExpression)
	for _, tok := range []token.TokenType{token.PLUS, token.MINUS, token.DIV, token.MULT, token.EQ, token.NOT_EQ, token.LT, token.GT} {
		parser.addInfixFnc(tok, parser.parseInfixExpression)
	}
	parser.addPrefixFnc(token.TRUE, parser.parseBoolean)
	parser.addPrefixFnc(token.FALSE, parser.parseBoolean)
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

func (parser *Parser) addPrefixFnc(tokenType token.TokenType, fnc prefixParseFnc) {
	parser.prefixParseFncs[tokenType] = fnc
}

func (parser *Parser) addInfixFnc(tokenType token.TokenType, fnc infixParseFnc) {
	parser.infixParseFncs[tokenType] = fnc
}

func (parser *Parser) peekPrecedence() int {
	if prec, ok := precedences[parser.peekToken.Type]; ok {
		return prec
	}
	return LOWEST
}

func (parser *Parser) currentPrecedence() int {
	if prec, ok := precedences[parser.currToken.Type]; ok {
		return prec
	}
	return LOWEST
}

func (parser *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for parser.currToken.Type != token.EOF {
		stmt := parser.parseStatement()
		program.Statements = append(program.Statements, stmt)
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
		return parser.parseExpressionStatement()
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

func (parser *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: parser.currToken}
	stmt.Expression = parser.parseExpression(LOWEST)
	if parser.peekTokenIs(token.SEMICOLON) {
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

func (parser *Parser) parseExpression(precedence int) ast.Expression {
	prefix := parser.prefixParseFncs[parser.currToken.Type]
	if prefix == nil {
		parser.noPrefixParseFuncFoundError(parser.currToken.Type)
		return nil
	}
	leftExp := prefix()
	for !parser.peekTokenIs(token.SEMICOLON) && precedence < parser.peekPrecedence() {
		infix := parser.infixParseFncs[parser.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		parser.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (parser *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: parser.currToken, Value: parser.currToken.Literal}
}

func (parser *Parser) parseIntegerLiteral() ast.Expression {
	inte := &ast.IntegerLiteral{Token: parser.currToken}

	val, err := strconv.ParseInt(parser.currToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as Integer", parser.currToken.Literal)
		parser.errors = append(parser.errors, msg)
		return nil
	}
	inte.Value = val
	return inte
}

func (parser *Parser) parseBoolean() ast.Expression {
	boolean := &ast.Boolean{Token: parser.currToken, Value: parser.currentTokenIs(token.TRUE)}
	return boolean
}

func (parser *Parser) parsePrefixExpression() ast.Expression {
	pref := &ast.PrefixExpression{
		Token:    parser.currToken,
		Operator: parser.currToken.Literal,
	}

	parser.nextToken()

	pref.Right = parser.parseExpression(PREFIX)
	return pref
}

func (parser *Parser) noPrefixParseFuncFoundError(ttype token.TokenType) {
	msg := fmt.Sprintf("No Prefix Parse Function found for %s", ttype)
	parser.errors = append(parser.errors, msg)
}

func (parser *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	inf := &ast.InfixExpression{
		Token:    parser.currToken,
		Operator: parser.currToken.Literal,
		Left:     left,
	}
	precedence := parser.currentPrecedence()
	parser.nextToken()
	inf.Right = parser.parseExpression(precedence)

	return inf
}
