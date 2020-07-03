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
	fmt.Println("Bắt đầu đo đạc")
}

func runRepl() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Chào người dùng %s!\n",
		user.Username)
	repl.Start(os.Stdin, os.Stdout)
}

func runScript(args []string) {
	env := object.NewEnvironment()
	s := object.NewCallStack()
	env.Set("tham_số", toArgsArray(args))
	evaluated := evaluator.RunFile(s, args[0], env)
	if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
		fmt.Println(evaluated.Inspect())
	}
}

func main() {
	var i int = 1
	for i < len(os.Args) {
		if os.Args[i][0] != '-' {
			break
		}
		if os.Args[i] == "-đo_đạc" {
			i++
			if i >= len(os.Args) {
				panic("Cờ -đo_đạc <file_kết_quả> cần có tham số")
			}
			setupProfiling(os.Args[i])
			defer pprof.StopCPUProfile()
		}
		i++
	}

	if len(os.Args) <= i {
		runRepl()
	} else {
		runScript(os.Args[1:])
	}
}
