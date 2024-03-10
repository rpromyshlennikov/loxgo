package main

import (
	"fmt"
	"os"

	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/tools"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: astgenerator <output directory>")
		os.Exit(64)
	}
	if len(os.Args) == 2 {
		generateAst(os.Args[1])
	}
}

func generateAst(outputDir string) {
	fmt.Println(outputDir)

	tools.DefineAst(
		outputDir,
		"Expr",
		[]string{
			"Binary   : Left Expr, Operator Token, Right Expr",
			"Grouping : Expression Expr",
			"Literal  : Value any",
			"Unary    : Operator Token, Right Expr",
		},
	)
}
