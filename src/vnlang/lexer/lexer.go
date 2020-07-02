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

func New(in io.Reader) *Lexer {
	var s scanner.Scanner
	s.Init(in)
	s.Whitespace ^= 1 << '\n' // don't skip new lines

	l := &Lexer{s: s}
	return l
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
		// p := l.s.Pos()
		lit := l.s.TokenText()
		// col := p.Column
		if la := l.s.Peek(); la == '?' || la == '!' {
			l.readRune()
			lit += l.s.TokenText()
			// col += 1
		}
		t = token.Token{
			Type:    token.LookupIdent(lit),
			Literal: lit,
		}
	case scanner.Int:
		// p := l.s.Pos()
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
		// p := l.s.Pos()
		lit := l.s.TokenText()
		t = token.Token{
			Type:    token.STRING,
			Literal: lit[1 : len(lit)-1],
		}
	case scanner.EOF:
		// p := l.s.Pos()
		t = token.Token{Type: token.EOF, Literal: ""}
	default:
		// p := l.s.Pos()
		lit := l.s.TokenText()
		t = token.Token{Type: token.ILLEGAL, Literal: lit}
	}
	return t
}

func (l *Lexer) readRune() {
	l.curr = l.s.Scan()
}

func (l *Lexer) token(ty token.TokenType) token.Token {
	// p := l.s.Pos()
	lit := l.s.TokenText()
	return token.Token{Type: ty, Literal: lit}
}

func (l *Lexer) either(lookAhead rune, option, alternative token.TokenType) token.Token {
	// p := l.s.Pos()
	lit := l.s.TokenText()
	// col := p.Column
	if l.s.Peek() == lookAhead {
		l.readRune()
		lit += l.s.TokenText()
		// col += 1
		return token.Token{Type: option, Literal: lit}
	} else {
		return token.Token{Type: alternative, Literal: lit}
	}
}
