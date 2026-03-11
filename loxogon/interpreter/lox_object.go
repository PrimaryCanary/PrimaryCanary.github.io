package interpreter

import (
	"fmt"
	"loxogon/ast"
	"time"
)

type LoxObject struct {
	Value any
}

type Callable interface {
	Call(*Interpreter, []LoxObject) (LoxObject, error)
	Arity() int
}

type LoxFunction struct {
	Decl ast.Stmt
}

func (lf LoxFunction) Call(i *Interpreter, args []LoxObject) (LoxObject, error) {
	funcEnv := NewEnv(i.globals)
	for i := range lf.Decl.Tokens[1:] {
		funcEnv.Define(lf.Decl.Tokens[i+1].Lexeme, args[i])
	}
	prev := i.env
	i.env = funcEnv
	err := i.EvaluateStmt(lf.Decl.Stmts[0])
	if err != nil {
		return LoxObject{}, err
	}
	i.env = prev
	return LoxObject{}, nil
}

func (lf LoxFunction) Arity() int {
	return len(lf.Decl.Tokens) - 1
}

type nativeClock struct{}

func (n nativeClock) Call(i *Interpreter, args []LoxObject) (LoxObject, error) {
	return LoxObject{float64(time.Now().Unix())}, nil
}

func (n nativeClock) Arity() int {
	return 0
}

func (lo LoxObject) String() string {
	if lo.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", lo.Value)
}
