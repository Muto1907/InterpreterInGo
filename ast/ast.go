package ast

import (
	"bytes"
	"strings"

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

type AssignmentStatement struct {
	Token token.Token
	Value Expression
	Name  *Identifier
}

func (as *AssignmentStatement) statementNode() {}
func (as *AssignmentStatement) TokenLiteral() string {
	return as.Token.Literal
}
func (as *AssignmentStatement) String() string {
	var output bytes.Buffer
	output.WriteString(as.Name.String() + " = ")
	if as.Value != nil {
		output.WriteString(as.Value.String())
	}
	output.WriteString(";")
	return output.String()
}

type WhileStatement struct {
	Token     token.Token
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) statementNode() {}
func (ws *WhileStatement) TokenLiteral() string {
	return ws.Token.Literal
}

func (ws *WhileStatement) String() string {
	var output bytes.Buffer
	output.WriteString(ws.Token.Literal)
	output.WriteString(ws.Condition.String() + " ")
	output.WriteString(ws.Body.String())
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

type StringLiteral struct {
	Token token.Token
	Value string
}

func (str *StringLiteral) expressionNode() {}
func (str *StringLiteral) TokenLiteral() string {
	return str.Token.Literal
}

func (str *StringLiteral) String() string {
	return str.Token.Literal
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pref *PrefixExpression) expressionNode() {}
func (pref *PrefixExpression) TokenLiteral() string {
	return pref.Token.Literal
}
func (pref *PrefixExpression) String() string {
	var output bytes.Buffer
	output.WriteString("(" + pref.Operator)
	output.WriteString(pref.Right.String() + ")")
	return output.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (inf *InfixExpression) expressionNode() {}
func (inf *InfixExpression) TokenLiteral() string {
	return inf.Token.Literal
}
func (inf *InfixExpression) String() string {
	var output bytes.Buffer
	output.WriteString("(" + inf.Left.String())
	output.WriteString(" " + inf.Operator + " ")
	output.WriteString(inf.Right.String() + ")")
	return output.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (bool *Boolean) expressionNode() {}
func (bool *Boolean) TokenLiteral() string {
	return bool.Token.Literal
}
func (bool *Boolean) String() string {
	return bool.Token.Literal
}

type IfExpression struct {
	Token     token.Token
	Condition Expression
	Then      *BlockStatement
	Alt       *BlockStatement
}

func (iff *IfExpression) expressionNode() {}
func (iff *IfExpression) TokenLiteral() string {
	return iff.Token.Literal
}
func (iff *IfExpression) String() string {
	var output bytes.Buffer
	output.WriteString("if")
	output.WriteString(iff.Condition.String() + " ")
	output.WriteString(iff.Then.String())
	if iff.Alt != nil {
		output.WriteString(" else ")
		output.WriteString(iff.Alt.String())
	}
	return output.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (blck *BlockStatement) statementNode() {}
func (blck *BlockStatement) TokenLiteral() string {
	return blck.Token.Literal
}
func (blck *BlockStatement) String() string {
	var output bytes.Buffer
	for _, stmt := range blck.Statements {
		output.WriteString(stmt.String())
	}
	return output.String()
}

type FuncLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fn *FuncLiteral) expressionNode() {}
func (fn *FuncLiteral) TokenLiteral() string {
	return fn.Token.Literal
}
func (fn *FuncLiteral) String() string {
	var output bytes.Buffer
	params := []string{}
	for _, pa := range fn.Parameters {
		params = append(params, pa.String())
	}
	output.WriteString(fn.TokenLiteral() + "( ")
	output.WriteString(strings.Join(params, ", ") + ") ")
	output.WriteString(fn.Body.String())
	return output.String()
}

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}
func (ce *CallExpression) String() string {
	var output bytes.Buffer
	var args []string
	for _, arg := range ce.Arguments {
		args = append(args, arg.String())
	}
	output.WriteString(ce.Function.String() + "(")
	output.WriteString(strings.Join(args, ", "))
	output.WriteString(")")
	return output.String()
}

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}
func (al *ArrayLiteral) TokenLiteral() string {
	return al.Token.Literal
}
func (al *ArrayLiteral) String() string {
	var output bytes.Buffer
	var elements []string
	for _, arg := range al.Elements {
		elements = append(elements, arg.String())
	}
	output.WriteString("[")
	output.WriteString(strings.Join(elements, ", "))
	output.WriteString("]")
	return output.String()
}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}
func (ie *IndexExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *IndexExpression) String() string {
	var output bytes.Buffer
	output.WriteString("(")
	output.WriteString(ie.Left.String())
	output.WriteString("[")
	output.WriteString(ie.Index.String())
	output.WriteString("])")
	return output.String()
}

type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (hsh *HashLiteral) expressionNode() {}
func (hsh *HashLiteral) TokenLiteral() string {
	return hsh.Token.Literal
}
func (hsh *HashLiteral) String() string {
	var output bytes.Buffer
	pairs := []string{}
	for key, value := range hsh.Pairs {
		pairs = append(pairs, key.String()+": "+value.String())
	}
	output.WriteString("{")
	output.WriteString(strings.Join(pairs, ", "))
	output.WriteString("}")
	return output.String()
}
