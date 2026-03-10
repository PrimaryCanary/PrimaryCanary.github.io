package parser

import (
	"fmt"
	"loxogon/ast"
)

type parser struct {
	toks    []ast.Token
	current int
}

type ParseError struct {
	tok     ast.Token
	message string
}

// Parses a set of tokens into an AST.
//
// Lox ENBF grammar:
//
// program        → declaration* EOF ;
// declaration    → varDecl | statement ;
// varDecl        → "var" IDENTIFIER ( "=" expression )? ";" ;
// program        → statement* EOF ;
// statement      → exprStmt | printStmt | block | ifStmt;
// exprStmt       → expression ";" ;
// printStmt      → "print" expression ";" ;
// ifStmt         → "if" "(" expression ")" statement ( "else" statement )? ;
// block          → "{" declaration* "}" ;
// expression     → assignment ;
// assignment     → IDENTIFIER "=" assignment | equality ;
// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
// term           → factor ( ( "-" | "+" ) factor )* ;
// factor         → unary ( ( "/" | "*" ) unary )* ;
// unary          → ( "!" | "-" ) unary | primary ;
// primary        → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" | IDENTIFIER;
func Parse(toks []ast.Token) ([]ast.Stmt, error) {
	p := parser{toks: toks, current: 0}
	statements := make([]ast.Stmt, 0)
	for !p.atEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}
	return statements, nil
}

// declaration → varDecl | statement ;
func (p *parser) declaration() (ast.Stmt, error) {
	if p.match(ast.VAR_TOK) {
		stmt, err := p.varDecl()
		if err != nil {
			p.synchronize()
			return ast.Stmt{}, err
		}
		return stmt, nil
	}

	stmt, err := p.statement()
	if err != nil {
		p.synchronize()
		return ast.Stmt{}, err
	}
	return stmt, nil
}

// varDecl → "var" IDENTIFIER ( "=" expression )? ";" ;
func (p *parser) varDecl() (ast.Stmt, error) {
	ident, err := p.consume(ast.IDENTIFIER, "expected identifier after var")
	if err != nil {
		return ast.Stmt{}, err
	}

	var decl ast.Stmt
	if p.match(ast.EQUAL) {
		initializer, err := p.expression()
		if err != nil {
			return ast.Stmt{}, err
		}
		decl = ast.NewVarDecl(ident, initializer)
	} else {
		decl = ast.NewVarDeclUninit(ident)
	}

	_, err = p.consume(ast.SEMICOLON, "expected semicolon after variable declaration")
	if err != nil {
		return ast.Stmt{}, err
	}
	return decl, nil
}

// statement → exprStmt | printStmt | block | ifStmt;
func (p *parser) statement() (ast.Stmt, error) {
	if p.match(ast.PRINT_TOK) {
		return p.printStatement()
	}
	if p.match(ast.IF_TOK) {
		return p.ifStmt()
	}
	if p.match(ast.LEFT_BRACE) {
		stmts, err := p.block()
		if err != nil {
			return ast.Stmt{}, err
		}
		return ast.NewBlock(stmts), nil
	}

	return p.expressionStatement()
}

// ifStmt → "if" "(" expression ")" statement ( "else" statement )? ;
func (p *parser) ifStmt() (ast.Stmt, error) {
	_, err := p.consume(ast.LEFT_PAREN, "expected '(' after if keyword")
	if err != nil {
		return ast.Stmt{}, err
	}

	cond, err := p.expression()
	if err != nil {
		return ast.Stmt{}, err
	}
	_, err = p.consume(ast.RIGHT_PAREN, "expected ')' after if condition")
	if err != nil {
		return ast.Stmt{}, err
	}

	thenBranch, err := p.statement()
	if err != nil {
		return ast.Stmt{}, err
	}

	if p.match(ast.ELSE) {
		elseBranch, err := p.statement()
		if err != nil {
			return ast.Stmt{}, err
		}
		return ast.NewIf(cond, thenBranch, elseBranch), nil
	}

	return ast.NewIf(cond, thenBranch), nil
}

// block → "{" declaration* "}" ;
func (p *parser) block() ([]ast.Stmt, error) {
	stmts := make([]ast.Stmt, 0)
	for !p.check(ast.RIGHT_BRACE) && !p.atEnd() {
		st, err := p.declaration()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, st)
	}
	_, err := p.consume(ast.RIGHT_BRACE, "expected '}' after block")
	if err != nil {
		return nil, err
	}
	return stmts, nil
}

// exprStmt → expression ";" ;
func (p *parser) expressionStatement() (ast.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return ast.Stmt{}, err
	}
	_, err = p.consume(ast.SEMICOLON, "expected ';' after expression")
	if err != nil {
		return ast.Stmt{}, err
	}
	return ast.NewExprStmt(expr), nil
}

// printStmt → "print" expression ";" ;
func (p *parser) printStatement() (ast.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return ast.Stmt{}, err
	}
	_, err = p.consume(ast.SEMICOLON, "expected ';' after expression")
	if err != nil {
		return ast.Stmt{}, err
	}
	return ast.NewPrintStmt(expr), nil
}

