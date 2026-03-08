package interpreter

import (
	"fmt"
	"loxogon/environment"
	"loxogon/parser"
	"loxogon/token"
)

type Interpreter struct {
	env environment.Env
}

type LoxObject struct {
	Value any
}

func New() Interpreter {
	return Interpreter{environment.New()}
}

func (i *Interpreter) Evaluate(expr parser.Expr) (LoxObject, error) {
	// TODO handle (0/0) - (0/0). correct semantics should not be IEEE-754 compliant
	switch expr.Kind {
	case parser.LITERAL:
		return LoxObject{expr.Data}, nil
	case parser.VARIABLE:
		value, err := i.env.Get(expr.Tok)
		if err != nil {
			return LoxObject{}, err
		}
		return LoxObject{Value: value}, nil
	case parser.GROUPING:
		value, err := i.Evaluate(expr.Children[0])
		if err != nil {
			return LoxObject{}, err
		}

		return value, nil
	case parser.UNARY:
		value, err := i.Evaluate(expr.Children[0])
		if err != nil {
			return LoxObject{}, err
		}

		switch expr.Tok.Kind {
		case token.MINUS:
			number, err := operandToNumber(expr.Tok, value)
			if err != nil {
				return LoxObject{}, err
			}
			return LoxObject{-number}, nil

		case token.BANG:
			return LoxObject{!isTruthy(value)}, nil
		}
	case parser.BINARY:
		left, err := i.Evaluate(expr.Children[0])
		if err != nil {
			return LoxObject{}, err
		}
		right, err := i.Evaluate(expr.Children[1])
		if err != nil {
			return LoxObject{}, err
		}
		switch expr.Tok.Kind {
		case token.PLUS:
			if l, r, err := operandsToNumbers(expr.Tok, left, right); err == nil {
				return LoxObject{l + r}, nil
			} else if l, r, err := operandsToStrings(expr.Tok, left, right); err == nil {
				return LoxObject{l + r}, nil
			} else {
				err := fmt.Errorf("[line %v] Runtime error: operator '+' requires number or string operands, got %v and %v",
					expr.Tok, left, right)
				return LoxObject{}, err
			}

		case token.MINUS:
			l, r, err := operandsToNumbers(expr.Tok, left, right)
			if err != nil {
				return LoxObject{}, err
			}
			return LoxObject{l - r}, nil
		case token.STAR:
			l, r, err := operandsToNumbers(expr.Tok, left, right)
			if err != nil {
				return LoxObject{}, err
			}
			return LoxObject{l * r}, nil
		case token.SLASH:
			// TODO divide by zero
			l, r, err := operandsToNumbers(expr.Tok, left, right)
			if err != nil {
				return LoxObject{}, err
			}
			return LoxObject{l / r}, nil
		case token.GREATER:
			l, r, err := operandsToNumbers(expr.Tok, left, right)
			if err != nil {
				return LoxObject{}, err
			}
			return LoxObject{l > r}, nil
		case token.GREATER_EQUAL:
			l, r, err := operandsToNumbers(expr.Tok, left, right)
			if err != nil {
				return LoxObject{}, err
			}
			return LoxObject{l >= r}, nil
		case token.LESS:
			l, r, err := operandsToNumbers(expr.Tok, left, right)
			if err != nil {
				return LoxObject{}, err
			}
			return LoxObject{l < r}, nil
		case token.LESS_EQUAL:
			l, r, err := operandsToNumbers(expr.Tok, left, right)
			if err != nil {
				return LoxObject{}, err
			}
			return LoxObject{l <= r}, nil
		case token.EQUAL_EQUAL:
			return LoxObject{isEqual(left, right)}, nil
		case token.BANG_EQUAL:
			return LoxObject{!isEqual(left, right)}, nil
		}
	}

	// Unreachable
	panic("Hit unreachable state in expression evaluation")
}

func (i *Interpreter) EvaluateStmt(stmt parser.Stmt) (LoxObject, error) {
	switch stmt.Kind {
	case parser.EXPR:
		result, err := i.Evaluate(stmt.Child)
		if err != nil {
			return LoxObject{}, err
		}
		return result, nil
	case parser.PRINT:
		_, err := i.Evaluate(stmt.Child)
		if err != nil {
			return LoxObject{}, err
		}
		return LoxObject{}, nil
	case parser.VAR_UNINIT:
		i.env.Define(stmt.Name.Lexeme, nil)
		return LoxObject{}, nil
	case parser.VAR:
		value, err := i.Evaluate(stmt.Child)
		if err != nil {
			return LoxObject{}, err
		}
		i.env.Define(stmt.Name.Lexeme, value.Value)
		return LoxObject{}, nil
	}

	// Unreachable
	panic("Hit unreachable state in statement evaluation")
}

// Lox semantics: nil == nil is true, nil == (any value) is false,
func isEqual(left, right any) bool {
	// TODO I'm not sure if this is a bug with nil interfaces
	if left == nil && right == nil {
		return true
	}
	if left == nil || right == nil {
		return false
	}
	return left == right
}

// Nil and false are the only falsey values
func isTruthy(value any) bool {
	// TODO I'm not sure if this is a bug
	if value == nil {
		return false
	}
	if boolean, ok := value.(bool); ok {
		return boolean
	}
	return true
}

func operandToNumber(operator token.Token, operand LoxObject) (float64, error) {
	if v, ok := operand.Value.(float64); ok {
		return v, nil
	}
	return 0, fmt.Errorf("[line %v] Runtime error: operator '%v' requires numeric operand, got %v",
		operator.Line, operator.Lexeme, operand)
}

func operandsToNumbers(operator token.Token, left, right LoxObject) (float64, float64, error) {
	l, leftOk := left.Value.(float64)
	r, rightOk := right.Value.(float64)
	if leftOk && rightOk {
		return l, r, nil
	}
	err := fmt.Errorf("[line %v] Runtime error: operator '%v' requires numeric operands, got %v and %v",
		operator.Line, operator.Lexeme, left, right)
	return 0, 0, err
}

func operandsToStrings(operator token.Token, left, right LoxObject) (string, string, error) {
	l, leftOk := left.Value.(string)
	r, rightOk := right.Value.(string)
	if leftOk && rightOk {
		return l, r, nil
	}
	err := fmt.Errorf("[line %v] Runtime error: operator '%v' requires string operands, got %v and %v",
		operator.Line, operator.Lexeme, left, right)
	return "", "", err
}

func (lo LoxObject) String() string {
	if lo.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", lo.Value)
}
