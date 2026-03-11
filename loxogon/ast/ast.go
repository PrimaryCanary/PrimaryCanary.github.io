package ast

import (
	"fmt"
	"slices"
	"strings"
)

type ExprKind int
type StmtKind int

const (
	BINARY ExprKind = iota
	LOGICAL
	UNARY
	LITERAL
	GROUPING
	VARIABLE
	ASSIGN
	CALL
)
const (
	EXPR StmtKind = iota
	PRINT
	VAR
	VAR_UNINIT
	IF
	WHILE
	FOR
	FUN
	RETURN
	RETURN_EMPTY
	BLOCK
)

// Fat struct representation of expressions
type Expr struct {
	Kind     ExprKind
	Tok      Token
	Data     any
	Children []Expr
}

// Fat struct representation of expressions
type Stmt struct {
	Kind   StmtKind
	Tokens []Token
	Child  Expr
	Stmts  []Stmt
}

// TODO see if unified struct repr works
//	type TestCombinedStmtExpr struct {
//		Kind  StmtKind
//		Name  Token
//		Data  any
//		Stmts []Stmt
//	}

func NewLiteral(data any) Expr {
	return Expr{Kind: LITERAL, Data: data}
}

func NewBinary(operator Token, left, right Expr) Expr {
	return Expr{Kind: BINARY, Tok: operator, Children: []Expr{left, right}}
}

func NewLogical(operator Token, left, right Expr) Expr {
	return Expr{Kind: LOGICAL, Tok: operator, Children: []Expr{left, right}}
}

func NewUnary(operator Token, e Expr) Expr {
	return Expr{Kind: UNARY, Tok: operator, Children: []Expr{e}}
}

func NewGrouping(e Expr) Expr {
	return Expr{Kind: GROUPING, Children: []Expr{e}}
}

func NewVariable(name Token) Expr {
	return Expr{Kind: VARIABLE, Tok: name}
}

func NewAssign(name Token, e Expr) Expr {
	return Expr{Kind: ASSIGN, Tok: name, Children: []Expr{e}}
}

func NewCall(callee Expr, paren Token, args []Expr) Expr {
	children := slices.Concat([]Expr{callee}, args)
	return Expr{Kind: CALL, Tok: paren, Children: children}
}

func NewExprStmt(e Expr) Stmt {
	return Stmt{Kind: EXPR, Child: e}
}

func NewPrintStmt(e Expr) Stmt {
	return Stmt{Kind: PRINT, Child: e}
}

func NewIf(condition Expr, branches ...Stmt) Stmt {
	switch len(branches) {
	case 1:
		return Stmt{Kind: IF, Child: condition, Stmts: []Stmt{branches[0]}}
	case 2:
		return Stmt{Kind: IF, Child: condition, Stmts: []Stmt{branches[0], branches[1]}}
	default:
		panic("Too many branches specified in if statement creation")
	}
}

func NewWhile(condition Expr, body Stmt) Stmt {
	return Stmt{Kind: WHILE, Child: condition, Stmts: []Stmt{body}}
}

func NewFun(name Token, params []Token, body []Stmt) Stmt {
	toks := slices.Concat([]Token{name}, params)
	b := NewBlock(body)
	return Stmt{Kind: FUN, Tokens: toks, Stmts: []Stmt{b}}
}

func NewReturn(keyword Token, e Expr) Stmt {
	return Stmt{Kind: RETURN, Tokens: []Token{keyword}, Child: e}
}

// TODO investigate SENTINEL ExprKind to recognize uninit exprs
func NewReturnEmpty(keyword Token) Stmt {
	return Stmt{Kind: RETURN_EMPTY, Tokens: []Token{keyword}}
}

func NewBlock(stmts []Stmt) Stmt {
	return Stmt{Kind: BLOCK, Stmts: stmts}
}
func NewVarDecl(name Token, e Expr) Stmt {
	return Stmt{Kind: VAR, Tokens: []Token{name}, Child: e}
}

func NewVarDeclUninit(name Token) Stmt {
	return Stmt{Kind: VAR_UNINIT, Tokens: []Token{name}}
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
	case VARIABLE:
		return fmt.Sprintf("%s", e.Tok.Lexeme)
	case ASSIGN:
		return fmt.Sprintf("(%v = %v)", e.Tok.Lexeme, e.Children[0])
	}
	// TODO Unreachable
	return ""
}

func (s Stmt) String() string {
	switch s.Kind {
	case EXPR:
		return fmt.Sprintf("%v", s.Child)
	case PRINT:
		return fmt.Sprintf("(print %v)", s.Child)
	case VAR:
		return fmt.Sprintf("(var %v=%v)", s.Tokens[0].Lexeme, s.Child)
	case VAR_UNINIT:
		return fmt.Sprintf("(var %v)", s.Tokens[0].Lexeme)
	case BLOCK:
		return fmt.Sprintf("{%v}", StmtsToString(s.Stmts))
	}
	return ""
}

func StmtsToString(stmts []Stmt) string {
	var sb strings.Builder
	for _, stmt := range stmts {
		sb.WriteString(stmt.String())
		sb.WriteString(" ")
	}
	// Remove trailing space
	return sb.String()[:sb.Len()-1]
}
