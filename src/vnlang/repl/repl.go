package repl

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"vnlang/lexer"
	"vnlang/object"
	"vnlang/parser"
)

const PROMPT = ">> "

func SetupInterrupt(e object.Evaluator) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	go func() {
		for s := range signalChannel {
			switch s {
			case os.Interrupt:
				e.Interrupt()
			}
		}
	}()
}

func ResetInterrupt() {
	signal.Reset(os.Interrupt)
}

func Start(e object.Evaluator, in io.Reader, out io.Writer) {
	SetupInterrupt(e)
	defer ResetInterrupt()

	l := lexer.New(in, "")
	p := parser.New(l)

	for {
		fmt.Printf(PROMPT)

		p.ResetLineCount()
		program := p.ParseOneStatementProgram()

		if program == nil {
			break
		}

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			p.ClearErrors()
			continue
		}

		e.ResetInterrupt()
		evaluated := e.Eval(program)
		if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
			errors, ok := evaluated.(*object.Error)
			if ok {
				errors.Stack.PrintCallStack(out, 10)
			}
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
