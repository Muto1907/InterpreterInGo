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
	case *ast.WhileStatement:
		return evalWhileStatement(node, env)
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
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
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

func evalWhileStatement(while *ast.WhileStatement, env *object.Environment) object.Object {
	condition := Eval(while.Condition, env)
	if isError(condition) {
		return condition
	}
	for isTruthy(condition) {
		Eval(while.Body, env)
		condition = Eval(while.Condition, env)
	}
	return NULL

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
	switch fnc := fnc.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnvironment(fnc, args)
		value := Eval(fnc.Body, extendedEnv)
		return unwrapReturnValue(value)
	case *object.BuiltIn:
		return fnc.Fnc(args...)
	default:
		return newError("not a Function %s", fnc.Type())
	}

}

func extendFunctionEnvironment(fnc *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fnc.Env)
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

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpr(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpr(left, index)
	default:
		return newError("Index Operator not supported for %s", left.Type())
	}
}

func evalArrayIndexExpr(array, index object.Object) object.Object {
	arr := array.(*object.Array)
	ind := index.(*object.Integer).Value
	length := int64(len(arr.Elements) - 1)
	if length < ind || ind < 0 {
		return NULL
	}
	return arr.Elements[ind]
}

func evalHashIndexExpr(hash, index object.Object) object.Object {
	hashObj := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("%s can not be used as HashKey", index.Type())
	}
	pair, ok := hashObj.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}
	return pair.Value
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("%s can not be used as HashKey", key.Type())
		}

		val := Eval(valNode, env)
		if isError(val) {
			return val
		}
		hash := hashKey.HashKey()
		pairs[hash] = object.HashPair{Key: key, Value: val}
	}
	return &object.Hash{Pairs: pairs}
}
