package evaluator

import (
	"bytes"
	"io"
	"os"
	"vnlang/lexer"
	"vnlang/object"
	"vnlang/parser"
)

const importKeyword = "sử_dụng"

func RunFile(path string, env *object.Environment) object.Object {
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return &object.Error{Message: "không thể mở file"}
	}
	pr := parser.New(lexer.New(file, file.Name()))
	prog := pr.ParseProgram()
	if len(pr.Errors()) != 0 {
		errStr := bytes.NewBufferString("")
		printParserErrors(errStr, pr.Errors())
		return &object.Error{Message: errStr.String()}
	}

	return Eval(prog, env)
}

func ImportFile(p *object.Import, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("Sai số lượng tham số truyền vào. nhận được = %d, mong muốn = 1",
			len(args))
	}
	if args[0].Type() != object.STRING_OBJ {
		return newError("Tham số truyền vào hàm `sử_dụng` phải là một xâu đường dẫn. Nhận được kiểu %s",
			args[0].Type())
	}

	path := args[0].(*object.String).Value
	newEnv := object.NewEnvironment()
	evaluated := RunFile(path, newEnv)
	if !isError(evaluated) {
		data, ok := newEnv.Get("xuất")
		if ok {
			return data
		} else {
			return NULL
		}
	}

	return evaluated
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "lỗi phân giải:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
