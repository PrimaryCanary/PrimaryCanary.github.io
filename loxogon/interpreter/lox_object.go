package interpreter

import "fmt"

type LoxObject struct {
	Value any
}

func (lo LoxObject) String() string {
	if lo.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", lo.Value)
}
