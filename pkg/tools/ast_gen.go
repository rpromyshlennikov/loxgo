package tools

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func DefineAst(outputDir string, base string, types []string) {
	names := strings.Split(base, ",")
	absPath := path.Join(outputDir, strings.ToLower(names[0])+".go")
	file, err := os.Create(absPath)
	if err != nil {
		panic(err)
	}

	strBuilder := strings.Builder{}
	_, err = strBuilder.WriteString("package ast\n")
	if err != nil {
		panic(err)
	}

	defineImports(&strBuilder, names, types)

	defineVisitor(&strBuilder, names, types)

	// The AST classes.
	for _, kind := range types {
		defineType(&strBuilder, names, kind)
	}

	n, err := file.WriteString(strBuilder.String())
	if n < strBuilder.Len() {
		panic("error writing")
	}
	if err != nil {
		panic(err)
	}
}

func defineImports(builder *strings.Builder, names []string, types []string) {
	_, baseReturnType := baseNames(names)
	usedTypes := append([]string{baseReturnType}, types...)
	addImport := false
	for _, usedType := range usedTypes {
		if strings.Contains(usedType, "scanner.") {
			addImport = true
		}

	}
	if !addImport {
		return
	}
	_, err := builder.WriteString(`
import (
	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/scanner"
)
`)
	if err != nil {
		panic(err)
	}
}

func defineVisitor(builder *strings.Builder, names []string, types []string) {
	baseName, baseReturnType := baseNames(names)
	// Producing AST Node and Visitor interfaces.
	_, err := builder.WriteString(
		fmt.Sprintf(`
type %s interface {
	Accept(visitor Visitor%s)%s
}

type Visitor%s interface {
`,
			baseName, baseName, baseReturnType, baseName,
		),
	)
	if err != nil {
		panic(err)
	}

	// Producing types.
	for _, kind := range types {
		className, _ := classProps(kind)
		_, err = builder.WriteString(
			fmt.Sprintf(
				"\tVisit%s(*%s)%s\n",
				className, className, baseReturnType,
			),
		)
		if err != nil {
			panic(err)
		}
	}

	_, err = builder.WriteString("}\n")
	if err != nil {
		panic(err)
	}
}

// gocyclo considers this function too difficult, but it's just tool.
// TODO: refactor to several private methods or use templates after book will be completed.
//
//gocyclo:ignore
func defineType(builder *strings.Builder, names []string, kind string) {
	className, fieldsList := classProps(kind)
	fields := strings.Split(fieldsList, ", ")
	baseName, baseReturnType := baseNames(names)

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
	_, err = builder.WriteString("func New" + className + "(")
	if err != nil {
		panic(err)
	}
	parameters := make([]string, 0, len(fields))
	for _, field := range fields {
		parameters = append(parameters, strings.ToLower(field[:1])+field[1:])
		if err != nil {
			panic(err)
		}
	}
	_, err = builder.WriteString(strings.Join(parameters, ", "))
	if err != nil {
		panic(err)
	}
	_, err = builder.WriteString(
		") *" + className + " {\n" +
			"\tthis := " + className + "{}\n",
	)
	if err != nil {
		panic(err)
	}
	for _, field := range fields {
		name := strings.Split(field, " ")[0]
		_, err = builder.WriteString(
			"\tthis." + name + " = " + strings.ToLower(name[:1]) + name[1:] + "\n",
		)
		if err != nil {
			panic(err)
		}
	}
	_, err = builder.WriteString("\treturn &this\n}\n\n")
	if err != nil {
		panic(err)
	}

	returnStr := ""
	if baseReturnType != "" {
		returnStr = "return "
	}
	// Producing accept method.
	_, err = builder.WriteString(
		fmt.Sprintf(
			"func (%s *%s) Accept(visitor Visitor%s)%s {\n"+
				"\t%svisitor.Visit%s(%s)\n"+
				"}\n",
			strings.ToLower(className[:1]), className, baseName, baseReturnType, returnStr, className, strings.ToLower(className[:1]),
		),
	)
	if err != nil {
		panic(err)
	}
}
func baseNames(names []string) (string, string) {
	baseName := names[0]
	baseReturnType := ""
	if len(names) > 1 {
		baseReturnType = " " + names[1]
	}
	return baseName, baseReturnType
}

func classProps(kind string) (string, string) {
	className := strings.Trim(strings.Split(kind, ":")[0], " ")
	fieldsList := strings.Trim(strings.Split(kind, ":")[1], " ")
	return className, fieldsList
}