// expression → assignment ;
func (p *parser) expression() (ast.Expr, error) {
	return p.assignment()
}

// assignment → IDENTIFIER "=" assignment | equality ;
func (p *parser) assignment() (ast.Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return ast.Expr{}, err
	}

	if p.match(ast.EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return ast.Expr{}, nil
		}

		if expr.Kind == ast.VARIABLE {
			name := expr.Tok
			return ast.NewAssign(name, value), nil
		}
		return ast.Expr{}, ParseError{equals, "invalid assignment target"}
	}

	return expr, nil
}

func (p *parser) leftAssocBinaryExpr(operand func() (ast.Expr, error), kinds ...ast.TokenKind) (ast.Expr, error) {
	left, err := operand()
	if err != nil {
		return ast.Expr{}, err
	}
	for p.match(kinds...) {
		operator := p.previous()
		right, err := operand()
		if err != nil {
			return ast.Expr{}, err
		}
		left = ast.NewBinary(operator, left, right)
	}
	return left, nil
}

// equality → comparison ( ( "!=" | "==" ) comparison )*
func (p *parser) equality() (ast.Expr, error) {
	return p.leftAssocBinaryExpr(p.comparison, ast.BANG_EQUAL, ast.EQUAL_EQUAL)
}

// comparison → term ( ( ">" | ">=" | "<" | "<=" ) term )*
func (p *parser) comparison() (ast.Expr, error) {
	return p.leftAssocBinaryExpr(p.term, ast.GREATER, ast.GREATER_EQUAL, ast.LESS, ast.LESS_EQUAL)
}

// term → factor ( ( "-" | "+" ) factor )*
func (p *parser) term() (ast.Expr, error) {
	return p.leftAssocBinaryExpr(p.factor, ast.MINUS, ast.PLUS)
}

// factor → unary ( ( "/" | "*" ) unary )*
func (p *parser) factor() (ast.Expr, error) {
	return p.leftAssocBinaryExpr(p.unary, ast.STAR, ast.SLASH)
}

// unary → ( "!" | "-" ) unary | primary
func (p *parser) unary() (ast.Expr, error) {
	if p.match(ast.BANG, ast.MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return ast.Expr{}, err
		}
		return ast.NewUnary(operator, right), nil
	}
	return p.primary()
}

// primary → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" | IDENTIFIER;
func (p *parser) primary() (ast.Expr, error) {
	if p.match(ast.NUMBER, ast.STRING) {
		data := p.previous().Literal
		return ast.NewLiteral(data), nil
	}
	if p.match(ast.TRUE) {
		return ast.NewLiteral(true), nil
	}
	if p.match(ast.FALSE) {
		return ast.NewLiteral(false), nil
	}
	if p.match(ast.NIL) {
		return ast.NewLiteral(nil), nil
	}
	if p.match(ast.IDENTIFIER) {
		return ast.NewVariable(p.previous()), nil
	}
	if p.match(ast.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return ast.Expr{}, err
		}
		_, err = p.consume(ast.RIGHT_PAREN, "expected ')' after expression")
		if err != nil {
			return ast.Expr{}, err
		}
		return ast.NewGrouping(expr), nil
	}

	return ast.Expr{}, ParseError{p.peek(), "expected expression"}
}

// Check if the current token is the same as the argument.
func (p *parser) check(kind ast.TokenKind) bool {
	return !p.atEnd() && (kind == p.peek().Kind)
}

// Check if the current token is the same as one of the arguments and advance if necessary.
func (p *parser) match(kinds ...ast.TokenKind) bool {
	for _, k := range kinds {
		if p.check(k) {
			p.advance()
			return true
		}
	}
	return false
}

// Discard tokens until parser is (probably) at a statement boundary
func (p *parser) synchronize() {
	// Skip unexpected token
	p.advance()

	for !p.atEnd() {
		if p.previous().Kind == ast.SEMICOLON {
			return
		}
		switch p.peek().Kind {
		case ast.CLASS, ast.FUN, ast.VAR_TOK, ast.FOR,
			ast.IF_TOK, ast.WHILE, ast.PRINT_TOK, ast.RETURN:
			return
		}
		p.advance()
	}
}

// Check if the current token is the same as the argument and advance if necessary.
// A failed match is an error.
func (p *parser) consume(kind ast.TokenKind, message string) (ast.Token, error) {
	if p.check(kind) {
		return p.advance(), nil
	}
	return ast.Token{}, ParseError{p.peek(), message}
}

func (p *parser) peek() ast.Token {
	return p.toks[p.current]
}

func (p *parser) previous() ast.Token {
	return p.toks[p.current-1]
}

func (p *parser) advance() ast.Token {
	if !p.atEnd() {
		p.current += 1
	}
	return p.previous()
}

func (p *parser) atEnd() bool {
	return p.peek().Kind == ast.EOF
}

func (pe ParseError) Error() string {
	if pe.tok.Kind == ast.EOF {
		return fmt.Sprintf("[line %v] Parse error at end of input: %v",
			pe.tok.Line, pe.message)
	}
	return fmt.Sprintf("[line %v] Parse error at '%v': %v",
		pe.tok.Line, pe.tok.Lexeme, pe.message)
}
