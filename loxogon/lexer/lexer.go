package lexer

import (
	"fmt"
	"loxogon/token"
	"strconv"
)

type Lexer struct {
	source               string
	toks                 []token.Token
	start, current, line int
}

var keywords = map[string]token.TokenKind{
	"and":    token.AND,
	"class":  token.CLASS,
	"else":   token.ELSE,
	"false":  token.FALSE,
	"for":    token.FOR,
	"fun":    token.FUN,
	"if":     token.IF,
	"nil":    token.NIL,
	"or":     token.OR,
	"print":  token.PRINT,
	"return": token.RETURN,
	"super":  token.SUPER,
	"this":   token.THIS,
	"true":   token.TRUE,
	"var":    token.VAR,
	"while":  token.WHILE,
}

func New(source string) Lexer {
	return Lexer{source: source,
		toks:  make([]token.Token, 0, 10),
		start: 0, current: 0, line: 1}
}

func (l *Lexer) ScanTokens() ([]token.Token, error) {
	for !l.atEnd() {
		l.start = l.current
		if err := l.scan(); err != nil {
			return nil, fmt.Errorf("could not lex string: %w", err)
		}
	}

	l.toks = append(l.toks,
		token.Token{Kind: token.EOF,
			Lexeme: "", Literal: "", Line: l.line})
	return l.toks, nil
}

func (l *Lexer) atEnd() bool {
	return l.current >= len(l.source)
}

func (l *Lexer) advance() byte {
	b := l.source[l.current]
	l.current++
	return b
}

func (l *Lexer) scan() error {
	c := l.advance()
	switch c {
	case '(':
		l.addToken(token.LEFT_PAREN)
	case ')':
		l.addToken(token.RIGHT_PAREN)
	case '{':
		l.addToken(token.LEFT_BRACE)
	case '}':
		l.addToken(token.RIGHT_BRACE)
	case ',':
		l.addToken(token.COMMA)
	case '.':
		l.addToken(token.DOT)
	case '-':
		l.addToken(token.MINUS)
	case '+':
		l.addToken(token.PLUS)
	case ';':
		l.addToken(token.SEMICOLON)
	case '*':
		l.addToken(token.STAR)
	case '!':
		if l.match('=') {
			l.addToken(token.BANG_EQUAL)
		} else {
			l.addToken(token.BANG)
		}
	case '=':
		if l.match('=') {
			l.addToken(token.EQUAL_EQUAL)
		} else {
			l.addToken(token.EQUAL)
		}
	case '<':
		if l.match('=') {
			l.addToken(token.LESS_EQUAL)
		} else {
			l.addToken(token.LESS)
		}
	case '>':
		if !l.match('=') {
			l.addToken(token.GREATER)
		} else {
			l.addToken(token.GREATER_EQUAL)
		}
	case '/':
		if l.match('/') {
			for l.peek() != '\n' && !l.atEnd() {
				l.advance()
			}
		} else {
			l.addToken(token.SLASH)
		}

	// Ignore whitespace.
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		l.line++
	case '"':
		if err := l.string(); err != nil {
			return err
		}
	default:
		if isDigit(c) {
			if err := l.number(); err != nil {
				return err
			}
		} else if isAlpha(c) {
			l.identifier()
		} else {
			return l.error(l.line, "unexpected character")
		}
	}
	return nil
}

func (l *Lexer) addToken(k token.TokenKind) {
	l.addLiteral(k, nil)
}

func (l *Lexer) addLiteral(k token.TokenKind, literal any) {
	text := l.source[l.start:l.current]
	l.toks = append(l.toks, token.Token{Kind: k, Literal: literal, Lexeme: text, Line: l.line})
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

func (l *Lexer) string() error {
	for l.peek() != '"' && !l.atEnd() {
		if l.peek() == '\n' {
			l.line++
		}
		l.advance()
	}
	if l.atEnd() {
		return l.error(l.line, "unterminated string")
	}
	l.advance()
	l.addLiteral(token.STRING, string(l.source[l.start+1:l.current-1]))
	return nil
}

func (l *Lexer) number() error {
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
		// TODO: Unreachable?
		return l.error(l.line, "could not parse number")
	}
	l.addLiteral(token.NUMBER, float)
	return nil
}

func (l *Lexer) identifier() {
	for isAlphaNumeric(l.peek()) {
		l.advance()
	}

	text := string(l.source[l.start:l.current])
	kind, ok := keywords[text]
	if !ok {
		kind = token.IDENTIFIER
	}
	l.addToken(kind)
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

func (l *Lexer) error(line int, message string) error {
	return fmt.Errorf("[line %d] %s", line, message)
}
