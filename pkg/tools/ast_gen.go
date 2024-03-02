package tools

import (
	"os"
	"path"
	"strings"
)

func DefineAst(outputDir string, baseName string, types []string) {
	absPath := path.Join(outputDir, strings.ToLower(baseName)+".go")
	file, err := os.Create(absPath)
	if err != nil {
		panic(err)
	}

	strBuilder := strings.Builder{}
	_, err = strBuilder.WriteString(`package parser

import (
	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/scanner"
)

type Token = scanner.Token
type TokenType = scanner.TokenType

type Expr interface {
	accept(visitor Visitor) any
}
`)
	if err != nil {
		panic(err)
	}

	defineVisitor(&strBuilder, types)

	// The AST classes.
	for _, kind := range types {
		defineType(&strBuilder, kind)
	}

	n, err := file.WriteString(strBuilder.String())
	if n < strBuilder.Len() {
		panic("error writing")
	}
	if err != nil {
		panic(err)
	}
}

func defineVisitor(builder *strings.Builder, types []string) {
	// Producing Visitor interface.
	_, err := builder.WriteString("\ntype Visitor interface {\n")
	if err != nil {
		panic(err)
	}

	// Producing types.
	for _, kind := range types {
		className, _ := classProps(kind)
		_, err = builder.WriteString(
			"\tvisit" + className + "(*" + className + ") any\n")
		if err != nil {
			panic(err)
		}
	}

	_, err = builder.WriteString("}\n")
	if err != nil {
		panic(err)
	}
}

func defineType(builder *strings.Builder, kind string) {
	className, fieldsList := classProps(kind)
	fields := strings.Split(fieldsList, ", ")

	// Producing type for Expression.
	_, err := builder.WriteString("\ntype " + className + " struct {\n")
	if err != nil {
		panic(err)
	}
	for _, field := range fields {
		name := strings.Split(field, " ")[0]
		_, err = builder.WriteString(
			"\t// " + name + " field.\n" +
				"\t" + field + "\n",
		)
		if err != nil {
			panic(err)
		}
	}
	_, err = builder.WriteString("}\n\n")
	if err != nil {
		panic(err)
	}

	// Producing constructor
	_, err = builder.WriteString(
		"func New" + className +
			"(" + fieldsList + ") *" +
			className + " {\n" +
			"\tthis := " + className + "{}\n")
	if err != nil {
		panic(err)
	}
	for _, field := range fields {
		name := strings.Split(field, " ")[0]
		_, err = builder.WriteString(
			"\tthis." + name + " = " + name + "\n")
		if err != nil {
			panic(err)
		}
	}
	_, err = builder.WriteString("\treturn &this\n}\n\n")
	if err != nil {
		panic(err)
	}

	// Producing accept method.
	_, err = builder.WriteString(
		"func (" + strings.ToLower(className[:1]) + " *" + className + ") accept(visitor Visitor) any {\n" +
			"\treturn visitor.visit" + className + "(" + strings.ToLower(className[:1]) + ")\n" +
			"}\n")
	if err != nil {
		panic(err)
	}
}

func classProps(kind string) (string, string) {
	className := strings.Trim(strings.Split(kind, ":")[0], " ")
	fieldsList := strings.Trim(strings.Split(kind, ":")[1], " ")
	return className, fieldsList
}
