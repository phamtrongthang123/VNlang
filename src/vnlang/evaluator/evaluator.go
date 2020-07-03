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

func Eval(s object.CallStack, node ast.Node, env *object.Environment) object.Object {
	if interrupt {
		ResetInterrupt()
		return newError(s, node, "Tiến trình bị ngắt")
	}

	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(s, node, env)

	case *ast.BlockStatement:
		return evalBlockStatement(s, node, env)

	case *ast.ExpressionStatement:
		return Eval(s, node.Expression, env)

	case *ast.ReturnStatement:
		val := Eval(s, node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.BreakStatement:
		return &object.BreakSignal{}

	case *ast.ContinueStatement:
		return &object.ContinueSignal{}

	case *ast.LetStatement:
		val := Eval(s, node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		return evalPrefixExpression(s, node, env)

	case *ast.InfixExpression:
		return evalInfixExpression(s, node, env)

	case *ast.IfExpression:
		return evalIfExpression(s, node, env)

	case *ast.WhileExpression:
		return evalWhileExpression(s, node, env)

	case *ast.Identifier:
		return evalIdentifier(s, node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}

	case *ast.CallExpression:
		return applyFunction(s, node, env)

	case *ast.ArrayLiteral:
		elements := evalExpressions(s, node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.IndexExpression:
		return evalIndexExpression(s, node, env)

	case *ast.HashLiteral:
		return evalHashLiteral(s, node, env)

	}

	return nil
}

func evalProgram(s object.CallStack, program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(s, statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		case *object.BreakSignal:
			return newError(s, statement, "Không thể 'ngắt' ngoài vòng lặp")
		case *object.ContinueSignal:
			return newError(s, statement, "Không thể 'tiếp' ngoài vòng lặp")
		}

	}

	return result
}

func evalBlockStatement(
	s object.CallStack,
	block *ast.BlockStatement,
	env *object.Environment,
) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(s, statement, env)

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

func evalPrefixExpression(s object.CallStack, node *ast.PrefixExpression, env *object.Environment) object.Object {
	right := Eval(s, node.Right, env)
	if isError(right) {
		return right
	}

	switch node.Operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(s, node, right)
	default:
		return newError(s, node, "toán tử lạ: %s%s", node.Operator, right.Type())
	}
}

