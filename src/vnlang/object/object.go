package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"io"
	"math/big"
	"strings"
	"vnlang/ast"
	"vnlang/scanner"
)

type Mutability int

const (
	IMMUTABLE Mutability = iota
	MUTABLE
)

type ActivationRecord struct {
	CallNode *ast.CallExpression
	Function Object
	Args     []Object
}

type CallStack []ActivationRecord

type Evaluator interface {
	Interrupt()
	ResetInterrupt()
	GetCallStack() CallStack
	GetEnvironment() *Environment
	Eval(node ast.Node) Object
	NewError(node ast.Node, format string, a ...interface{}) *Error
	CloneClean() Evaluator
}

func (s CallStack) PrintCallStack(out io.Writer, level int) {
	n := len(s) - level
	if n < 0 {
		n = 0
	}

	for i := len(s) - 1; i >= n; i-- {
		fmt.Fprintf(out, "%d: %v %v\n", i, s[i].CallNode.Position(), s[i].CallNode)
	}

	if len(s) > level {
		fmt.Fprintf(out, "...\n")
	}
}

type BuiltinFunction func(e Evaluator, node ast.Node, args ...Object) Object
type BuiltinFnMap map[string]*Builtin

type ObjectType string

const (
	NULL_OBJ  = "NULL"
	ERROR_OBJ = "LỖI"

	INTEGER_OBJ = "SỐ_NGUYÊN"
	FLOAT_OBJ   = "SỐ_THỰC"
	BOOLEAN_OBJ = "BOOLEAN"
	STRING_OBJ  = "XÂU"

	RETURN_VALUE_OBJ    = "GIÁ_TRỊ_TRẢ_VỀ"
	BREAK_SIGNAL_OBJ    = "TÍN_HIỆU_NGẮT"
	CONTINUE_SIGNAL_OBJ = "TÍN_HIỆU_TIẾP"

	FUNCTION_OBJ = "HÀM"
	BUILTIN_OBJ  = "CÓ_SẴN"

	ARRAY_OBJ = "MẢNG"
	HASH_OBJ  = "BĂM"
)

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type Hashable interface {
	HashKey() HashKey
}

type Object interface {
	Type() ObjectType
	Inspect() string
	Mutable() Mutability
}

func IsError(obj Object) bool {
	_, ok := obj.(*Error)
	return ok
}

type Integer struct {
	Value *big.Int
}

func (i *Integer) Type() ObjectType    { return INTEGER_OBJ }
func (i *Integer) Inspect() string     { return i.Value.Text(10) }
func (i *Integer) Mutable() Mutability { return IMMUTABLE }

func (i *Integer) HashKey() HashKey {
	h := fnv.New64a()
	data := i.Value.Bytes()
	h.Write(data)
	var sign byte = 0
	if i.Value.Sign() < 0 {
		sign = 43
	}
	h.Write([]byte{sign})
	return HashKey{Type: i.Type(), Value: h.Sum64()}
}

type Float struct {
	Value float64
}

func (i *Float) Type() ObjectType    { return FLOAT_OBJ }
func (i *Float) Inspect() string     { return fmt.Sprintf("%f", i.Value) }
func (i *Float) Mutable() Mutability { return IMMUTABLE }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType    { return BOOLEAN_OBJ }
func (b *Boolean) Mutable() Mutability { return IMMUTABLE }
func (b *Boolean) Inspect() string {
	if b.Value {
		return "đúng"
	} else {
		return "sai"
	}
}

func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

type Null struct{}

func (n *Null) Type() ObjectType    { return NULL_OBJ }
func (n *Null) Inspect() string     { return "rỗng" }
func (n *Null) Mutable() Mutability { return IMMUTABLE }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType    { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string     { return rv.Value.Inspect() }
func (rv *ReturnValue) Mutable() Mutability { return IMMUTABLE }

type BreakSignal struct {
}

func (bs *BreakSignal) Type() ObjectType    { return BREAK_SIGNAL_OBJ }
func (bs *BreakSignal) Inspect() string     { return "ngắt" }
func (bs *BreakSignal) Mutable() Mutability { return IMMUTABLE }

type ContinueSignal struct {
}

func (cs *ContinueSignal) Type() ObjectType    { return CONTINUE_SIGNAL_OBJ }
func (cs *ContinueSignal) Inspect() string     { return "ngắt" }
func (cs *ContinueSignal) Mutable() Mutability { return IMMUTABLE }

type Error struct {
	Stack   CallStack
	Pos     scanner.Position
	Message string
}

func (e *Error) Type() ObjectType    { return ERROR_OBJ }
func (e *Error) Inspect() string     { return "LỖI " + e.Pos.String() + ": " + e.Message }
func (e *Error) Mutable() Mutability { return IMMUTABLE }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("hàm ")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}
func (f *Function) Mutable() Mutability { return IMMUTABLE }

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}
func (s *String) Mutable() Mutability { return IMMUTABLE }

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType    { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string     { return "hàm có sẵn" }
func (b *Builtin) Mutable() Mutability { return IMMUTABLE }

type Array struct {
	Elements []Object
	Mut      Mutability
}

func (ao *Array) Type() ObjectType { return ARRAY_OBJ }
func (ao *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
func (ao *Array) Mutable() Mutability { return ao.Mut }
func (ao *Array) Clone(mut Mutability) *Array {
	length := len(ao.Elements)
	newElements := make([]Object, length, length)
	copy(newElements, ao.Elements)

	return &Array{Elements: newElements, Mut: mut}
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
	Mut   Mutability
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
func (h *Hash) Mutable() Mutability { return h.Mut }
func (h *Hash) Clone(mut Mutability) *Hash {
	// Create the target map
	newMap := make(map[HashKey]HashPair)

	// Copy from the original map to the target map
	for key, value := range h.Pairs {
		newMap[key] = value
	}

	return &Hash{Pairs: newMap, Mut: mut}
}

type RefObject struct {
	Obj *Object
}

func (ro *RefObject) Type() ObjectType { return "TRỎ(" + (*ro.Obj).Type() + ")" }
func (ro *RefObject) Inspect() string {
	return (*ro.Obj).Inspect()
}
func (ro *RefObject) Mutable() Mutability { return MUTABLE }

func UnwrapReference(obj Object) Object {
	switch obj := obj.(type) {
	case *RefObject:
		return *obj.Obj
	default:
		return obj
	}
}
