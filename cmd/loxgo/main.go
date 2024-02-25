package main

import (
	"fmt"
	"os"

	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/interpreter"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: loxgo [script]")
		os.Exit(64)
	}
	lox := interpreter.New()
	if len(os.Args) == 2 {
		lox.RunFile(os.Args[1])
		return
	}
	lox.RunPrompt()
	fmt.Println("Exiting...")
}
