package parser

import (
	"fmt"
	"loxogon/token"
)

type ExprKind int

const (
	BINARY ExprKind = iota
	UNARY
	LITERAL
	GROUPING
)

// Fat struct representation of expressions
type Expr struct {
	Kind     ExprKind
	Operator token.Token
	Data     any
	Children []Expr
}

func (e Expr) String() string {
	switch e.Kind {
	case BINARY:
		return fmt.Sprintf("(%s %s %v)", e.Operator.Lexeme, e.Children[0], e.Children[1])
	case UNARY:
		return fmt.Sprintf("(%s %v)", e.Operator.Lexeme, e.Children[0])
	case LITERAL:
		return fmt.Sprint(e.Data)
	case GROUPING:
		return fmt.Sprintf("(group %v)", e.Children[0])
	}
	return ""
}

type parser struct {
	tokens  []token.Token
	current int
}

func Parse(tokens []token.Token) (Expr, error) {
	p := parser{tokens: tokens, current: 0}
	return p.expression(), nil
}

func (p *parser) expression() Expr {
	return p.equality()
}

func (p *parser) leftAssocBinaryExpr(operand func() Expr, kinds ...token.TokenKind) Expr {
	left := operand()
	for p.match(kinds...) {
		operator := p.previous()
		right := operand()
		left = NewBinary(operator, left, right)
	}
	return left
}

func (p *parser) equality() Expr {
	return p.leftAssocBinaryExpr(p.comparison, token.BANG_EQUAL, token.EQUAL_EQUAL)
}

func (p *parser) comparison() Expr {
	return p.leftAssocBinaryExpr(p.term, token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL)
}

func (p *parser) term() Expr {
	return p.leftAssocBinaryExpr(p.factor, token.MINUS, token.PLUS)
}

func (p *parser) factor() Expr {
	return p.leftAssocBinaryExpr(p.unary, token.STAR, token.SLASH)
}

func (p *parser) unary() Expr {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right := p.unary()
		return NewUnary(operator, right)
	}
	return p.primary()
}

func (p *parser) primary() Expr {
	if p.match(token.NUMBER, token.STRING) {
		data := p.previous().Literal
		return NewLiteral(data)
	}
	if p.match(token.TRUE) {
		return NewLiteral(true)
	}
	if p.match(token.FALSE) {
		return NewLiteral(false)
	}
	if p.match(token.NIL) {
		return NewLiteral(nil)
	}
	if p.match(token.LEFT_PAREN) {
		expr := p.expression()
		// TODO error handling
		p.match(token.RIGHT_PAREN)
		return NewGrouping(expr)
	}

	panic("TODO error handling in the parser")
	// return Expr{} // Unreachable, TODO panic
}

func (p *parser) match(kinds ...token.TokenKind) bool {
	for _, k := range kinds {
		if (k == p.peek().Kind) && !p.atEnd() {
			p.advance()
			return true
		}
	}
	return false
}

func (p *parser) peek() token.Token {
	return p.tokens[p.current]
}

func (p *parser) previous() token.Token {
	return p.tokens[p.current-1]
}

func (p *parser) advance() token.Token {
	if !p.atEnd() {
		p.current += 1
	}
	return p.previous()
}

func (p *parser) atEnd() bool {
	return p.peek().Kind == token.EOF
}

func NewLiteral(data any) Expr {
	return Expr{Kind: LITERAL, Data: data}
}

func NewBinary(operator token.Token, left, right Expr) Expr {
	return Expr{Kind: BINARY, Operator: operator, Children: []Expr{left, right}}
}

func NewUnary(operator token.Token, e Expr) Expr {
	return Expr{Kind: UNARY, Operator: operator, Children: []Expr{e}}
}

func NewGrouping(e Expr) Expr {
	return Expr{Kind: GROUPING, Children: []Expr{e}}
}
