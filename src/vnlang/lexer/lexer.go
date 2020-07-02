package lexer

import (
	"io"
	"vnlang/scanner"
	"vnlang/token"
)

type Lexer struct {
	s    scanner.Scanner
	curr rune
}

func New(in io.Reader, filename string) *Lexer {
	var s scanner.Scanner
	s.Init(in)
	s.Filename = filename
	s.Whitespace ^= 1 << '\n' // don't skip new lines

	l := &Lexer{s: s}
	return l
}

func (l *Lexer) ResetLineCount() {
	l.s.ResetLineCount()
}

func (l *Lexer) NextToken() token.Token {
	var t token.Token

	l.readRune()
	switch l.curr {
	case '=':
		t = l.either('=', token.EQ, token.ASSIGN)
	case '+':
		t = l.token(token.PLUS)
	case '-':
		t = l.token(token.MINUS)
	case '*':
		t = l.token(token.ASTERISK)
	case '/':
		t = l.token(token.SLASH)
	case '.':
		t = l.token(token.DOT)
	case '%':
		t = l.token(token.MOD)
	case '&':
		t = l.either('&', token.AND, token.ILLEGAL)
	case '|':
		t = l.either('|', token.OR, token.ILLEGAL)
	case '!':
		t = l.either('=', token.NOT_EQ, token.BANG)
	case '>':
		t = l.either('=', token.GE, token.GT)
	case '<':
		t = l.either('=', token.LE, token.LT)
	case ';':
		t = l.token(token.SEMICOLON)
	case ',':
		t = l.token(token.COMMA)
	case ':':
		t = l.token(token.COLON)
	case '(':
		t = l.token(token.LPAREN)
	case ')':
		t = l.token(token.RPAREN)
	case '[':
		t = l.token(token.LBRACKET)
	case ']':
		t = l.token(token.RBRACKET)
	case '{':
		t = l.token(token.LBRACE)
	case '}':
		t = l.token(token.RBRACE)
	case '\n':
		t = l.token(token.NEWLINE)
	case scanner.Ident:
		lit := l.s.TokenText()
		if la := l.s.Peek(); la == '?' || la == '!' {
			l.readRune()
			lit += l.s.TokenText()
		}
		t = token.Token{
			Type:    token.LookupIdent(lit),
			Literal: lit,
		}
	case scanner.Int:
		lit := l.s.TokenText()
		t = token.Token{
			Type:    token.INT,
			Literal: lit,
		}
	case scanner.Float:
		lit := l.s.TokenText()
		t = token.Token{
			Type:    token.FLOAT,
			Literal: lit,
		}
	case scanner.String:
		lit := l.s.TokenText()
		t = token.Token{
			Type:    token.STRING,
			Literal: lit[1 : len(lit)-1],
		}
	case scanner.EOF:
		t = token.Token{Type: token.EOF, Literal: ""}
	default:
		lit := l.s.TokenText()
		t = token.Token{Type: token.ILLEGAL, Literal: lit}
	}
	t.Pos = l.s.Position
	return t
}

func (l *Lexer) readRune() {
	l.curr = l.s.Scan()
}

func (l *Lexer) token(ty token.TokenType) token.Token {
	lit := l.s.TokenText()
	return token.Token{Type: ty, Literal: lit}
}

func (l *Lexer) either(lookAhead rune, option, alternative token.TokenType) token.Token {
	lit := l.s.TokenText()
	if l.s.Peek() == lookAhead {
		l.readRune()
		lit += l.s.TokenText()
		return token.Token{Type: option, Literal: lit}
	} else {
		return token.Token{Type: alternative, Literal: lit}
	}
}
