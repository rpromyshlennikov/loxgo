package interpreter

import (
	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/scanner"
)

type valuesStorage = map[string]any

type Environment struct {
	values    valuesStorage
	enclosing *Environment
}

func NewEnvironment(environment *Environment) Environment {
	return Environment{
		values:    make(valuesStorage),
		enclosing: environment,
	}
}

func (e Environment) get(name scanner.Token) (any, error) {
	v, ok := e.values[name.Lexeme()]
	if ok {
		return v, nil
	}
	if e.enclosing != nil {
		return e.enclosing.get(name)
	}
	return nil, NewRuntimeError(name, "Undefined variable '"+name.Lexeme()+"'.")
}

func (e Environment) define(name string, value any) {
	e.values[name] = value
}

func (e Environment) assign(name scanner.Token, value any) error {
	_, ok := e.values[name.Lexeme()]
	if ok {
		e.values[name.Lexeme()] = value
		return nil
	}
	if e.enclosing != nil {
		return e.enclosing.assign(name, value)
	}
	return NewRuntimeError(
		name,
		"Undefined variable '"+name.Lexeme()+"'.",
	)
}
