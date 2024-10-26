package ast

import (
	"bytes"

	"github.com/Muto1907/interpreterInGo/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Expression interface {
	Node
	expressionNode()
}

type Statement interface {
	Node
	statementNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var output bytes.Buffer
	for _, stmt := range p.Statements {
		output.WriteString(stmt.String())
	}
	return output.String()
}

type LetStatement struct {
	Token token.Token
	Value Expression
	Name  *Identifier
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *LetStatement) String() string {
	var output bytes.Buffer
	output.WriteString(ls.TokenLiteral() + " ")
	output.WriteString(ls.Name.String())
	output.WriteString(" = ")
	if ls.Value != nil {
		output.WriteString(ls.Value.String())
	}
	output.WriteString(";")
	return output.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (id *Identifier) expressionNode() {}
func (id *Identifier) TokenLiteral() string {
	return id.Token.Literal
}
func (id *Identifier) String() string {
	return id.Value
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (inte *IntegerLiteral) expressionNode() {}
func (inte *IntegerLiteral) TokenLiteral() string {
	return inte.Token.Literal
}
func (inte *IntegerLiteral) String() string {
	return inte.Token.Literal
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (ret *ReturnStatement) statementNode() {}
func (ret *ReturnStatement) TokenLiteral() string {
	return ret.Token.Literal
}

func (ret *ReturnStatement) String() string {
	var output bytes.Buffer
	output.WriteString(ret.TokenLiteral() + " ")
	if ret.ReturnValue != nil {
		output.WriteString(ret.ReturnValue.String())
	}
	output.WriteString(";")
	return output.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (expr *ExpressionStatement) statementNode() {}
func (expr *ExpressionStatement) TokenLiteral() string {
	return expr.Token.Literal
}
func (expr *ExpressionStatement) String() string {
	if expr.Expression != nil {
		return expr.Expression.String()
	}
	return ""
}
