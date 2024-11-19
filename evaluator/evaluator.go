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
		return evalProgram(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBooltoBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return EvalPrefixExpr(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return EvalInfixExpr(node.Operator, left, right)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)
		return &object.ReturnValue{Value: val}
	case *ast.BlockStatement:
		return evalStatements(node.Statements)
	case *ast.IfExpression:
		return EvalIfExpression(node)

	}

	return nil
}

func evalProgram(statements []ast.Statement) object.Object {
	var obj object.Object

	for _, stmt := range statements {
		obj = Eval(stmt)
		if returnValue, ok := obj.(*object.ReturnValue); ok {
			return returnValue.Value
		}
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
	case "-":
		return evalPrefixMinusExpr(right)
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

func evalPrefixMinusExpr(obj object.Object) object.Object {
	if obj.Type() != object.INTEGER_OBJ {
		return NULL
	}
	value := obj.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func EvalInfixExpr(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBooltoBooleanObject(left == right)
	case operator == "!=":
		return nativeBooltoBooleanObject(left != right)
	default:
		return NULL
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer)
	rightVal := right.(*object.Integer)

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal.Value + rightVal.Value}
	case "-":
		return &object.Integer{Value: leftVal.Value - rightVal.Value}
	case "*":
		return &object.Integer{Value: leftVal.Value * rightVal.Value}
	case "/":
		if rightVal.Value != 0 {
			return &object.Integer{Value: leftVal.Value / rightVal.Value}
		}
		return NULL
	case "<":
		return nativeBooltoBooleanObject(leftVal.Value < rightVal.Value)
	case ">":
		return nativeBooltoBooleanObject(leftVal.Value > rightVal.Value)
	case "==":
		return nativeBooltoBooleanObject(leftVal.Value == rightVal.Value)
	case "!=":
		return nativeBooltoBooleanObject(leftVal.Value != rightVal.Value)
	default:
		return NULL
	}

}

func EvalIfExpression(ifExpression *ast.IfExpression) object.Object {
	condition := Eval(ifExpression.Condition)
	if isTruthy(condition) {
		return Eval(ifExpression.Then)
	} else if ifExpression.Alt != nil {
		return Eval(ifExpression.Alt)
	} else {
		return NULL
	}
}

func isTruthy(object object.Object) bool {
	switch object {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}
