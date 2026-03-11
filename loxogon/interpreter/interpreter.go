package interpreter

import (
	"fmt"
	"io"
	"loxogon/ast"
	"os"
)

type Interpreter struct {
	env      environment
	output   io.Writer
	LastExpr LoxObject
	globals  environment
}

func New() Interpreter {
	globals := NewEnv()
	globals.Define("clock", LoxObject{Value: nativeClock{}})
	i := Interpreter{globals, os.Stdout, LoxObject{}, globals}
	return i
}

func NewWithWriter(w io.Writer) Interpreter {
	i := New()
	i.output = w
	return i
}

func (i *Interpreter) Evaluate(expr ast.Expr) (LoxObject, error) {
	// TODO handle (0/0) - (0/0). correct semantics should not be IEEE-754 compliant
	switch expr.Kind {
	case ast.LITERAL:
		return LoxObject{Value: expr.Data}, nil
	case ast.VARIABLE:
		value, err := i.env.Get(expr.Tok)
		if err != nil {
			return LoxObject{}, err
		}
		return value, nil
	case ast.ASSIGN:
		value, err := i.Evaluate(expr.Children[0])
		if err != nil {
			return LoxObject{}, err
		}
		i.env.Assign(expr.Tok, value)
		return value, nil
	case ast.GROUPING:
		value, err := i.Evaluate(expr.Children[0])
		if err != nil {
			return LoxObject{}, err
		}

		return value, nil
	case ast.UNARY:
		value, err := i.Evaluate(expr.Children[0])
		if err != nil {
			return LoxObject{}, err
		}

		switch expr.Tok.Kind {
		case ast.MINUS:
			number, err := operandToNumber(expr.Tok, value)
			if err != nil {
				return LoxObject{}, err
			}
			return LoxObject{Value: -number}, nil

		case ast.BANG:
			return LoxObject{Value: !isTruthy(value)}, nil
		}
	case ast.LOGICAL:
		left, err := i.Evaluate(expr.Children[0])
		if err != nil {
			return LoxObject{}, err
		}
		switch expr.Tok.Kind {
		case ast.OR:
			if isTruthy(left) {
				return left, nil
			}
		case ast.AND:
			if !isTruthy(left) {
				return left, nil
			}
		}
		right, err := i.Evaluate(expr.Children[1])
		if err != nil {
			return LoxObject{}, err
		}
		return right, nil
	case ast.BINARY:
		left, err := i.Evaluate(expr.Children[0])
		if err != nil {
			return LoxObject{}, err
		}
		right, err := i.Evaluate(expr.Children[1])
		if err != nil {
			return LoxObject{}, err
		}
		switch expr.Tok.Kind {
		case ast.PLUS:
			if l, r, err := operandsToNumbers(expr.Tok, left, right); err == nil {
				return LoxObject{Value: l + r}, nil
			} else if l, r, err := operandsToStrings(expr.Tok, left, right); err == nil {
				return LoxObject{Value: l + r}, nil
			} else {
				str := fmt.Sprintf("operator '+' requires number or string operands, got '%v' and '%v'",
					left, right)
				return LoxObject{}, RuntimeError{expr.Tok, str}
			}

		case ast.MINUS:
			l, r, err := operandsToNumbers(expr.Tok, left, right)
			if err != nil {
				return LoxObject{}, err
			}
			return LoxObject{Value: l - r}, nil
		case ast.STAR:
			l, r, err := operandsToNumbers(expr.Tok, left, right)
			if err != nil {
				return LoxObject{}, err
			}
			return LoxObject{Value: l * r}, nil
		case ast.SLASH:
			l, r, err := operandsToNumbers(expr.Tok, left, right)
			if err != nil {
				return LoxObject{}, err
			}
			if r == 0.0 {
				return LoxObject{}, RuntimeError{expr.Tok, "divided by zero"}
			}
			return LoxObject{Value: l / r}, nil
		case ast.GREATER:
			l, r, err := operandsToNumbers(expr.Tok, left, right)
			if err != nil {
				return LoxObject{}, err
			}
			return LoxObject{Value: l > r}, nil
		case ast.GREATER_EQUAL:
			l, r, err := operandsToNumbers(expr.Tok, left, right)
			if err != nil {
				return LoxObject{}, err
			}
			return LoxObject{Value: l >= r}, nil
		case ast.LESS:
			l, r, err := operandsToNumbers(expr.Tok, left, right)
			if err != nil {
				return LoxObject{}, err
			}
			return LoxObject{Value: l < r}, nil
		case ast.LESS_EQUAL:
			l, r, err := operandsToNumbers(expr.Tok, left, right)
			if err != nil {
				return LoxObject{}, err
			}
			return LoxObject{Value: l <= r}, nil
		case ast.EQUAL_EQUAL:
			return LoxObject{Value: isEqual(left, right)}, nil
		case ast.BANG_EQUAL:
			return LoxObject{Value: !isEqual(left, right)}, nil
		}
	case ast.CALL:
		callee, err := i.Evaluate(expr.Children[0])
		if err != nil {
			return LoxObject{}, err
		}
		args := make([]LoxObject, 0, len(expr.Children)-1)
		for _, arg := range expr.Children[1:] {
			a, err := i.Evaluate(arg)
			if err != nil {
				return LoxObject{}, err
			}
			args = append(args, a)
		}
		fn, ok := callee.Value.(Callable)
		if !ok {
			return LoxObject{}, RuntimeError{expr.Tok, "can only call functions and classes"}
		}
		if arity, n := fn.Arity(), len(args); arity != n {
			str := fmt.Sprintf("expected %v arguments but got %v", arity, n)
			return LoxObject{}, RuntimeError{expr.Tok, str}
		}
		return fn.Call(i, args)
	}

	// Unreachable
	panic("Hit unreachable state in expression evaluation")
}

