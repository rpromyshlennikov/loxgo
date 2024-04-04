package interpreter

import (
	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/scanner"
)

type valuesStorage = map[string]any

type Environment struct {
	values valuesStorage
}

func NewEnvironment() Environment {
	return Environment{
		values: make(valuesStorage),
	}
}

func (e Environment) get(name scanner.Token) (any, error) {
	v, ok := e.values[name.Lexeme()]
	if !ok {
		return nil, NewRuntimeError(name, "Undefined variable '"+name.Lexeme()+"'.")
	}
	return v, nil
}

func (e Environment) define(name string, value any) {
	e.values[name] = value
}

func (e Environment) assign(name scanner.Token, value any) error {
	_, ok := e.values[name.Lexeme()]
	if !ok {
		return NewRuntimeError(
			name,
			"Undefined variable '"+name.Lexeme()+"'.",
		)
	}
	e.values[name.Lexeme()] = value
	return nil
}
