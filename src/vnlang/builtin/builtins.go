package builtin

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"syscall"
	"vnlang/ast"
	"vnlang/lexer"
	"vnlang/object"
	"vnlang/parser"
	"vnlang/repl"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

var Builtin = object.BuiltinFnMap{
	// len
	"độ_dài": {Fn: func(e object.Evaluator, node ast.Node, args ...object.Object) object.Object {
		if len(args) != 1 {
			return e.NewError(node, "Sai số lượng tham số truyền vào. nhận được = %d, mong muốn = 1",
				len(args))
		}

		switch arg := args[0].(type) {
		case *object.Array:
			return &object.Integer{Value: big.NewInt(int64(len(arg.Elements)))}
		case *object.String:
			return &object.Integer{Value: big.NewInt(int64(len(arg.Value)))}
		default:
			return e.NewError(node, "Tham số truyền vào `độ_dài` không được hỗ trợ lấy độ dài (chỉ có Mảng hoặc Chuỗi được hỗ trợ), kiểu tham số %s.",
				args[0].Type())
		}
	},
	},
	"kiểu": {Fn: func(e object.Evaluator, node ast.Node, args ...object.Object) object.Object {
		if len(args) != 1 {
			return e.NewError(node, "Sai số lượng tham số truyền vào. nhận được = %d, mong muốn = 1",
				len(args))
		}

		return args[0].Type()
	},
	},
	// convert big int to float
	"thực": {Fn: func(e object.Evaluator, node ast.Node, args ...object.Object) object.Object {
		if len(args) != 1 {
			return e.NewError(node, "Sai số lượng tham số truyền vào. nhận được = %d, mong muốn = 1",
				len(args))
		}

		switch val := args[0].(type) {
		case *object.Float:
			return val
		case *object.Integer:
			valFloat, _ := new(big.Float).SetInt(val.Value).Float64()
			return &object.Float{Value: valFloat}
		default:
			return e.NewError(node, "Tham số truyền vào `thực` phải là số nguyên hoặc số thực, kiểu tham số %s.",
				args[0].Type())
		}
	},
	},
	"nguyên": {Fn: func(e object.Evaluator, node ast.Node, args ...object.Object) object.Object {
		if len(args) != 1 {
			return e.NewError(node, "Sai số lượng tham số truyền vào. nhận được = %d, mong muốn = 1",
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
			return e.NewError(node, "Tham số truyền vào `nguyên` phải là số nguyên hoặc số thực, kiểu tham số %s.",
				args[0].Type())
		}
	},
	},
	// string() // convert object to string
	"xâu": {Fn: func(e object.Evaluator, node ast.Node, args ...object.Object) object.Object {
		if len(args) != 1 {
			return e.NewError(node, "Sai số lượng tham số truyền vào. nhận được = %d, mong muốn = 1",
				len(args))
		}

		return &object.String{Value: args[0].Inspect()}
	},
	},
	// puts
	"in_ra": {
		Fn: func(e object.Evaluator, node ast.Node, args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Print(arg.Inspect(), " ")
			}
			fmt.Println()

			return NULL
		},
	},
	//first
	"đầu": {
		Fn: func(e object.Evaluator, node ast.Node, args ...object.Object) object.Object {
			if len(args) != 1 {
				return e.NewError(node, "Sai số lượng tham số truyền vào. nhận được = %d, mong muốn = 1",
					len(args))
			}

			arr, ok := args[0].(*object.Array)
			if !ok {
				return e.NewError(node, "Tham số truyền vào hàm lấy `đầu` của mảng phải thuộc kiểu Mảng. Nhận được kiểu %s",
					args[0].Type())
			}

			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}

			return NULL
		},
	},
	// last
	"đuôi": {
		Fn: func(e object.Evaluator, node ast.Node, args ...object.Object) object.Object {
			if len(args) != 1 {
				return e.NewError(node, "Sai số lượng tham số truyền vào. nhận được = %d, mong muốn = 1",
					len(args))
			}

			arr, ok := args[0].(*object.Array)
			if !ok {
				return e.NewError(node, "Tham số truyền vào hàm lấy `đuôi` của mảng phải thuộc kiểu Mảng. Nhận được kiểu %s",
					args[0].Type())
			}

			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}

			return NULL
		},
	},
	// rest, rest returns a new array containing all elements of the array passed as argument, except the `first one`.
	"trừ_đầu": {
		Fn: func(e object.Evaluator, node ast.Node, args ...object.Object) object.Object {
			if len(args) != 1 {
				return e.NewError(node, "Sai số lượng tham số truyền vào. nhận được = %d, mong muốn = 1",
					len(args))
			}

			arr, ok := args[0].(*object.Array)

			if !ok {
				return e.NewError(node, "Tham số truyền vào hàm lấy `trừ_đầu` của mảng phải thuộc kiểu Mảng. Nhận được kiểu %s",
					args[0].Type())
			}

			if len(arr.Elements) == 0 {
				return NULL
			}
			return &object.Array{Elements: arr.Elements[1:]}
		},
	},
	// push
	"đẩy": {
		Fn: func(e object.Evaluator, node ast.Node, args ...object.Object) object.Object {
			if len(args) != 2 {
				return e.NewError(node, "Sai số lượng tham số truyền vào. nhận được = %d, mong muốn = 2",
					len(args))
			}

			arr, ok := args[0].(*object.Array)
			if !ok {
				return e.NewError(node, "Tham số truyền vào hàm lấy `đẩy` của mảng phải thuộc kiểu Mảng. Nhận được kiểu %s",
					args[0].Type())
			}

			if arr.Mut == object.IMMUTABLE {
				length := len(arr.Elements)
				newElements := make([]object.Object, length+1, length+1)

				copy(newElements, arr.Elements)
				newElements[length] = args[1]

				return &object.Array{Elements: newElements}
			} else if arr.Mut == object.MUTABLE {
				return &object.Array{Elements: append(arr.Elements, args[1]), Mut: object.MUTABLE}
			} else {
				return e.NewError(node, "Không thể xảy ra ??!!")
			}
		},
	},
	"sử_dụng": {
		Fn: func(e object.Evaluator, node ast.Node, args ...object.Object) object.Object {
			if len(args) != 1 {
				return e.NewError(node, "Sai số lượng tham số truyền vào. nhận được = %d, mong muốn = 1",
					len(args))
			}

			path, ok := args[0].(*object.String)

			if !ok {
				return e.NewError(node, "Tham số truyền vào hàm `sử_dụng` phải là một xâu đường dẫn. Nhận được kiểu %s",
					args[0].Type())
			}

			newE := e.CloneClean()
			evaluated := RunFile(newE, path.Value)
			if !object.IsError(evaluated) {
				data, ok := newE.GetEnvironment().Get("xuất")
				if ok {
					return data
				} else {
					return NULL
				}
			}

			return evaluated
		},
	},
	"thoát": {
		Fn: func(e object.Evaluator, node ast.Node, args ...object.Object) object.Object {
			if len(args) > 1 {
				return e.NewError(node, "Sai số lượng tham số truyền vào. nhận được = %d, mong muốn = 0 hoặc 1",
					len(args))
			}

			exitCode := 0
			if len(args) > 0 {
				arg, ok := args[0].(*object.Integer)
				if !ok {
					return e.NewError(node, "Tham số là một số nguyên (exit code). Nhận được %s",
						args[0].Type())
				}
				if arg.Value.IsInt64() {
					exitCode = int(arg.Value.Int64())
				}
			}
			os.Exit(exitCode)
			return NULL
		},
	},
	"thăm_dò": {
		Fn: func(e object.Evaluator, node ast.Node, args ...object.Object) object.Object {
			if len(args) > 0 {
				return e.NewError(node, "Sai số lượng tham số truyền vào. nhận được = %d, mong muốn = 0",
					len(args))
			}
			fmt.Printf("Bắt đầu thăm dò ở: %v\n", node.Position())
			repl.Start(e, os.Stdin, os.Stdout)
			fmt.Printf("Kết thúc thăm dò\n")

			return NULL
		},
	},
	// load(dll filename: str) -> dict{syscall}
	"thư_viện": {
		Fn: func(e object.Evaluator, node ast.Node, args ...object.Object) object.Object {
			if len(args) != 1 {
				return e.NewError(node, "Sai số lượng tham số truyền vào. nhận được = %d, mong muốn = 1",
					len(args))
			}

			path, ok := args[0].(*object.String)

			if !ok {
				return e.NewError(node, "Tham số truyền vào hàm `thư_viện` phải là một xâu đường dẫn. Nhận được kiểu %s",
					args[0].Type())
			}
			h, err := syscall.LoadLibrary(path.Value) //Make sure this DLL follows Golang machine bit architecture (64-bit in my case)
			if err != nil {
				log.Fatal(err)
			}
			// syscall
			define_syscall_func := "hàm(func_name, arg1, arg2) { trả_về gọi_hàm_ngoài_toàn_chuỗi_2( " + fmt.Sprint(h) + ",func_name, arg1, arg2);}"
			l := lexer.New(strings.NewReader(define_syscall_func), "<test>")
			p := parser.New(l)
			program := p.ParseProgram()
			evaluated := e.Eval(program)
			syscall_func, ok := evaluated.(*object.Function)
			// base_pairs := map[object.HashKey]object.Function{
			// 	(&object.String{Value: "gọi_hàm_ngoài_toàn_chuỗi_2"}).HashKey(): *syscall_func,
			// }
			pairs := make(map[object.HashKey]object.HashPair)
			key := &object.String{Value: "gọi_hàm_ngoài_toàn_chuỗi_2"}
			hashed := key.HashKey()
			pairs[hashed] = object.HashPair{Key: key, Value: syscall_func}
			return &object.Hash{Pairs: pairs, Mut: object.IMMUTABLE}
		},
	},
	// Fix builtin syscall2 (dll: handle, function name: str, arg1: str, arg2: str) -> str
	"gọi_hàm_ngoài_toàn_chuỗi_2": {
		Fn: func(e object.Evaluator, node ast.Node, args ...object.Object) object.Object {
			if len(args) != 4 {
				return e.NewError(node, "Sai số lượng tham số truyền vào. nhận được = %d, mong muốn = 4",
					len(args))
			}
			handle, ok := args[0].(*object.Integer)
			ffi_func, ok := args[1].(*object.String)
			if !ok {
				return e.NewError(node, "Tham số truyền vào hàm `thư_viện` phải là một xâu đường dẫn. Nhận được kiểu %s",
					args[0].Type())
			}
			arg1, ok := args[2].(*object.String)
			arg2, ok := args[3].(*object.String)
			// Add more if later
			res := syscall2_str_helper(syscall.Handle(handle.Value.Int64()), ffi_func.Value, arg1.Value, arg2.Value)
			return &object.String{Value: res}
		},
	},
}
