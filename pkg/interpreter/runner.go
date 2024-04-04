package interpreter

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/parser"
	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/scanner"
)

type LoxGo struct {
	hadError        bool
	hadRuntimeError bool

	interpreter Interpreter
}

func New() *LoxGo {
	return &LoxGo{
		hadError:        false,
		hadRuntimeError: false,
		interpreter:     NewInterpreter(),
	}
}

func (lox *LoxGo) RunFile(fileName string) {
	sources, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("running file: %s\n", fileName)
	lox.Run(string(sources))
	if lox.hadError {
		os.Exit(65)
	}
	if lox.hadRuntimeError {
		os.Exit(70)
	}
}

func (lox *LoxGo) RunPrompt() {
	log.Println("running prompt")
	scannr := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		ok := scannr.Scan()
		if !ok {
			break
		}
		input := scannr.Text()
		if input == "" {
			continue
		}
		lox.Run(input)
		lox.hadError = false
		lox.hadRuntimeError = false
	}
	if err := scannr.Err(); err != nil {
		log.Println(err)
	}
}

func (lox *LoxGo) Run(sources string) {
	errRepCallback := func(line int, message string) {
		lox.erro(line, message)
	}
	scannr := scanner.NewScanner(sources, errRepCallback)
	tokens := scannr.ScanTokens()

	parsr := parser.NewParser(tokens, errRepCallback)
	statements := parsr.Parse()

	// For now, just print the AST.
	//fmt.Println((&parser.AstPrinter{}).Sprint(astTree))
	//fmt.Println(plugins.AstPrinter{}.Sprint(astTree))
	if lox.hadError {
		return
	}

	// Trying to interpret.
	err := lox.interpreter.Interpret(statements)
	if err != nil {
		lox.runtimeError(err)
	}
}

func (lox *LoxGo) erro(line int, message string) {
	lox.report(line, "", message)
}

func (lox *LoxGo) report(line int, where string, message string) {
	_, err := fmt.Fprintln(os.Stderr, "[line ", line, "] Error"+where+": "+message)
	if err != nil {
		log.Fatalln(err)
	}
	lox.hadError = true
}

func (lox *LoxGo) runtimeError(err error) {
	rErr, ok := err.(*RuntimeError)
	if !ok {
		log.Fatalln(err)
	}
	_, err = fmt.Fprintln(os.Stderr, "[line ", rErr.token.Line(), "] Runtime error: "+rErr.message)
	if err != nil {
		log.Fatalln(err)
	}
	lox.hadRuntimeError = true
}
