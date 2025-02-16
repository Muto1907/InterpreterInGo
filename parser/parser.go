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
	INDEX
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.GT:       LESSGREATER,
	token.LT:       LESSGREATER,
	token.MINUS:    SUM,
	token.PLUS:     SUM,
	token.MULT:     PRODUCT,
	token.DIV:      PRODUCT,
	token.PARENL:   CALL,
	token.BRACKETL: INDEX,
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
	parser.addPrefixFnc(token.STRING, parser.parseStringLiteral)
	parser.addPrefixFnc(token.PARENL, parser.ParseGroupedExpr)
	parser.addPrefixFnc(token.IF, parser.parseIfExpression)
	parser.addPrefixFnc(token.FUNCTION, parser.parseFunctionLiteral)
	parser.addPrefixFnc(token.BRACKETL, parser.parseArrayLiteral)
	parser.addInfixFnc(token.PARENL, parser.parseCallExpression)
	parser.addInfixFnc(token.BRACKETL, parser.parseIndexExpr)
	parser.addPrefixFnc(token.BRACEL, parser.parseHashLiteral)
	parser.addPrefixFnc(token.AMPERSAND, parser.parsePrefixExpression)
	parser.addPrefixFnc(token.MULT, parser.parsePrefixExpression)
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
	case token.WHILE:
		return parser.parseWhileStatement()
	default:
		return parser.parseExpressionOrAssignmentStatement()
	}
}

func (parser *Parser) parseExpressionOrAssignmentStatement() ast.Statement {
	leftExp := parser.parseExpression(LOWEST)
	if leftExp == nil {
		return nil
	}

	if parser.peekTokenIs(token.ASSIGN) {
		parser.nextToken()

		assignStmt := &ast.AssignmentStatement{
			Token: parser.currToken,
			Left:  leftExp,
		}
		parser.nextToken()
		assignStmt.Value = parser.parseExpression(LOWEST)

		if parser.peekTokenIs(token.SEMICOLON) {
			parser.nextToken()
		}
		return assignStmt
	}

	stmt := &ast.ExpressionStatement{Token: parser.currToken}
	stmt.Expression = leftExp
	if parser.peekTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}
	return stmt
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
	parser.nextToken()
	stmt.Value = parser.parseExpression(LOWEST)
	if parser.peekTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}
	return stmt
}
func (parser *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: parser.currToken}
	parser.nextToken()
	stmt.ReturnValue = parser.parseExpression(LOWEST)
	if parser.peekTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}
	return stmt
}

