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

type Evaluator struct {
	Heap        object.Heap
	NextAddress uint64
	Threshold   int
	visitedEnvs map[*object.Environment]bool
}

func NewEval() *Evaluator {
	return &Evaluator{Heap: make(map[uint64]object.HeapObject), NextAddress: 0, Threshold: 10}
}

func (eva *Evaluator) MarkandSweep(env *object.Environment) {
	eva.visitedEnvs = make(map[*object.Environment]bool)
	eva.mark(env)
	eva.Sweep()
}

func (eva *Evaluator) mark(env *object.Environment) {
	if env == nil {
		return
	}
	if eva.visitedEnvs[env] {
		return
	}
	eva.visitedEnvs[env] = true

	for _, val := range env.State {
		eva.markValue(val)
	}
	eva.mark(env.Outer)
}

func (eva *Evaluator) markValue(obj object.Object) {
	switch o := obj.(type) {

	case *object.Pointer:
		eva.markObject(o.Value)

	case *object.Array:
		for _, elem := range o.Elements {
			eva.markValue(elem)
		}

	case *object.Hash:
		for _, pair := range o.Pairs {
			eva.markValue(pair.Key)
			eva.markValue(pair.Value)
		}

	case *object.Function:
		eva.mark(o.Env)
	}
}

func (eva *Evaluator) markObject(adress uint64) {
	heapObj, ok := eva.Heap[adress]
	if !ok {
		return
	}
	if heapObj.IsMarked {
		return
	}
	heapObj.IsMarked = true
	eva.Heap[adress] = heapObj

	eva.markValue(heapObj.Object)
}

func (eva *Evaluator) Sweep() {
	for addr, heapObj := range eva.Heap {
		if !heapObj.IsMarked {
			delete(eva.Heap, addr)
		} else {
			heapObj.IsMarked = false
			eva.Heap[addr] = heapObj
		}
	}
}

