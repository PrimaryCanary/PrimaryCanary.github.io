package interpreter

import (
	"fmt"
	"loxogon/ast"
)

type environment struct {
	// Bindings in the local scope
	bindings map[string]LoxObject
	// Environment of enclosing scope
	parent *environment
}

type UndefVarError struct {
	token ast.Token
}

func NewEnv(parent ...environment) environment {
	if len(parent) == 0 {
		return environment{bindings: make(map[string]LoxObject), parent: nil}
	} else {
		return environment{bindings: make(map[string]LoxObject), parent: &parent[0]}

	}
}

func (e *environment) Define(name string, value LoxObject) {
	e.bindings[name] = value
}

func (e *environment) Get(name ast.Token) (LoxObject, error) {
	if value, ok := e.bindings[name.Lexeme]; ok {
		return value, nil
	}
	if e.parent != nil {
		return e.parent.Get(name)
	}
	return LoxObject{}, UndefVarError{token: name}
}

func (e *environment) Assign(name ast.Token, value LoxObject) error {
	if _, ok := e.bindings[name.Lexeme]; ok {
		e.bindings[name.Lexeme] = value
		return nil
	}
	if e.parent != nil {
		return e.parent.Assign(name, value)
	}
	return UndefVarError{name}
}

func (uve UndefVarError) Error() string {
	return fmt.Sprintf("[line %v] Undefined variable '%v'", uve.token.Line, uve.token.Lexeme)
}
