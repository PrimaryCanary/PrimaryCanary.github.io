package parser

import (
	"fmt"
	"loxogon/token"
)

type ExprKind int
type StmtKind int

const (
	BINARY ExprKind = iota
	UNARY
	LITERAL
	GROUPING
)
const (
	EXPR StmtKind = iota
	PRINT
)

// Fat struct representation of expressions
type Expr struct {
	Kind     ExprKind
	Tok      token.Token
	Data     any
	Children []Expr
}

// Fat struct representation of expressions
type Stmt struct {
	Kind  StmtKind
	Child Expr
}

type parser struct {
	toks    []token.Token
	current int
}

type ParseError struct {
	tok     token.Token
	message string
}

// Parses a set of tokens into an AST.
//
// Lox ENBF grammar:
//
// program        → statement* EOF ;
// statement      → exprStmt | printStmt ;
// exprStmt       → expression ";" ;
// printStmt      → "print" expression ";" ;
// expression     → equality ;
// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
// term           → factor ( ( "-" | "+" ) factor )* ;
// factor         → unary ( ( "/" | "*" ) unary )* ;
// unary          → ( "!" | "-" ) unary | primary ;
// primary        → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ;
func Parse(toks []token.Token) ([]Stmt, error) {
	p := parser{toks: toks, current: 0}
	statements := make([]Stmt, 0)
	for !p.atEnd() {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}
	return statements, nil
}

// statement → exprStmt | printStmt ;
func (p *parser) statement() (Stmt, error) {
	if p.match(token.PRINT) {
		return p.printStatement()
	}
	return p.expressionStatement()
}

// exprStmt → expression ";" ;
func (p *parser) expressionStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return Stmt{}, err
	}
	_, err = p.consume(token.SEMICOLON, "expected ';' after expression")
	if err != nil {
		return Stmt{}, err
	}
	return newExprStmt(expr), nil
}

// printStmt → "print" expression ";" ;
func (p *parser) printStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return Stmt{}, err
	}
	_, err = p.consume(token.SEMICOLON, "expected ';' after expression")
	if err != nil {
		return Stmt{}, err
	}
	return newPrintStmt(expr), nil
}

// expression → equality
func (p *parser) expression() (Expr, error) {
	return p.equality()
}

func (p *parser) leftAssocBinaryExpr(operand func() (Expr, error), kinds ...token.TokenKind) (Expr, error) {
	left, err := operand()
	if err != nil {
		return Expr{}, err
	}
	for p.match(kinds...) {
		operator := p.previous()
		right, err := operand()
		if err != nil {
			return Expr{}, err
		}
		left = newBinary(operator, left, right)
	}
	return left, nil
}

// equality → comparison ( ( "!=" | "==" ) comparison )*
func (p *parser) equality() (Expr, error) {
	return p.leftAssocBinaryExpr(p.comparison, token.BANG_EQUAL, token.EQUAL_EQUAL)
}

// comparison → term ( ( ">" | ">=" | "<" | "<=" ) term )*
func (p *parser) comparison() (Expr, error) {
	return p.leftAssocBinaryExpr(p.term, token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL)
}

// term → factor ( ( "-" | "+" ) factor )*
func (p *parser) term() (Expr, error) {
	return p.leftAssocBinaryExpr(p.factor, token.MINUS, token.PLUS)
}

// factor → unary ( ( "/" | "*" ) unary )*
func (p *parser) factor() (Expr, error) {
	return p.leftAssocBinaryExpr(p.unary, token.STAR, token.SLASH)
}

// unary → ( "!" | "-" ) unary | primary
func (p *parser) unary() (Expr, error) {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return Expr{}, err
		}
		return newUnary(operator, right), nil
	}
	return p.primary()
}

// primary → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")"
func (p *parser) primary() (Expr, error) {
	if p.match(token.NUMBER, token.STRING) {
		data := p.previous().Literal
		return newLiteral(data), nil
	}
	if p.match(token.TRUE) {
		return newLiteral(true), nil
	}
	if p.match(token.FALSE) {
		return newLiteral(false), nil
	}
	if p.match(token.NIL) {
		return newLiteral(nil), nil
	}
	if p.match(token.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return Expr{}, err
		}
		_, err = p.consume(token.RIGHT_PAREN, "expected ')' after expression")
		if err != nil {
			return Expr{}, err
		}
		return newGrouping(expr), nil
	}

	return Expr{}, ParseError{p.peek(), "expected expression"}
}

// Check if the current token is the same as the argument.
func (p *parser) check(kind token.TokenKind) bool {
	return !p.atEnd() && (kind == p.peek().Kind)
}

// Check if the current token is the same as one of the arguments and advance if necessary.
func (p *parser) match(kinds ...token.TokenKind) bool {
	for _, k := range kinds {
		if p.check(k) {
			p.advance()
			return true
		}
	}
	return false
}

// Check if the current token is the same as the argument and advance if necessary.
// A failed match is an error.
func (p *parser) consume(kind token.TokenKind, message string) (token.Token, error) {
	if p.check(kind) {
		return p.advance(), nil
	}
	return token.Token{}, ParseError{p.peek(), message}
}

func (p *parser) peek() token.Token {
	return p.toks[p.current]
}

func (p *parser) previous() token.Token {
	return p.toks[p.current-1]
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

func newLiteral(data any) Expr {
	return Expr{Kind: LITERAL, Data: data}
}

func newBinary(operator token.Token, left, right Expr) Expr {
	return Expr{Kind: BINARY, Tok: operator, Children: []Expr{left, right}}
}

func newUnary(operator token.Token, e Expr) Expr {
	return Expr{Kind: UNARY, Tok: operator, Children: []Expr{e}}
}

func newGrouping(e Expr) Expr {
	return Expr{Kind: GROUPING, Children: []Expr{e}}
}

func newExprStmt(e Expr) Stmt {
	return Stmt{Kind: EXPR, Child: e}
}

func newPrintStmt(e Expr) Stmt {
	return Stmt{Kind: PRINT, Child: e}
}
func (e Expr) String() string {
	switch e.Kind {
	case BINARY:
		return fmt.Sprintf("(%s %s %v)", e.Tok.Lexeme, e.Children[0], e.Children[1])
	case UNARY:
		return fmt.Sprintf("(%s %v)", e.Tok.Lexeme, e.Children[0])
	case LITERAL:
		return fmt.Sprint(e.Data)
	case GROUPING:
		return fmt.Sprintf("(group %v)", e.Children[0])
	}
	return ""
}

func (pe ParseError) Error() string {
	if pe.tok.Kind == token.EOF {
		return fmt.Sprintf("[line %v] Parse error at end of input: %v",
			pe.tok.Line, pe.message)
	}
	return fmt.Sprintf("[line %v] Parse error at '%v': %v",
		pe.tok.Line, pe.tok.Lexeme, pe.message)
}