func (eva *Evaluator) Eval(node ast.Node, env *object.Environment) object.Object {
	if len(eva.Heap) >= eva.Threshold {
		eva.MarkandSweep(env)
	}
	switch node := node.(type) {
	case *ast.Program:
		return eva.evalProgram(node, env)
	case *ast.ExpressionStatement:
		return eva.Eval(node.Expression, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.Boolean:
		return nativeBooltoBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := eva.Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return eva.EvalPrefixExpr(node.Operator, right)
	case *ast.InfixExpression:
		left := eva.Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := eva.Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return EvalInfixExpr(node.Operator, left, right)
	case *ast.BlockStatement:
		return eva.evalBlockStatement(node, env)
	case *ast.IfExpression:
		return eva.evalIfExpression(node, env)
	case *ast.WhileStatement:
		val := eva.evalWhileStatement(node, env)
		if isError(val) {
			return val
		}
	case *ast.ReturnStatement:
		val := eva.Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := eva.Eval(node.Value, env)
		if isError(val) {
			return val
		}
		_, ok := env.GetLocal(node.Name.Value)
		if ok {
			return newError("Variable already initialized: %s", node.Name.Value)
		}
		env.Set(node.Name.Value, val)
	case *ast.AssignmentStatement:
		val := eva.evalAssignmentStatement(node, env)
		if isError(val) {
			return val
		}
	case *ast.Identifier:
		return eva.evalIdentifier(node, env)
	case *ast.FuncLiteral:
		return &object.Function{Params: node.Parameters, Body: node.Body, Env: env}
	case *ast.CallExpression:
		function := eva.Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := eva.evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return eva.callFunction(function, args)
	case *ast.ArrayLiteral:
		elements := eva.evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.HashLiteral:
		return eva.evalHashLiteral(node, env)
	case *ast.IndexExpression:
		left := eva.Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := eva.Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	}

	return nil
}

func (eva *Evaluator) evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var obj object.Object

	for _, stmt := range program.Statements {
		obj = eva.Eval(stmt, env)
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

func (eva *Evaluator) EvalPrefixExpr(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpr(right)
	case "-":
		return evalPrefixMinusExpr(right)
	case "&":
		return eva.evalAmpersandExpr(right)
	case "*":
		return eva.evalDereference(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func (eva *Evaluator) evalAmpersandExpr(obj object.Object) object.Object {
	eva.Heap[eva.NextAddress] = object.NewHeapOject(obj)
	ptr := &object.Pointer{Value: eva.NextAddress}
	eva.NextAddress += 1
	return ptr
}

func (eva *Evaluator) evalDereference(obj object.Object) object.Object {
	if obj.Type() != object.POINTER_OBJ {
		return newError("unknown operator: *%s", obj.Type())
	}
	ptr := obj.(*object.Pointer)
	ret := eva.Heap[ptr.Value].Object
	return ret
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

func (eva *Evaluator) evalIfExpression(ifExpression *ast.IfExpression, env *object.Environment) object.Object {
	condition := eva.Eval(ifExpression.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return eva.Eval(ifExpression.Then, env)
	} else if ifExpression.Alt != nil {
		return eva.Eval(ifExpression.Alt, env)
	} else {
		return NULL
	}
}

func (eva *Evaluator) evalWhileStatement(while *ast.WhileStatement, env *object.Environment) object.Object {
	condition := eva.Eval(while.Condition, env)
	if isError(condition) {
		return condition
	}

	for isTruthy(condition) {
		val := eva.Eval(while.Body, env)
		if val != nil {
			if isError(val) {
				return val
			}
			if val.Type() == object.RETURN_OBJ {
				return val
			}
		}
		condition = eva.Eval(while.Condition, env)
		if isError(condition) {
			return condition
		}
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

func (eva *Evaluator) evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var obj object.Object
	var blockEnv *object.Environment
	if block.IsFunctionBody {
		blockEnv = env
	} else {
		blockEnv = object.NewEnclosedEnvironment(env)
	}
	for _, stmt := range block.Statements {
		obj = eva.Eval(stmt, blockEnv)
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

func (eva *Evaluator) evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if ok {
		return val
	}
	if builtin, ok := builtIns[node.Value]; ok {
		return builtin
	}
	return newError("identifier not found: %s", node.Value)
}

func (eva *Evaluator) evalExpressions(expressions []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, expr := range expressions {
		value := eva.Eval(expr, env)
		if isError(value) {
			return []object.Object{value}
		}
		result = append(result, value)
	}
	return result
}

func (eva *Evaluator) callFunction(fnc object.Object, args []object.Object) object.Object {
	switch fnc := fnc.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnvironment(fnc, args)
		value := eva.Eval(fnc.Body, extendedEnv)
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

func (eva *Evaluator) evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valNode := range node.Pairs {
		key := eva.Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("%s can not be used as HashKey", key.Type())
		}

		val := eva.Eval(valNode, env)
		if isError(val) {
			return val
		}
		hash := hashKey.HashKey()
		pairs[hash] = object.HashPair{Key: key, Value: val}
	}
	return &object.Hash{Pairs: pairs}
}

func (eva *Evaluator) evalAssignmentStatement(stmt *ast.AssignmentStatement, env *object.Environment) object.Object {
	val := eva.Eval(stmt.Value, env)
	if isError(val) {
		return val
	}

	switch left := stmt.Left.(type) {
	case *ast.Identifier:
		_, localOk := env.GetLocal(left.Value)
		if localOk {
			env.Set(left.Value, val)
			return val
		}
		_, ok := env.Get(left.Value)
		if !ok {
			return newError("Variable not initialized: %s", left.Value)
		}
		env.SetOuter(left.Value, val)
		return val

	case *ast.PrefixExpression:
		if left.Operator == "*" {
			pointerObj := eva.Eval(left.Right, env)
			if isError(pointerObj) {
				return pointerObj
			}
			if pointerObj.Type() != object.POINTER_OBJ {
				return newError("cannot assign through non-pointer type: %s", pointerObj.Type())
			}

			ptr := pointerObj.(*object.Pointer)

			eva.Heap[ptr.Value] = object.NewHeapOject(val)
			return val
		}

		return newError("unsupported prefix operator in assignment: %s", left.Operator)
	case *ast.IndexExpression:
		arrayObj := eva.Eval(left.Left, env)
		if isError(arrayObj) {
			return arrayObj
		}

		indexObj := eva.Eval(left.Index, env)
		if isError(indexObj) {
			return indexObj
		}
		arr, ok := arrayObj.(*object.Array)
		if ok {
			intIdx, ok := indexObj.(*object.Integer)
			if !ok {
				return newError("array index is not an integer: %s", indexObj.Type())
			}

			idx := intIdx.Value
			if idx < 0 || idx >= int64(len(arr.Elements)) {
				return newError("array index out of bounds: %d", idx)
			}

			arr.Elements[idx] = val
			return val
		}

		return newError("index assignment not supported for %s", arrayObj.Type())

	default:
		return newError("invalid assignment target: %T", left)
	}

}
