package lexer

import (
	"loxogon/bad"
	Token "loxogon/token"
	"strconv"
)

type Lexer struct {
	source               string
	tokens               []Token.Token
	start, current, line int
}

var keywords = map[string]Token.TokenType{
	"and":    Token.AND,
	"class":  Token.CLASS,
	"else":   Token.ELSE,
	"false":  Token.FALSE,
	"for":    Token.FOR,
	"fun":    Token.FUN,
	"if":     Token.IF,
	"nil":    Token.NIL,
	"or":     Token.OR,
	"print":  Token.PRINT,
	"return": Token.RETURN,
	"super":  Token.SUPER,
	"this":   Token.THIS,
	"true":   Token.TRUE,
	"var":    Token.VAR,
	"while":  Token.WHILE,
}

func New(source string) Lexer {
	return Lexer{source: source, tokens: make([]Token.Token, 0, 10), start: 0, current: 0, line: 0}
}

func (l *Lexer) ScanTokens() []Token.Token {
	for !l.atEnd() {
		l.start = l.current
		l.scan()
	}

	l.tokens = append(l.tokens, Token.Token{Ty: Token.EOF, Lexeme: "", Literal: "", Line: l.line})
	return l.tokens
}

func (l *Lexer) atEnd() bool {
	return l.current >= len(l.source)
}

func (l *Lexer) advance() byte {
	b := l.source[l.current]
	l.current++
	return b
}

func (l *Lexer) scan() {
	c := l.advance()
	switch c {
	case '(':
		l.addToken(Token.LEFT_PAREN)
	case ')':
		l.addToken(Token.RIGHT_PAREN)
	case '{':
		l.addToken(Token.LEFT_BRACE)
	case '}':
		l.addToken(Token.RIGHT_BRACE)
	case ',':
		l.addToken(Token.COMMA)
	case '.':
		l.addToken(Token.DOT)
	case '-':
		l.addToken(Token.MINUS)
	case '+':
		l.addToken(Token.PLUS)
	case ';':
		l.addToken(Token.SEMICOLON)
	case '*':
		l.addToken(Token.STAR)
	case '!':
		if l.match('=') {
			l.addToken(Token.BANG_EQUAL)
		} else {
			l.addToken(Token.BANG)
		}
	case '=':
		if l.match('=') {
			l.addToken(Token.EQUAL_EQUAL)
		} else {
			l.addToken(Token.EQUAL)
		}
	case '<':
		if l.match('=') {
			l.addToken(Token.LESS_EQUAL)
		} else {
			l.addToken(Token.LESS)
		}
	case '>':
		if l.match('=') {
			l.addToken(Token.GREATER)
		} else {
			l.addToken(Token.GREATER_EQUAL)
		}
	case '/':
		if l.match('/') {
			for l.peek() != '\n' && !l.atEnd() {
				l.advance()
			}
		} else {
			l.addToken(Token.SLASH)
		}

	// Ignore whitespace.
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		l.line++
	case '"':
		l.string()
	default:
		if isDigit(c) {
			l.number()
		} else if isAlpha(c) {
			l.identifier()
		} else {
			bad.Raise(l.line, "Unexpected character.")
		}
	}
}

func (l *Lexer) addToken(t Token.TokenType) {
	l.addLiteral(t, nil)
}

func (l *Lexer) addLiteral(t Token.TokenType, literal any) {
	text := l.source[l.start:l.current]
	l.tokens = append(l.tokens, Token.Token{Ty: t, Literal: literal, Lexeme: text, Line: l.line})
}

func (l *Lexer) match(expected byte) bool {
	if l.atEnd() {
		return false
	}
	if l.source[l.current] != expected {
		return false
	}
	l.current++
	return true
}

func (l *Lexer) peek() byte {
	if l.atEnd() {
		return 0x00
	}
	return l.source[l.current]
}
func (l *Lexer) peekNext() byte {
	if l.current+1 >= len(l.source) {
		return 0x00
	}
	return l.source[l.current+1]
}

func (l *Lexer) string() {
	for l.peek() != '"' && !l.atEnd() {
		if l.peek() == '\n' {
			l.line++
		}
		l.advance()
	}
	if l.atEnd() {
		bad.Raise(l.line, "Unterminated string.")
		return
	}
	l.advance()
	l.addLiteral(Token.STRING, string(l.source[l.start+1:l.current-1]))
}

func (l *Lexer) number() {
	for isDigit(l.peek()) {
		l.advance()
	}
	if l.peek() == '.' && isDigit(l.peekNext()) {
		l.advance()
		for isDigit(l.peek()) {
			l.advance()
		}
	}
	float, err := strconv.ParseFloat(string(l.source[l.start:l.current]), 64)
	if err != nil {
		bad.Raise(l.line, "Could not parse number.")
	}
	l.addLiteral(Token.NUMBER, float)
}

func (l *Lexer) identifier() {
	for isAlphaNumeric(l.peek()) {
		l.advance()
	}

	text := string(l.source[l.start:l.current])
	ty, ok := keywords[text]
	if !ok {
		ty = Token.IDENTIFIER
	}
	l.addToken(ty)
}

func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func isAlpha(char byte) bool {
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		char == '_'
}

func isAlphaNumeric(char byte) bool {
	return isAlpha(char) || isDigit(char)
}
