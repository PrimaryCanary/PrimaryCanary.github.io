package environment

import (
	"fmt"
	"loxogon/ast"
)

type Env struct {
	Bindings map[string]ast.LoxObject
}

type UndefVarError struct {
	token ast.Token
}

func New() Env {
	return Env{Bindings: make(map[string]ast.LoxObject)}
}

func (e *Env) Define(name string, value ast.LoxObject) {
	e.Bindings[name] = value
}

func (e *Env) Get(name ast.Token) (ast.LoxObject, error) {
	if value, ok := e.Bindings[name.Lexeme]; ok {
		return value, nil
	}
	return ast.LoxObject{}, UndefVarError{token: name}
}

func (uve UndefVarError) Error() string {
	return fmt.Sprintf("[line %v] Undefined variable '%v'", uve.token.Line, uve.token.Lexeme)
}
