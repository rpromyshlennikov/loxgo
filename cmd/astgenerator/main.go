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
			"Binary   : left Expr, operator Token, right Expr",
			"Grouping : expression Expr",
			"Literal  : value any",
			"Unary    : operator Token, right Expr",
		},
	)
}
