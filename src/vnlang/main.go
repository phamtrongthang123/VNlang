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

func toArgsArray(args []string) *object.Array {
	res := &object.Array{
		Elements: make([]object.Object, len(args)),
	}

	for i, arg := range args {
		res.Elements[i] = &object.String{Value: arg}
	}
	return res
}

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
		var i int = 1
		for {
			if os.Args[i][0] != '-' {
				break
			}
			if os.Args[i] == "-profile" {
				i++
				setupProfiling(os.Args[i])
				defer pprof.StopCPUProfile()
			}
			i++
		}

		env := object.NewEnvironment()
		env.Set("tham_số", toArgsArray(os.Args[i:]))
		evaluated := evaluator.RunFile(os.Args[i], env)
		if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
			fmt.Println(evaluated.Inspect())
		}
	}
}
