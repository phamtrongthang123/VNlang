package evaluator

import (
	"fmt"
	"vnlang/ast"
	"vnlang/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

type Evaluator struct {
	interrupt bool
	Builtin   object.BuiltinFnMap
	Stack     object.CallStack
	Env       *object.Environment
}

func New(Builtin object.BuiltinFnMap) *Evaluator {
	return &Evaluator{
		Builtin: Builtin,
		Stack:   object.NewCallStack(),
		Env:     object.NewEnvironment(),
	}
}

func (e *Evaluator) Interrupt() {
	e.interrupt = true
}

func (e *Evaluator) ResetInterrupt() {
	e.interrupt = false
}

func (e *Evaluator) GetCallStack() object.CallStack {
	return e.Stack
}

func (e *Evaluator) GetEnvironment() *object.Environment {
	return e.Env
}

func (e *Evaluator) CloneClean() object.Evaluator {
	newE := *e
	newE.Env = object.NewEnvironment()
	return &newE
}

func (e *Evaluator) Eval(node ast.Node) object.Object {
	if e.interrupt {
		e.ResetInterrupt()
		return e.NewError(node, "Tiến trình bị ngắt")
	}

	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return e.evalProgram(node)

	case *ast.BlockStatement:
		return e.evalBlockStatement(node)

	case *ast.ExpressionStatement:
		return e.Eval(node.Expression)

	case *ast.ReturnStatement:
		val := e.Eval(node.ReturnValue)
		if object.IsError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.BreakStatement:
		return &object.BreakSignal{}

	case *ast.ContinueStatement:
		return &object.ContinueSignal{}

	case *ast.LetStatement:
		val := e.Eval(node.Value)
		if object.IsError(val) {
			return val
		}
		e.Env.Set(node.Name.Value, val)
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
		return e.evalPrefixExpression(node)

	case *ast.InfixExpression:
		return e.evalInfixExpression(node)

	case *ast.IfExpression:
		return e.evalIfExpression(node)

	case *ast.WhileExpression:
		return e.evalWhileExpression(node)

	case *ast.Identifier:
		return e.evalIdentifier(node)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: e.Env, Body: body}

	case *ast.CallExpression:
		return e.applyFunction(node)

	case *ast.ArrayLiteral:
		elements := e.evalExpressions(node.Elements)
		if len(elements) == 1 && object.IsError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.IndexExpression:
		return e.evalIndexExpression(node)

	case *ast.HashLiteral:
		return e.evalHashLiteral(node)
	default:
		return e.NewError(node, "phiên dịch thất bại, lệnh lạ!!")
	}

	return NULL
}

func (e *Evaluator) evalProgram(program *ast.Program) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = e.Eval(statement)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		case *object.BreakSignal:
			return e.NewError(statement, "Không thể 'ngắt' ngoài vòng lặp")
		case *object.ContinueSignal:
			return e.NewError(statement, "Không thể 'tiếp' ngoài vòng lặp")
		}

	}

	return result
}

func (e *Evaluator) evalBlockStatement(
	block *ast.BlockStatement,
) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = e.Eval(statement)

		switch result := result.(type) {
		case *object.Error:
			return result
		case *object.ReturnValue:
			return result
		case *object.BreakSignal:
			return result
		case *object.ContinueSignal:
			return result
		}
	}

	return result
}

func (e *Evaluator) evalIfExpression(
	ie *ast.IfExpression,
) object.Object {
	for i, cond := range ie.Condition {
		condition := e.Eval(cond)
		if object.IsError(condition) {
			return condition
		}

		if isTruthy(condition) {
			return e.Eval(ie.Consequence[i])
		}
	}

	if ie.Alternative != nil {
		return e.Eval(ie.Alternative)
	} else {
		return NULL
	}
}

func (e *Evaluator) evalWhileExpression(
	ie *ast.WhileExpression,
) object.Object {

	var result object.Object

	for {
		condition := e.Eval(ie.Condition)
		if object.IsError(condition) {
			return condition
		}

		if !isTruthy(condition) {
			break
		}

		result = e.Eval(ie.Body)

		switch result.(type) {
		case *object.ReturnValue:
			return result
		case *object.Error:
			return result
		case *object.BreakSignal:
			return result
		}
	}

	return unwrapWhileSignal(result)
}

