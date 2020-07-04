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

func (e *Evaluator) evalPrefixExpression(node *ast.PrefixExpression) object.Object {
	right := e.Eval(node.Right)
	if object.IsError(right) {
		return right
	}

	switch node.Operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return e.evalMinusPrefixOperatorExpression(node, right)
	default:
		return e.NewError(node, "toán tử lạ: %s%s", node.Operator, right.Type())
	}
}

func (e *Evaluator) evalInfixExpression(
	node *ast.InfixExpression,
) object.Object {
	left := e.Eval(node.Left)
	if object.IsError(left) {
		return left
	}

	right := e.Eval(node.Right)
	if object.IsError(right) {
		return right
	}

	leftType := left.Type()
	rightType := right.Type()

	if leftType != rightType {
		return e.NewError(node, "kiểu không tương thích: %s %s %s", leftType, node.Operator, rightType)
	}

	switch left := left.(type) {
	case *object.Integer:
		return e.evalIntegerInfixExpression(node, left, right.(*object.Integer))
	case *object.Float:
		return e.evalFloatInfixExpression(node, left, right.(*object.Float))
	case *object.String:
		return e.evalStringInfixExpression(node, left, right.(*object.String))
	case *object.Boolean:
		return e.evalBooleanInfixExpression(node, left, right.(*object.Boolean))
	case *object.Array:
		return e.evalArrayInfixExpression(node, left, right.(*object.Array))
	}

	switch node.Operator {
	case "==":
		return nativeBoolToBooleanObject(left == right)
	case "!=":
		return nativeBoolToBooleanObject(left != right)
	}

	return e.NewError(node, "toán tử lạ: %s %s %s",
		leftType, node.Operator, rightType)
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
