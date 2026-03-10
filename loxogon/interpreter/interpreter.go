package interpreter

import (
	"fmt"
	"io"
	"loxogon/ast"
	"loxogon/environment"
	"os"
)

type Interpreter struct {
	env    environment.Env
	output io.Writer
}

func New() Interpreter {
	return Interpreter{environment.New(), os.Stdout}
}

func NewWithWriter(w io.Writer) Interpreter {
	return Interpreter{environment.New(), w}
}

func (i *Interpreter) Evaluate(expr ast.Expr) (ast.LoxObject, error) {
	// TODO handle (0/0) - (0/0). correct semantics should not be IEEE-754 compliant
	switch expr.Kind {
	case ast.LITERAL:
		return ast.LoxObject{Value: expr.Data}, nil
	case ast.VARIABLE:
		value, err := i.env.Get(expr.Tok)
		if err != nil {
			return ast.LoxObject{}, err
		}
		return value, nil
	case ast.ASSIGN:
		value, err := i.Evaluate(expr.Children[0])
		if err != nil {
			return ast.LoxObject{}, err
		}
		i.env.Assign(expr.Tok, value)
		return value, nil
	case ast.GROUPING:
		value, err := i.Evaluate(expr.Children[0])
		if err != nil {
			return ast.LoxObject{}, err
		}

		return value, nil
	case ast.UNARY:
		value, err := i.Evaluate(expr.Children[0])
		if err != nil {
			return ast.LoxObject{}, err
		}

		switch expr.Tok.Kind {
		case ast.MINUS:
			number, err := operandToNumber(expr.Tok, value)
			if err != nil {
				return ast.LoxObject{}, err
			}
			return ast.LoxObject{Value: -number}, nil

		case ast.BANG:
			return ast.LoxObject{Value: !isTruthy(value)}, nil
		}
	case ast.BINARY:
		left, err := i.Evaluate(expr.Children[0])
		if err != nil {
			return ast.LoxObject{}, err
		}
		right, err := i.Evaluate(expr.Children[1])
		if err != nil {
			return ast.LoxObject{}, err
		}
		switch expr.Tok.Kind {
		case ast.PLUS:
			if l, r, err := operandsToNumbers(expr.Tok, left, right); err == nil {
				return ast.LoxObject{Value: l + r}, nil
			} else if l, r, err := operandsToStrings(expr.Tok, left, right); err == nil {
				return ast.LoxObject{Value: l + r}, nil
			} else {
				str := fmt.Sprintf("operator '+' requires number or string operands, got '%v' and '%v'",
					left, right)
				return ast.LoxObject{}, RuntimeError{expr.Tok, str}
			}

		case ast.MINUS:
			l, r, err := operandsToNumbers(expr.Tok, left, right)
			if err != nil {
				return ast.LoxObject{}, err
			}
			return ast.LoxObject{Value: l - r}, nil
		case ast.STAR:
			l, r, err := operandsToNumbers(expr.Tok, left, right)
			if err != nil {
				return ast.LoxObject{}, err
			}
			return ast.LoxObject{Value: l * r}, nil
		case ast.SLASH:
			l, r, err := operandsToNumbers(expr.Tok, left, right)
			if err != nil {
				return ast.LoxObject{}, err
			}
			if r == 0.0 {
				return ast.LoxObject{}, RuntimeError{expr.Tok, "divided by zero"}
			}
			return ast.LoxObject{Value: l / r}, nil
		case ast.GREATER:
			l, r, err := operandsToNumbers(expr.Tok, left, right)
			if err != nil {
				return ast.LoxObject{}, err
			}
			return ast.LoxObject{Value: l > r}, nil
		case ast.GREATER_EQUAL:
			l, r, err := operandsToNumbers(expr.Tok, left, right)
			if err != nil {
				return ast.LoxObject{}, err
			}
			return ast.LoxObject{Value: l >= r}, nil
		case ast.LESS:
			l, r, err := operandsToNumbers(expr.Tok, left, right)
			if err != nil {
				return ast.LoxObject{}, err
			}
			return ast.LoxObject{Value: l < r}, nil
		case ast.LESS_EQUAL:
			l, r, err := operandsToNumbers(expr.Tok, left, right)
			if err != nil {
				return ast.LoxObject{}, err
			}
			return ast.LoxObject{Value: l <= r}, nil
		case ast.EQUAL_EQUAL:
			return ast.LoxObject{Value: isEqual(left, right)}, nil
		case ast.BANG_EQUAL:
			return ast.LoxObject{Value: !isEqual(left, right)}, nil
		}
	}

	// Unreachable
	panic("Hit unreachable state in expression evaluation")
}

func (i *Interpreter) EvaluateStmt(stmt ast.Stmt) (ast.LoxObject, error) {
	switch stmt.Kind {
	case ast.EXPR:
		result, err := i.Evaluate(stmt.Child)
		if err != nil {
			return ast.LoxObject{}, err
		}
		return result, nil
	case ast.PRINT:
		result, err := i.Evaluate(stmt.Child)
		if err != nil {
			return ast.LoxObject{}, err
		}
		_, err = i.output.Write([]byte(result.String() + "\n"))
		if err != nil {
			return ast.LoxObject{}, fmt.Errorf("failed writing to io.Writer: %w", err)
		}
		return ast.LoxObject{}, nil
	case ast.VAR_UNINIT:
		// TODO runtime error on uninitialized values
		i.env.Define(stmt.Name.Lexeme, ast.LoxObject{})
		return ast.LoxObject{}, nil
	case ast.VAR:
		value, err := i.Evaluate(stmt.Child)
		if err != nil {
			return ast.LoxObject{}, err
		}
		i.env.Define(stmt.Name.Lexeme, value)
		return ast.LoxObject{}, nil
	case ast.BLOCK:
		outerScope := i.env
		i.env = environment.NewWithParent(outerScope)
		for _, st := range stmt.Stmts {
			_, err := i.EvaluateStmt(st)
			if err != nil {
				return ast.LoxObject{}, err
			}
		}
		i.env = outerScope
		return ast.LoxObject{}, nil
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

func operandToNumber(operator ast.Token, operand ast.LoxObject) (float64, error) {
	if v, ok := operand.Value.(float64); ok {
		return v, nil
	}

	str := fmt.Sprintf("operator '%v' requires numeric operand, got %v",
		operator.Lexeme, operand)
	return 0, RuntimeError{operator, str}
}

func operandsToNumbers(operator ast.Token, left, right ast.LoxObject) (float64, float64, error) {
	l, leftOk := left.Value.(float64)
	r, rightOk := right.Value.(float64)
	if leftOk && rightOk {
		return l, r, nil
	}

	str := fmt.Sprintf("operator '%v' requires numeric operands, got '%v' and '%v'",
		operator.Lexeme, left, right)
	return 0, 0, RuntimeError{operator, str}
}

func operandsToStrings(operator ast.Token, left, right ast.LoxObject) (string, string, error) {
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
	return fmt.Sprintf("[line %v] Runtime error: %v", re.tok.Line, re.message)
}