func (e *Evaluator) evalIdentifier(
	node *ast.Identifier,
) object.Object {
	if val, ok := e.Env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := e.Builtin[node.Value]; ok {
		return builtin
	}

	return e.NewError(node, "không tìm thấy tên định danh: "+node.Value)
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

func (e *Evaluator) NewError(node ast.Node, format string, a ...interface{}) *object.Error {
	return &object.Error{Stack: e.Stack, Pos: node.Position(), Message: fmt.Sprintf(format, a...)}
}

func (e *Evaluator) evalExpressions(
	exps []ast.Expression,
) []object.Object {
	var result []object.Object

	for _, exp := range exps {
		evaluated := e.Eval(exp)
		if object.IsError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func (e *Evaluator) applyFunction(node *ast.CallExpression) object.Object {
	fn := e.Eval(node.Function)
	if object.IsError(fn) {
		return fn
	}

	args := e.evalExpressions(node.Arguments)
	if len(args) == 1 && object.IsError(args[0]) {
		return args[0]
	}

	s := e.Stack
	newS := append(e.Stack, object.ActivationRecord{
		CallNode: node,
		Function: fn,
		Args:     args,
	})

	switch fn := fn.(type) {

	case *object.Function:
		oldEnv := e.Env
		e.Env = extendFunctionEnv(fn, args)
		e.Stack = newS
		evaluated := e.Eval(fn.Body)
		e.Stack = s
		e.Env = oldEnv

		switch evaluated.(type) {
		case *object.BreakSignal:
			return e.NewError(node, "không thể 'ngắt' ngoài vòng lặp")
		case *object.ContinueSignal:
			return e.NewError(node, "không thể 'tiếp' ngoài vòng lặp")
		}

		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return fn.Fn(e, node, args...)
	default:
		return e.NewError(node, "không phải là một hàm: %s", fn.Type())
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
		return NULL
	case *object.ContinueSignal:
		return NULL
	default:
		return obj
	}
}

func (e *Evaluator) evalIndexExpression(node *ast.IndexExpression) object.Object {
	left := object.UnwrapReference(e.Eval(node.Left))
	if object.IsError(left) {
		return left
	}
	index := object.UnwrapReference(e.Eval(node.Index))
	if object.IsError(index) {
		return index
	}

	switch left := left.(type) {
	case *object.Array:
		index, ok := index.(*object.Integer)
		if !ok {
			return e.NewError(node, "chỉ có thể lấy chỉ mục 'nguyên' của 'mảng', nhận kiểu chỉ mục '%s'", index.Type())
		}
		return evalArrayIndexExpression(left, index)
	case *object.Hash:
		return e.evalHashIndexExpression(node, left, index)
	default:
		return e.NewError(node, "toán tử chỉ mục không hỗ trợ cho: %s", left.Type())
	}
}

func evalArrayIndexExpression(array *object.Array, index *object.Integer) object.Object {
	idxBig := index.Value
	if !idxBig.IsInt64() {
		return NULL
	}

	idx := idxBig.Int64()
	max := int64(len(array.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	if array.Mut == object.MUTABLE {
		return &object.RefObject{Obj: &(array.Elements[idx])}
	} else {
		return array.Elements[idx]
	}
}

func (e *Evaluator) evalHashLiteral(
	node *ast.HashLiteral,
) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := e.Eval(keyNode)
		if object.IsError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return e.NewError(node, "không thể dùng như khóa băm: %s", key.Type())
		}

		value := e.Eval(valueNode)
		if object.IsError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func (e *Evaluator) evalHashIndexExpression(node *ast.IndexExpression, hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return e.NewError(node, "không thể dùng như khóa băm: %s", index.Type())
	}

	hashKey := key.HashKey()
	pair, ok := hashObject.Pairs[hashKey]
	if hashObject.Mut == object.MUTABLE {
		if !ok {
			pair = object.HashPair{Key: index, Value: NULL}
			hashObject.Pairs[hashKey] = pair
		}
		return &object.RefObject{Obj: &pair.Value}
	} else {
		if !ok {
			return NULL
		}
		return pair.Value
	}
}
