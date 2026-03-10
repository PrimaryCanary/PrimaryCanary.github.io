package ast

import (
	"fmt"
	"strings"
)

type ExprKind int
type StmtKind int

const (
	BINARY ExprKind = iota
	UNARY
	LITERAL
	GROUPING
	VARIABLE
	ASSIGN
)
const (
	EXPR StmtKind = iota
	PRINT
	VAR
	VAR_UNINIT
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
	Kind  StmtKind
	Name  Token
	Child Expr
	Stmts []Stmt
}

type LoxObject struct {
	Value any
}

func NewLiteral(data any) Expr {
	return Expr{Kind: LITERAL, Data: data}
}

func NewBinary(operator Token, left, right Expr) Expr {
	return Expr{Kind: BINARY, Tok: operator, Children: []Expr{left, right}}
}

func NewUnary(operator Token, e Expr) Expr {
	return Expr{Kind: UNARY, Tok: operator, Children: []Expr{e}}
}

func NewGrouping(e Expr) Expr {
	return Expr{Kind: GROUPING, Children: []Expr{e}}
}

func NewVariable(t Token) Expr {
	return Expr{Kind: VARIABLE, Tok: t}
}

func NewAssign(t Token, e Expr) Expr {
	return Expr{Kind: ASSIGN, Tok: t, Children: []Expr{e}}
}

func NewExprStmt(e Expr) Stmt {
	return Stmt{Kind: EXPR, Child: e}
}

func NewPrintStmt(e Expr) Stmt {
	return Stmt{Kind: PRINT, Child: e}
}

func NewBlock(stmts []Stmt) Stmt {
	return Stmt{Kind: BLOCK, Stmts: stmts}
}

func NewVarDecl(name Token, e Expr) Stmt {
	return Stmt{Kind: VAR, Name: name, Child: e}
}

func NewVarDeclUninit(name Token) Stmt {
	return Stmt{Kind: VAR_UNINIT, Name: name}
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
		return fmt.Sprintf("(var %v=%v)", s.Name.Lexeme, s.Child)
	case VAR_UNINIT:
		return fmt.Sprintf("(var %v)", s.Name.Lexeme)
	case BLOCK:
		return fmt.Sprintf("{%v}", StmtsToString(s.Stmts))
	}
	return ""
}

func (lo LoxObject) String() string {
	if lo.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", lo.Value)
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
