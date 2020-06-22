package main

import (
	"fmt"
	"os"
	"os/user"
	"vnlang/evaluator"
	"vnlang/object"
	"vnlang/repl"
)

func main() {
	if len(os.Args) == 1 {
		user, err := user.Current()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Chào người dùng %s!\n",
			user.Username)
		repl.Start(os.Stdin, os.Stdout)
	} else {
		evaluated := evaluator.ImportFile(&object.Import{Env: object.NewEnvironment()}, &object.String{Value: os.Args[1]})
		if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
			fmt.Println(evaluated.Inspect())
		}
	}
}
