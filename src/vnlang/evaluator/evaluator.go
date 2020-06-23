package evaluator

import (
	"fmt"
	"math/big"
	"vnlang/ast"
	"vnlang/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

var interrupt = false

func Interrupt() {
	interrupt = true
}

func ResetInterrupt() {
	interrupt = false
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	if interrupt {
		ResetInterrupt()
		return newError("Tiến trình bị ngắt")
	}

	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.BreakStatement:
		return &object.BreakSignal{}

	case *ast.ContinueStatement:
		return &object.ContinueSignal{}

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.WhileExpression:
		return evalWhileExpression(node, env)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

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

	case *ast.HashLiteral:
		return evalHashLiteral(node, env)

	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		case *object.BreakSignal:
			return newError("Không thể 'ngắt' ngoài vòng lặp")
		case *object.ContinueSignal:
			return newError("Không thể 'tiếp' ngoài vòng lặp")
		}

	}

	return result
}

func evalBlockStatement(
	block *ast.BlockStatement,
	env *object.Environment,
) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.BREAK_SIGNAL_OBJ || rt == object.CONTINUE_SIGNAL_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("toán tử lạ: %s%s", operator, right.Type())
	}
}

func evalInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftType := left.Type()
	rightType := right.Type()

	if leftType != rightType {
		return newError("kiểu không tương thích: %s %s %s", leftType, operator, rightType)
	}

	switch leftType {
	case object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case object.BOOLEAN_OBJ:
		return evalBooleanInfixExpression(operator, left, right)
	}

	switch operator {
	case "==":
		return nativeBoolToBooleanObject(left == right)
	case "!=":
		return nativeBoolToBooleanObject(left != right)
	}

	return newError("toán tử lạ: %s %s %s",
		leftType, operator, rightType)
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
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

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	value, ok := right.(*object.Integer)

	if !ok {
		return newError("toán tử lạ: -%s", right.Type())
	}

	var newValue big.Int
	newValue.Neg(value.Value)
	return &object.Integer{Value: &newValue}
}

func evalIntegerInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "<":
		return nativeBoolToBooleanObject(leftVal.Cmp(rightVal) < 0)
	case ">":
		return nativeBoolToBooleanObject(leftVal.Cmp(rightVal) > 0)
	case "<=":
		return nativeBoolToBooleanObject(leftVal.Cmp(rightVal) <= 0)
	case ">=":
		return nativeBoolToBooleanObject(leftVal.Cmp(rightVal) >= 0)
	case "==":
		return nativeBoolToBooleanObject(leftVal.Cmp(rightVal) == 0)
	case "!=":
		return nativeBoolToBooleanObject(leftVal.Cmp(rightVal) != 0)
	}
	var resVal big.Int
	switch operator {
	case "+":
		resVal.Add(leftVal, rightVal)
	case "-":
		resVal.Sub(leftVal, rightVal)
	case "*":
		resVal.Mul(leftVal, rightVal)
	case "/":
		resVal.Div(leftVal, rightVal)
	case "%":
		resVal.Mod(leftVal, rightVal)
	default:
		return newError("toán tử lạ: %s %s %s",
			left.Type(), operator, right.Type())
	}
	return &object.Integer{Value: &resVal}
}

func evalStringInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "==":
		return &object.Boolean{Value: leftVal == rightVal}
	case "!=":
		return &object.Boolean{Value: leftVal != rightVal}
	case "<":
		return &object.Boolean{Value: leftVal < rightVal}
	case ">":
		return &object.Boolean{Value: leftVal > rightVal}
	case "<=":
		return &object.Boolean{Value: leftVal <= rightVal}
	case ">=":
		return &object.Boolean{Value: leftVal >= rightVal}
	default:
		return newError("toán tử lạ: %s %s %s",
			left.Type(), operator, right.Type())
	}

}

func evalBooleanInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Boolean).Value
	rightVal := right.(*object.Boolean).Value

	switch operator {
	case "||":
		return &object.Boolean{Value: leftVal || rightVal}
	case "&&":
		return &object.Boolean{Value: leftVal && rightVal}
	case "==":
		return &object.Boolean{Value: leftVal == rightVal}
	case "!=":
		return &object.Boolean{Value: leftVal != rightVal}
	default:
		return newError("toán tử lạ: %s %s %s",
			left.Type(), operator, right.Type())
	}

}

func evalIfExpression(
	ie *ast.IfExpression,
	env *object.Environment,
) object.Object {
	for i, cond := range ie.Condition {
		condition := Eval(cond, env)
		if isError(condition) {
			return condition
		}

		if isTruthy(condition) {
			return Eval(ie.Consequence[i], env)
		}
	}

	if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalWhileExpression(
	ie *ast.WhileExpression,
	env *object.Environment,
) object.Object {

	var result object.Object

	for {
		condition := Eval(ie.Condition, env)
		if isError(condition) {
			return condition
		}

		if !isTruthy(condition) {
			break
		}

		result = Eval(ie.Body, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}

			if rt == object.BREAK_SIGNAL_OBJ {
				break
			}
		}
	}

	return unwrapWhileSignal(result)
}

func evalIdentifier(
	node *ast.Identifier,
	env *object.Environment,
) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	if node.Value == importKeyword {
		return &object.Import{Env: env}
	}

	return newError("không tìm thấy tên định danh: " + node.Value)
}

func isTruthy(obj object.Object) bool {
	switch obj {
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

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	_, ok := obj.(*object.Error)
	return ok
}

func evalExpressions(
	exps []ast.Expression,
	env *object.Environment,
) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {

	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		switch evaluated.(type) {
		case *object.BreakSignal:
			return newError("không thể ngắt ngoài vòng lặp")
		case *object.ContinueSignal:
			return newError("không thể ngắt tiếptiếp vòng lặp")

		}

		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return fn.Fn(args...)
	case *object.Import:

		return ImportFile(fn, args...)
	default:
		return newError("không phải là một hàm: %s", fn.Type())
	}
}

func extendFunctionEnv(
	fn *object.Function,
	args []object.Object,
) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func unwrapWhileSignal(obj object.Object) object.Object {
	switch obj.(type) {
	case *object.BreakSignal:
		return nil
	case *object.ContinueSignal:
		return nil
	default:
		return obj
	}
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("toán tử chỉ mục không hỗ trợ cho: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idxBig := index.(*object.Integer).Value
	if !idxBig.IsInt64() {
		return NULL
	}

	idx := idxBig.Int64()
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrayObject.Elements[idx]
}

func evalHashLiteral(
	node *ast.HashLiteral,
	env *object.Environment,
) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("không thể dùng như khóa băm: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("không thể dùng như khóa băm: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}
