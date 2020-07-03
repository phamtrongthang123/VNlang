package builtin

import (
	"bytes"
	"io"
	"os"
	"vnlang/lexer"
	"vnlang/object"
	"vnlang/parser"
)

func RunFile(e object.Evaluator, path string) object.Object {
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

	return e.Eval(prog)
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "lỗi phân giải:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
