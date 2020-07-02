package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"math/big"
	"strings"
	"vnlang/ast"
	"vnlang/scanner"
)

type BuiltinFunction func(node ast.Node, args ...Object) Object

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
	IMPORT_OBJ   = "SỬ_DỤNG"

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
}

type Integer struct {
	Value *big.Int
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return i.Value.Text(10) }
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

func (i *Float) Type() ObjectType { return FLOAT_OBJ }
func (i *Float) Inspect() string  { return fmt.Sprintf("%f", i.Value) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
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

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "rỗng" }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type BreakSignal struct {
}

func (rv *BreakSignal) Type() ObjectType { return BREAK_SIGNAL_OBJ }
func (rv *BreakSignal) Inspect() string  { return "ngắt" }

type ContinueSignal struct {
}

func (rv *ContinueSignal) Type() ObjectType { return CONTINUE_SIGNAL_OBJ }
func (rv *ContinueSignal) Inspect() string  { return "ngắt" }

type Error struct {
	Pos     scanner.Position
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "LỖI " + e.Pos.String() + ": " + e.Message }

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

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "hàm có sẵn" }

type Import struct {
	Env *Environment
}

func (b *Import) Type() ObjectType { return IMPORT_OBJ }
func (b *Import) Inspect() string  { return "sử dụng" }

type Array struct {
	Elements []Object
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

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
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
