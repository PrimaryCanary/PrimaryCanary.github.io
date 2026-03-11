package interpreter

import (
	"fmt"
	"time"
)

type LoxObject struct {
	Value any
}

type Callable interface {
	Call(*Interpreter, []LoxObject) (LoxObject, error)
	Arity() int
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
