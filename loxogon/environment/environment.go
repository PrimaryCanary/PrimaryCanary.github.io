package environment

import (
	"fmt"
	"loxogon/token"
)

type Env struct {
	Bindings map[string]any
}

type UndefVarError struct {
	token token.Token
}

func New() Env {
	return Env{Bindings: make(map[string]any)}
}

func (e *Env) Define(name string, value any) {
	e.Bindings[name] = value
}

func (e *Env) Get(name token.Token) (any, error) {
	if value, ok := e.Bindings[name.Lexeme]; ok {
		return value, nil
	}
	return nil, UndefVarError{token: name}
}

func (uve UndefVarError) Error() string {
	return fmt.Sprintf("[line %v] Undefined variable '%v'", uve.token.Line, uve.token.Lexeme)
}