func (parser *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{Token: parser.currToken}
	if !parser.expectPeek(token.PARENL) {
		return nil
	}
	parser.nextToken()
	stmt.Condition = parser.parseExpression(LOWEST)
	if !parser.expectPeek(token.PARENR) {
		return nil
	}

	if !parser.expectPeek(token.BRACEL) {
		return nil
	}
	stmt.Body = parser.parseBlockStatement()

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

func (parser *Parser) parseStringLiteral() ast.Expression {
	str := &ast.StringLiteral{Token: parser.currToken, Value: parser.currToken.Literal}
	return str
}

func (parser *Parser) parseArrayLiteral() ast.Expression {
	arr := &ast.ArrayLiteral{Token: parser.currToken}
	arr.Elements = parser.parseExpressionList(token.BRACKETR)
	return arr
}

func (parser *Parser) parseExpressionList(tok token.TokenType) []ast.Expression {
	res := []ast.Expression{}
	if parser.peekTokenIs(tok) {
		parser.nextToken()
		return res
	}
	parser.nextToken()
	res = append(res, parser.parseExpression(LOWEST))
	for parser.peekTokenIs(token.COMMA) {
		parser.nextToken()
		parser.nextToken()
		res = append(res, parser.parseExpression(LOWEST))
	}
	if !parser.expectPeek(tok) {
		return nil
	}
	return res
}

func (parser *Parser) parseIndexExpr(left ast.Expression) ast.Expression {
	indexExpr := &ast.IndexExpression{Token: parser.currToken, Left: left}
	parser.nextToken()
	indexExpr.Index = parser.parseExpression(LOWEST)
	if !parser.expectPeek(token.BRACKETR) {
		return nil
	}
	return indexExpr
}

func (parser *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: parser.currToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !parser.peekTokenIs(token.BRACER) {
		parser.nextToken()
		key := parser.parseExpression(LOWEST)

		if !parser.expectPeek(token.COLON) {
			return nil
		}

		parser.nextToken()
		val := parser.parseExpression(LOWEST)
		hash.Pairs[key] = val

		if !parser.peekTokenIs(token.BRACER) && !parser.expectPeek(token.COMMA) {
			return nil
		}
	}
	if !parser.expectPeek(token.BRACER) {
		return nil
	}
	return hash
}

func (parser *Parser) parseBoolean() ast.Expression {
	boolean := &ast.Boolean{Token: parser.currToken, Value: parser.currentTokenIs(token.TRUE)}
	return boolean
}

func (parser *Parser) ParseGroupedExpr() ast.Expression {
	parser.nextToken()

	expr := parser.parseExpression(LOWEST)

	if !parser.expectPeek(token.PARENR) {
		return nil
	}
	return expr
}

func (parser *Parser) parseIfExpression() ast.Expression {
	iff := &ast.IfExpression{
		Token: parser.currToken,
	}
	if !parser.expectPeek(token.PARENL) {
		return nil
	}

	parser.nextToken()
	iff.Condition = parser.parseExpression(LOWEST)
	if !parser.expectPeek(token.PARENR) {
		return nil
	}

	if !parser.expectPeek(token.BRACEL) {
		return nil
	}
	iff.Then = parser.parseBlockStatement()

	if parser.peekTokenIs(token.ELSE) {
		parser.nextToken()
		if !parser.expectPeek(token.BRACEL) {
			return nil
		}

		iff.Alt = parser.parseBlockStatement()
	}

	return iff

}

func (parser *Parser) parseFunctionLiteral() ast.Expression {
	fnc := &ast.FuncLiteral{Token: parser.currToken}
	if !parser.expectPeek(token.PARENL) {
		return nil
	}
	fnc.Parameters = parser.parseFunctionParameters()
	if !parser.expectPeek(token.BRACEL) {
		return nil
	}
	fnc.Body = parser.parseBlockStatement()
	fnc.Body.IsFunctionBody = true
	return fnc
}

func (parser *Parser) parseBlockStatement() *ast.BlockStatement {
	blck := &ast.BlockStatement{Token: parser.currToken, Statements: []ast.Statement{}}

	parser.nextToken()

	for !parser.currentTokenIs(token.BRACER) && !parser.currentTokenIs(token.EOF) {
		stmt := parser.parseStatement()
		if stmt != nil {
			blck.Statements = append(blck.Statements, stmt)
		}
		parser.nextToken()
	}
	return blck
}

func (parser *Parser) parseFunctionParameters() []*ast.Identifier {
	parameters := []*ast.Identifier{}
	if parser.peekTokenIs(token.PARENR) {
		parser.nextToken()
		return parameters
	}
	parser.nextToken()
	identifier := &ast.Identifier{Token: parser.currToken, Value: parser.currToken.Literal}
	parameters = append(parameters, identifier)
	for parser.peekTokenIs(token.COMMA) {
		parser.nextToken()
		parser.nextToken()
		identifier = &ast.Identifier{Token: parser.currToken, Value: parser.currToken.Literal}
		parameters = append(parameters, identifier)
	}
	if !parser.expectPeek(token.PARENR) {
		return nil
	}
	return parameters
}

func (parser *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	expr := &ast.CallExpression{Token: parser.currToken, Function: function}
	expr.Arguments = parser.parseExpressionList(token.PARENR)
	return expr
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
