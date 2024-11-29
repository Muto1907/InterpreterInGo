package evaluator

import (
	"fmt"

	"github.com/Muto1907/interpreterInGo/ast"
	"github.com/Muto1907/interpreterInGo/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.NULL{}
)
var builtIns = map[string]*object.BuiltIn{
	"len": &object.BuiltIn{
		Fnc: func(args ...object.Object) object.Object {
			return NULL
		},
	},
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.Boolean:
		return nativeBooltoBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return EvalPrefixExpr(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return EvalInfixExpr(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FuncLiteral:
		return &object.Function{Params: node.Parameters, Body: node.Body, Env: env}
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return callFunction(function, args)
	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var obj object.Object

	for _, stmt := range program.Statements {
		obj = Eval(stmt, env)
		switch obj := obj.(type) {
		case *object.ReturnValue:
			return obj.Value
		case *object.Error:
			return obj
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
		return newError("unknown operator: %s%s", operator, right.Type())
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
		return newError("unknown operator: -%s", obj.Type())
	}
	value := obj.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func EvalInfixExpr(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBooltoBooleanObject(left == right)
	case operator == "!=":
		return nativeBooltoBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
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
		return newError("zero division: %d / %d", rightVal.Value, leftVal.Value)
	case "<":
		return nativeBooltoBooleanObject(leftVal.Value < rightVal.Value)
	case ">":
		return nativeBooltoBooleanObject(leftVal.Value > rightVal.Value)
	case "==":
		return nativeBooltoBooleanObject(leftVal.Value == rightVal.Value)
	case "!=":
		return nativeBooltoBooleanObject(leftVal.Value != rightVal.Value)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.String)
	rightVal := right.(*object.String)

	switch operator {
	case "+":
		return &object.String{Value: leftVal.Value + rightVal.Value}
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ifExpression *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ifExpression.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(ifExpression.Then, env)
	} else if ifExpression.Alt != nil {
		return Eval(ifExpression.Alt, env)
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

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var obj object.Object
	for _, stmt := range block.Statements {
		obj = Eval(stmt, env)
		if obj != nil {
			ot := obj.Type()
			if ot == object.RETURN_OBJ || ot == object.ERROR_OBJ {
				return obj
			}
		}
	}
	return obj
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if ok {
		return val
	}
	if builtin, ok := builtIns[node.Value]; ok {
		return builtin
	}
	return newError("identifier not found: %s", node.Value)

}

func evalExpressions(expressions []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, expr := range expressions {
		value := Eval(expr, env)
		if isError(value) {
			return []object.Object{value}
		}
		result = append(result, value)
	}
	return result
}

func callFunction(fnc object.Object, args []object.Object) object.Object {
	function, ok := fnc.(*object.Function)
	if !ok {
		return newError("not a Function %s", fnc.Type())
	}
	extendedEnv := extendFunctionEnvironment(function, args)
	value := Eval(function.Body, extendedEnv)
	return unwrapReturnValue(value)
}

func extendFunctionEnvironment(fnc *object.Function, args []object.Object) *object.Environment {
	env := object.NeweEnclosedEnvironment(fnc.Env)
	for paramId, param := range fnc.Params {
		env.Set(param.Value, args[paramId])
	}
	return env
}

func unwrapReturnValue(val object.Object) object.Object {
	if ret, ok := val.(*object.ReturnValue); ok {
		return ret.Value
	}
	return val
}