func evalInfixExpression(
	s object.CallStack,
	node *ast.InfixExpression,
	env *object.Environment,
) object.Object {
	left := Eval(s, node.Left, env)
	if isError(left) {
		return left
	}

	right := Eval(s, node.Right, env)
	if isError(right) {
		return right
	}

	leftType := left.Type()
	rightType := right.Type()

	if leftType != rightType {
		return newError(s, node, "kiểu không tương thích: %s %s %s", leftType, node.Operator, rightType)
	}

	switch leftType {
	case object.INTEGER_OBJ:
		return evalIntegerInfixExpression(s, node, left, right)
	case object.FLOAT_OBJ:
		return evalFloatInfixExpression(s, node, left, right)
	case object.STRING_OBJ:
		return evalStringInfixExpression(s, node, left, right)
	case object.BOOLEAN_OBJ:
		return evalBooleanInfixExpression(s, node, left, right)
	}

	switch node.Operator {
	case "==":
		return nativeBoolToBooleanObject(left == right)
	case "!=":
		return nativeBoolToBooleanObject(left != right)
	}

	return newError(s, node, "toán tử lạ: %s %s %s",
		leftType, node.Operator, rightType)
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

func evalMinusPrefixOperatorExpression(s object.CallStack, node *ast.PrefixExpression, right object.Object) object.Object {
	switch value := right.(type) {
	case *object.Integer:
		var newValue big.Int
		newValue.Neg(value.Value)
		return &object.Integer{Value: &newValue}
	case *object.Float:
		return &object.Float{Value: -value.Value}
	default:
		return newError(s, node, "toán tử lạ: -%s", right.Type())
	}
}

func evalIntegerInfixExpression(
	s object.CallStack,
	node *ast.InfixExpression,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch node.Operator {
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
	switch node.Operator {
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
		return newError(s, node, "toán tử lạ: %s %s %s",
			left.Type(), node.Operator, right.Type())
	}
	return &object.Integer{Value: &resVal}
}

func evalFloatInfixExpression(
	s object.CallStack,
	node *ast.InfixExpression,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value

	switch node.Operator {
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		return &object.Float{Value: leftVal / rightVal}
	default:
		return newError(s, node, "toán tử lạ: %s %s %s",
			left.Type(), node.Operator, right.Type())
	}
}

func evalStringInfixExpression(
	s object.CallStack,
	node *ast.InfixExpression,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch node.Operator {
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
		return newError(s, node, "toán tử lạ: %s %s %s",
			left.Type(), node.Operator, right.Type())
	}

}

func evalBooleanInfixExpression(
	s object.CallStack,
	node *ast.InfixExpression,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Boolean).Value
	rightVal := right.(*object.Boolean).Value

	switch node.Operator {
	case "||":
		return &object.Boolean{Value: leftVal || rightVal}
	case "&&":
		return &object.Boolean{Value: leftVal && rightVal}
	case "==":
		return &object.Boolean{Value: leftVal == rightVal}
	case "!=":
		return &object.Boolean{Value: leftVal != rightVal}
	default:
		return newError(s, node, "toán tử lạ: %s %s %s",
			left.Type(), node.Operator, right.Type())
	}

}

func evalIfExpression(
	s object.CallStack,
	ie *ast.IfExpression,
	env *object.Environment,
) object.Object {
	for i, cond := range ie.Condition {
		condition := Eval(s, cond, env)
		if isError(condition) {
			return condition
		}

		if isTruthy(condition) {
			return Eval(s, ie.Consequence[i], env)
		}
	}

	if ie.Alternative != nil {
		return Eval(s, ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalWhileExpression(
	s object.CallStack,
	ie *ast.WhileExpression,
	env *object.Environment,
) object.Object {

	var result object.Object

	for {
		condition := Eval(s, ie.Condition, env)
		if isError(condition) {
			return condition
		}

		if !isTruthy(condition) {
			break
		}

		result = Eval(s, ie.Body, env)

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
	s object.CallStack,
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

	return newError(s, node, "không tìm thấy tên định danh: "+node.Value)
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

func newError(s object.CallStack, node ast.Node, format string, a ...interface{}) *object.Error {
	return &object.Error{Stack: s, Pos: node.Position(), Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	_, ok := obj.(*object.Error)
	return ok
}

func evalExpressions(
	s object.CallStack,
	exps []ast.Expression,
	env *object.Environment,
) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(s, e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(s object.CallStack, node *ast.CallExpression, env *object.Environment) object.Object {
	fn := Eval(s, node.Function, env)
	if isError(fn) {
		return fn
	}

	args := evalExpressions(s, node.Arguments, env)
	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}

	newS := append(s, object.ActivationRecord{
		CallNode: node,
		Function: fn,
		Args:     args,
	})

	switch fn := fn.(type) {

	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(newS, fn.Body, extendedEnv)
		switch evaluated.(type) {
		case *object.BreakSignal:
			return newError(s, node, "không thể ngắt ngoài vòng lặp")
		case *object.ContinueSignal:
			return newError(s, node, "không thể tiếp ngoài vòng lặp")

		}

		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return fn.Fn(newS, node, args...)
	case *object.Import:

		return ImportFile(newS, node, fn, args...)
	default:
		return newError(s, node, "không phải là một hàm: %s", fn.Type())
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

func evalIndexExpression(s object.CallStack, node *ast.IndexExpression, env *object.Environment) object.Object {
	left := Eval(s, node.Left, env)
	if isError(left) {
		return left
	}
	index := Eval(s, node.Index, env)
	if isError(index) {
		return index
	}

	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(s, node, left, index)
	default:
		return newError(s, node, "toán tử chỉ mục không hỗ trợ cho: %s", left.Type())
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
	s object.CallStack,
	node *ast.HashLiteral,
	env *object.Environment,
) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(s, keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError(s, node, "không thể dùng như khóa băm: %s", key.Type())
		}

		value := Eval(s, valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func evalHashIndexExpression(s object.CallStack, node *ast.IndexExpression, hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError(s, node, "không thể dùng như khóa băm: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}
