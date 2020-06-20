package main

import (
	"fmt"
	"os"
	"os/user"
	"vnlang/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Chào người dùng %s!\n",
		user.Username)
	repl.Start(os.Stdin, os.Stdout)
}
