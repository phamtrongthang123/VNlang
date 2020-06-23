package evaluator

import (
	"fmt"
	"math/big"
	"os"
	"vnlang/object"
)

var builtins = map[string]*object.Builtin{
	// len
	"độ_dài": {Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("Sai số lượng tham số truyền vào. nhận được = %d, mong muốn = 1",
				len(args))
		}

		switch arg := args[0].(type) {
		case *object.Array:
			return &object.Integer{Value: big.NewInt(int64(len(arg.Elements)))}
		case *object.String:
			return &object.Integer{Value: big.NewInt(int64(len(arg.Value)))}
		default:
			return newError("Tham số truyền vào `độ_dài` không được hỗ trợ lấy độ dài (chỉ có Mảng hoặc Chuỗi được hỗ trợ), kiểu tham số %s.",
				args[0].Type())
		}
	},
	},
	// convert big int to float
	"thực": {Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("Sai số lượng tham số truyền vào. nhận được = %d, mong muốn = 1",
				len(args))
		}

		switch val := args[0].(type) {
		case *object.Float:
			return val
		case *object.Integer:
			valFloat, _ := new(big.Float).SetInt(val.Value).Float64()
			return &object.Float{Value: valFloat}
		default:
			return newError("Tham số truyền vào `thực` phải là số nguyên hoặc số thực, kiểu tham số %s.",
				args[0].Type())
		}
	},
	},
	"nguyên": {Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("Sai số lượng tham số truyền vào. nhận được = %d, mong muốn = 1",
				len(args))
		}

		switch val := args[0].(type) {
		case *object.Float:
			valInt := new(big.Int)
			big.NewFloat(val.Value).Int(valInt)
			return &object.Integer{Value: valInt}
		case *object.Integer:
			return val
		default:
			return newError("Tham số truyền vào `nguyên` phải là số nguyên hoặc số thực, kiểu tham số %s.",
				args[0].Type())
		}
	},
	},
	// string() // convert object to string
	"xâu": {Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("Sai số lượng tham số truyền vào. nhận được = %d, mong muốn = 1",
				len(args))
		}

		return &object.String{Value: args[0].Inspect()}
	},
	},
	// puts
	"in_ra": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Print(arg.Inspect(), " ")
			}
			fmt.Println()

			return NULL
		},
	},
	//first
	"đầu": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Sai số lượng tham số truyền vào. nhận được = %d, mong muốn = 1",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("Tham số truyền vào hàm lấy `đầu` của mảng phải thuộc kiểu Mảng. Nhận được kiểu %s",
					args[0].Type())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}

			return NULL
		},
	},
	// last
	"đuôi": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Sai số lượng tham số truyền vào. nhận được = %d, mong muốn = 1",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("Tham số truyền vào hàm lấy `đuôi` của mảng phải thuộc kiểu Mảng. Nhận được kiểu %s",
					args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}

			return NULL
		},
	},
	// rest, rest returns a new array containing all elements of the array passed as argument, except the `first one`.
	"trừ_đầu": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Sai số lượng tham số truyền vào. nhận được = %d, mong muốn = 1",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("Tham số truyền vào hàm lấy `trừ_đầu` của mảng phải thuộc kiểu Mảng. Nhận được kiểu %s",
					args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				newElements := make([]object.Object, length-1, length-1)
				copy(newElements, arr.Elements[1:length])
				return &object.Array{Elements: newElements}
			}

			return NULL
		},
	},
	// push
	"đẩy": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("Sai số lượng tham số truyền vào. nhận được = %d, mong muốn = 2",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("Tham số truyền vào hàm lấy `đẩy` của mảng phải thuộc kiểu Mảng. Nhận được kiểu %s",
					args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)

			newElements := make([]object.Object, length+1, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]

			return &object.Array{Elements: newElements}
		},
	},
	"thoát": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) > 1 {
				return newError("Sai số lượng tham số truyền vào. nhận được = %d, mong muốn = 0 hoặc 1",
					len(args))
			}
			if len(args) > 0 && args[0].Type() != object.INTEGER_OBJ {
				return newError("Tham số là một số nguyên (exit code). Nhận được %s",
					args[0].Type())
			}
			exitCode := 0
			if len(args) > 0 {
				exitCodeBig := args[0].(*object.Integer).Value
				if exitCodeBig.IsInt64() {
					exitCode = int(args[0].(*object.Integer).Value.Int64())
				}
			}
			os.Exit(exitCode)
			return NULL
		},
	},
}
