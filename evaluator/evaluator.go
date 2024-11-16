package evaluator

import (
	"github.com/Muto1907/interpreterInGo/ast"
	"github.com/Muto1907/interpreterInGo/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.NULL{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBooltoBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return EvalPrefixExpr(node.Operator, right)
	}
	return nil
}

func evalStatements(statements []ast.Statement) object.Object {
	var obj object.Object

	for _, stmt := range statements {
		obj = Eval(stmt)
	}
	return obj
}

func nativeBooltoBooleanObject(b bool) *object.Boolean {
	if b {
		return TRUE
	}
	return FALSE
}

func EvalPrefixExpr(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpr(right)
	default:
		return NULL
	}
}

func evalBangOperatorExpr(obj object.Object) object.Object {
	switch obj {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}
