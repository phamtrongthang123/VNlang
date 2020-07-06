package evaluator

import (
	"math/big"
	"vnlang/ast"
	"vnlang/object"
)

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func (e *Evaluator) changeMut(node ast.Node, obj object.Object, mut object.Mutability) object.Object {
	if obj.Mutable() == mut {
		return obj
	}

	switch obj := obj.(type) {
	case *object.Array:
		return obj.Clone(mut)
	case *object.Hash:
		return obj.Clone(mut)
	default:
		return e.NewError(node, "Không thể đổi kiểu cho '%s'", obj.Type())
	}
}

func (e *Evaluator) evalPrefixExpression(node *ast.PrefixExpression) object.Object {
	right := object.UnwrapReference(e.Eval(node.Right))
	if object.IsError(right) {
		return right
	}

	switch node.Operator {
	case "hằng":
		return e.changeMut(node, right, object.IMMUTABLE)
	case "biến":
		return e.changeMut(node, right, object.MUTABLE)
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return e.evalMinusPrefixOperatorExpression(node, right)
	default:
		return e.NewError(node, "toán tử lạ: %s%s", node.Operator, right.Type())
	}
}

func (e *Evaluator) evalInOperator(node *ast.InfixExpression, hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return e.NewError(node, "không thể dùng như khóa băm: %s", index.Type())
	}

	hashKey := key.HashKey()
	_, ok = hashObject.Pairs[hashKey]
	return &object.Boolean{Value: ok}
}

func (e *Evaluator) evalInfixExpression(
	node *ast.InfixExpression,
) object.Object {
	left := e.Eval(node.Left)
	if object.IsError(left) {
		return left
	}

	right := object.UnwrapReference(e.Eval(node.Right))
	if object.IsError(right) {
		return right
	}

	refLeft, ref := left.(*object.RefObject)
	switch node.Operator {
	case "=":
		if !ref || refLeft.Mutable() != object.MUTABLE {
			return e.NewError(node, "Vế trái không phải giá trị có thể gán")
		}
		*refLeft.Obj = right
		return left
	}

	// Unwrap reference
	if ref {
		left = *refLeft.Obj
	}

	switch left := left.(type) {
	case *object.Integer:
		if right, ok := right.(*object.Integer); ok {
			return e.evalIntegerInfixExpression(node, left, right)
		}
	case *object.Float:
		if right, ok := right.(*object.Float); ok {
			return e.evalFloatInfixExpression(node, left, right)
		}
	case *object.String:
		if right, ok := right.(*object.String); ok {
			return e.evalStringInfixExpression(node, left, right)
		}
	case *object.Boolean:
		if right, ok := right.(*object.Boolean); ok {
			return e.evalBooleanInfixExpression(node, left, right)
		}
	case *object.Array:
		if right, ok := right.(*object.Array); ok {
			return e.evalArrayInfixExpression(node, left, right)
		}
	}

	switch node.Operator {
	case "thuộc":
		right, ok := right.(*object.Hash)
		if !ok {
			return e.NewError(node, "toán tử 'thuộc' cần hạng tử bên phải thuộc kiểu 'băm' (nhận kiểu %s)", right.Type())
		}
		return e.evalInOperator(node, right, left)
	}

	return e.NewError(node, "toán tử lạ: %s %s %s",
		left.Type(), node.Operator, right.Type())
}

func evalBangOperatorExpression(right object.Object) *object.Boolean {
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

func (e *Evaluator) evalMinusPrefixOperatorExpression(node *ast.PrefixExpression, right object.Object) object.Object {
	switch value := right.(type) {
	case *object.Integer:
		var newValue big.Int
		newValue.Neg(value.Value)
		return &object.Integer{Value: &newValue}
	case *object.Float:
		return &object.Float{Value: -value.Value}
	default:
		return e.NewError(node, "toán tử lạ: -%s", right.Type())
	}
}

func (e *Evaluator) evalIntegerInfixExpression(
	node *ast.InfixExpression,
	left, right *object.Integer,
) object.Object {
	leftVal := left.Value
	rightVal := right.Value

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
		return e.NewError(node, "toán tử lạ: %s %s %s",
			left.Type(), node.Operator, right.Type())
	}
	return &object.Integer{Value: &resVal}
}

func (e *Evaluator) evalFloatInfixExpression(
	node *ast.InfixExpression,
	left, right *object.Float,
) object.Object {
	leftVal := left.Value
	rightVal := right.Value

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
		return e.NewError(node, "toán tử lạ: %s %s %s",
			left.Type(), node.Operator, right.Type())
	}
}

func (e *Evaluator) evalStringInfixExpression(
	node *ast.InfixExpression,
	left, right *object.String,
) object.Object {
	leftVal := left.Value
	rightVal := right.Value

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
		return e.NewError(node, "toán tử lạ: %s %s %s",
			left.Type(), node.Operator, right.Type())
	}

}

func (e *Evaluator) evalArrayInfixExpression(
	node *ast.InfixExpression,
	left, right *object.Array,
) object.Object {
	leftVal := left.Elements
	rightVal := right.Elements

	switch node.Operator {
	case "+":
		length := len(leftVal) + len(rightVal)
		newElements := make([]object.Object, length, length)

		copy(newElements, leftVal)
		copy(newElements[len(leftVal):], rightVal)
		return &object.Array{Elements: newElements}
	default:
		return e.NewError(node, "toán tử lạ: %s %s %s",
			left.Type(), node.Operator, right.Type())
	}
}

func (e *Evaluator) evalBooleanInfixExpression(
	node *ast.InfixExpression,
	left, right *object.Boolean,
) object.Object {
	leftVal := left.Value
	rightVal := right.Value

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
		return e.NewError(node, "toán tử lạ: %s %s %s",
			left.Type(), node.Operator, right.Type())
	}

}
