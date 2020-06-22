package repl

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"vnlang/evaluator"
	"vnlang/lexer"
	"vnlang/object"
	"vnlang/parser"
)

const PROMPT = ">> "

func SetupInterrupt() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	go func() {
		for s := range signalChannel {
			switch s {
			case os.Interrupt:
				evaluator.Interrupt()
			}
		}
	}()
}

func Start(in io.Reader, out io.Writer) {
	SetupInterrupt()
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

		evaluator.ResetInterrupt()
		evaluated := evaluator.Eval(program, env)
		if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
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
