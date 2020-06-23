package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"runtime/pprof"
	"vnlang/evaluator"
	"vnlang/object"
	"vnlang/repl"
)

func setupProfiling(file string) {
	f, err := os.Create(file)
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	fmt.Println("Start profiling")
}

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
		if len(os.Args) >= 3 {
			setupProfiling(os.Args[2])
			defer pprof.StopCPUProfile()
		}
		evaluated := evaluator.ImportFile(&object.Import{Env: object.NewEnvironment()}, &object.String{Value: os.Args[1]})
		if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
			fmt.Println(evaluated.Inspect())
		}
	}
}
