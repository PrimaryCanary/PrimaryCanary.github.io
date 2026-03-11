package ast

import "fmt"

type TokenKind int

const (
	// Single-character tokens.
	LEFT_PAREN TokenKind = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	// One or two character tokens.
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// Literals.
	IDENTIFIER
	STRING
	NUMBER

	// Keywords.
	AND
	CLASS
	ELSE
	FALSE
	FUN_TOK
	FOR_TOK
	IF_TOK
	NIL
	OR
	PRINT_TOK
	RETURN_TOK
	SUPER
	THIS
	TRUE
	VAR_TOK
	WHILE_TOK

	EOF
)

var token_names = [...]string{
	// Single-character tokens.
	"LEFT_PAREN",
	"RIGHT_PAREN",
	"LEFT_BRACE",
	"RIGHT_BRACE",
	"COMMA",
	"DOT",
	"MINUS",
	"PLUS",
	"SEMICOLON",
	"SLASH",
	"STAR",

	// One or two character tokens.
	"BANG",
	"BANG_EQUAL",
	"EQUAL",
	"EQUAL_EQUAL",
	"GREATER",
	"GREATER_EQUAL",
	"LESS",
	"LESS_EQUAL",

	// Literals.
	"IDENTIFIER",
	"STRING",
	"NUMBER",

	// Keywords.
	"AND",
	"CLASS",
	"ELSE",
	"FALSE",
	"FUN",
	"FOR",
	"IF",
	"NIL",
	"OR",
	"PRINT",
	"RETURN",
	"SUPER",
	"THIS",
	"TRUE",
	"VAR",
	"WHILE",

	"EOF",
}

var keywords = map[string]TokenKind{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR_TOK,
	"fun":    FUN_TOK,
	"if":     IF_TOK,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT_TOK,
	"return": RETURN_TOK,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR_TOK,
	"while":  WHILE_TOK,
}

func IsKeyword(text string) (TokenKind, bool) {
	kw, ok := keywords[text]
	return kw, ok
}
func (t TokenKind) String() string {
	return token_names[t]
}

type Token struct {
	Kind    TokenKind
	Lexeme  string
	Literal any
	Line    int
}

func (t Token) String() string {
	return fmt.Sprintf("%s %s %v", t.Kind, t.Lexeme, t.Literal)
}