func (i *Interpreter) EvaluateStmt(stmt ast.Stmt) error {
	i.LastExpr = LoxObject{}
	switch stmt.Kind {
	case ast.EXPR:
		result, err := i.Evaluate(stmt.Child)
		if err != nil {
			return err
		}
		i.LastExpr = result
		return nil
	case ast.PRINT:
		result, err := i.Evaluate(stmt.Child)
		if err != nil {
			return err
		}
		_, err = i.output.Write([]byte(result.String() + "\n"))
		if err != nil {
			return fmt.Errorf("failed writing to io.Writer: %w", err)
		}
		return nil
	case ast.VAR_UNINIT:
		// TODO runtime error on uninitialized values
		i.env.Define(stmt.Tokens[0].Lexeme, LoxObject{})
		return nil
	case ast.VAR:
		value, err := i.Evaluate(stmt.Child)
		if err != nil {
			return err
		}
		i.env.Define(stmt.Tokens[0].Lexeme, value)
		return nil
	case ast.BLOCK:
		outerScope := i.env
		i.env = NewEnv(outerScope)
		for _, st := range stmt.Stmts {
			if err := i.EvaluateStmt(st); err != nil {
				return err
			}
		}
		i.env = outerScope
		return nil
	case ast.IF:
		cond, err := i.Evaluate(stmt.Child)
		if err != nil {
			return err
		}
		if isTruthy(cond) {
			if err := i.EvaluateStmt(stmt.Stmts[0]); err != nil {
				return err
			}
			return nil
		} else if len(stmt.Stmts) > 1 {
			if err := i.EvaluateStmt(stmt.Stmts[1]); err != nil {
				return err
			}
			return nil
		}
		return nil
	case ast.WHILE:
		for {
			cond, err := i.Evaluate(stmt.Child)
			if err != nil {
				return err
			}
			if isTruthy(cond) {
				err = i.EvaluateStmt(stmt.Stmts[0])
				if err != nil {
					return err
				}
			} else {
				return nil
			}
		}
	case ast.FUN:
		fn := LoxFunction{Decl: stmt}
		i.env.Define(stmt.Tokens[0].Lexeme, LoxObject{fn})
		return nil
	case ast.RETURN_EMPTY:
		// Hijack error handling to return control flow to the caller
		// It's hacky and I hate it buuuuuuut....
		return ReturnError{returnValue: nil}
	case ast.RETURN:
		value, err := i.Evaluate(stmt.Child)
		if err != nil {
			return err
		}
		// Hijack error handling to return control flow to the caller
		// It's hacky and I hate it buuuuuuut....
		return ReturnError{returnValue: &value}
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
func isTruthy(obj LoxObject) bool {
	// TODO I'm not sure if this is a bug
	if obj.Value == nil {
		return false
	}
	if boolean, ok := obj.Value.(bool); ok {
		return boolean
	}
	return true
}

func operandToNumber(operator ast.Token, operand LoxObject) (float64, error) {
	if v, ok := operand.Value.(float64); ok {
		return v, nil
	}

	str := fmt.Sprintf("operator '%v' requires numeric operand, got %v",
		operator.Lexeme, operand)
	return 0, RuntimeError{operator, str}
}

func operandsToNumbers(operator ast.Token, left, right LoxObject) (float64, float64, error) {
	l, leftOk := left.Value.(float64)
	r, rightOk := right.Value.(float64)
	if leftOk && rightOk {
		return l, r, nil
	}

	str := fmt.Sprintf("operator '%v' requires numeric operands, got '%v' and '%v'",
		operator.Lexeme, left, right)
	return 0, 0, RuntimeError{operator, str}
}

func operandsToStrings(operator ast.Token, left, right LoxObject) (string, string, error) {
	l, leftOk := left.Value.(string)
	r, rightOk := right.Value.(string)
	if leftOk && rightOk {
		return l, r, nil
	}

	str := fmt.Sprintf("operator '%v' requires string operands, got '%v' and '%v'",
		operator.Lexeme, left, right)
	return "", "", RuntimeError{operator, str}
}

type RuntimeError struct {
	tok     ast.Token
	message string
}

func (re RuntimeError) Error() string {
	return fmt.Sprintf("[line %v at '%v'] Runtime error: %v", re.tok.Line, re.tok.Lexeme, re.message)
}

type ReturnError struct {
	returnValue *LoxObject
}

func (r ReturnError) Error() string {
	return ""
}
