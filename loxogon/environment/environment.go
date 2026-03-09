package environment

import (
	"fmt"
	"loxogon/ast"
)

type Env struct {
	// Bindings in the local scope
	bindings map[string]ast.LoxObject
	// Environment of enclosing scope
	parent *Env
}

type UndefVarError struct {
	token ast.Token
}

func New() Env {
	return Env{bindings: make(map[string]ast.LoxObject), parent: nil}
}

func NewWithParent(par Env) Env {
	return Env{bindings: make(map[string]ast.LoxObject), parent: &par}
}

func (e *Env) Define(name string, value ast.LoxObject) {
	e.bindings[name] = value
}

func (e *Env) Get(name ast.Token) (ast.LoxObject, error) {
	if value, ok := e.bindings[name.Lexeme]; ok {
		return value, nil
	}
	if e.parent != nil {
		return e.parent.Get(name)
	}
	return ast.LoxObject{}, UndefVarError{token: name}
}

func (e *Env) Assign(name ast.Token, value ast.LoxObject) error {
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
