package repl

import (
	"fmt"
	"io"
	"vnlang/evaluator"
	"vnlang/lexer"
	"vnlang/object"
	"vnlang/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	env := object.NewEnvironment()

	l := lexer.New(in)
	p := parser.New(l)

	for {
		fmt.Printf(PROMPT)

		program := p.ParseOneStatementProgram()

		if program == nil {
			break
		}

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			p.ClearErrors()
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Ôi! Có vấn đề xảy ra rồi!\n")
	io.WriteString(out, " lỗi phân giải:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
