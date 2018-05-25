package main

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// ************************
//  Credits for design of
//  this lexer:
//  Golang Talk:
//  Lexical Scanning in Go
//  by Rob Pike
// ************************
type Lexer struct {
	input  string
	start  int
	pos    int
	width  int
	tokens chan Token
	state  stateFn
}

var eof = rune(0)

type stateFn func(*Lexer) stateFn

func (l *Lexer) run() {
	for state := lexBase; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

func Lex(input string) *Lexer {
	l := &Lexer{
		input:  input,
		tokens: make(chan Token),
	}
	go l.run()
	return l
}

// are these all the cases????
func lexBase(l *Lexer) stateFn {
	switch r := l.next(); {
	case r == eof:
		fmt.Println("Reached eof")
		l.emit(T_EOF)
		return nil
	case r == '\n':
		l.emit(T_NEWLINE)
		return lexBase
	case r == '/':
		if l.peek() == '/' {
			return lexComment
		}
		return lexText
	case unicode.IsSpace(r):
		l.ignore()
		return lexBase
	case r == '.' || r == '-' || strings.IndexRune("0123456789", r) >= 0:
		l.backup()
		return lexNumber
		// Is this what I want?? hb letter/number only idrk
	case unicode.IsPrint(r):
		return lexText
	}
	// probably not the behavior i want lol
	return nil
}

func lexNumber(l *Lexer) stateFn {
	var ttype TokenType
	l.accept("-")
	digits := "0123456789"
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
		ttype = T_FLOAT
	} else {
		ttype = T_INT
	}
	if unicode.IsLetter(l.peek()) || unicode.IsPunct(l.peek()) {
		return l.errorf("Invalid number")
	}
	l.emit(ttype)
	return lexBase
}

func lexText(l *Lexer) stateFn {
	// need to check:
	// legal chars, if not has to be a text, else might be identifier
	// continue reading until whitespace
	r := l.next()
	// stop lexing when encountering whitespace
	for unicode.IsPrint(r) && !unicode.IsSpace(r) {
		r = l.next()
	}
	l.backup()
	ttype := FindOp(l.input[l.start:l.pos])
	if ttype != T_ILLEGAL {
		l.emit(T_IDENTIFIER)
	} else {
		// need to check if identifier
		l.emit(T_STRING)
	}
	return lexBase
}

func lexComment(l *Lexer) stateFn {
	r := l.next()
	if r == '\n' {
		l.ignore()
		return lexBase
	}
	if r == eof {
		return nil
	}
	return lexComment
}

func (l *Lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, width := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = width
	l.pos += width
	return r
}

func (l *Lexer) NextToken() Token {
	return <-l.tokens
}

func (l *Lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *Lexer) acceptRun(valid string) {
	for l.accept(valid) {
	}
}

func (l *Lexer) backup() {
	l.pos -= l.width
}

func (l *Lexer) ignore() {
	l.start = l.pos
}

func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *Lexer) errorf(err string) stateFn {
	// More robust error message lol
	l.tokens <- Token{
		ttype: T_ERROR,
		val:   fmt.Sprintf("Syntax error: %s", err),
	}
	return nil
}

func (l *Lexer) emit(ttype TokenType) {
	l.tokens <- Token{
		ttype: ttype,
		val:   l.input[l.start:l.pos],
	}
	l.start = l.pos
}

func (t Token) String() {
	fmt.Sprintf("(%d: %s)", t.ttype, t.val)
}
