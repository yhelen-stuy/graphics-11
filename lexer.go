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

type TokenType int

// TODO: also comment this code lol
//maybe i should move this
const (
	T_ILLEGAL TokenType = iota
	T_ERROR
	T_COMMENT
	T_INT
	T_FLOAT
	T_STRING
	T_NEWLINE
	T_IDENTIFIER
	T_EOF

	LIHT
	CONSTANTS
	SAVECOORDS
	CAMERA
	AMBIENT
	TORUS
	SPHERE
	BOX
	LINE
	MESH
	TEXTURE
	SET
	MOVE
	SCALE
	ROTATE
	BASENAME
	SAVEKNOBS
	TWEEN
	FRAMES
	VARY
	SETKNOBS
	FOCAL
	DISPLAY
	WEB

	PUSH
	POP

	SAVE
	GENERATERAYFILES
	SHADING
	SHADINGTYPE
)

type Token struct {
	ttype TokenType
	val   string
}

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

func Lex(input string) (*Lexer, chan Token) {
	l := &Lexer{
		input:  input,
		tokens: make(chan Token),
	}
	go l.run()
	return l, l.tokens
}

// are these all the cases????
func lexBase(l *Lexer) stateFn {
	switch r := l.next(); {
	case r == eof:
		l.emit(T_EOF)
		return nil
	case r == '\n':
		l.emit(T_NEWLINE)
	case r == '/':
		if l.peek() == '/' {
			return lexComment
		}
		return lexString
	case unicode.IsSpace(r):
		l.ignore()
	case r == '-' || 0 <= r && r <= 9:
		l.backup()
		return lexNumber
		// Is this what I want?? hb letter/number only idrk
	case unicode.IsPrint(r):
		return lexString
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
	if unicode.IsLetter(l.peek()) {
		return l.errorf("Invalid number")
	}
	l.emit(ttype)
	return lexBase
}

func lexString(l *Lexer) stateFn {
	// need to check:
	// legal chars, if not has to be a string, else might be identifier
	// continue reading until whitespace
	return lexBase
}

func lexComment(l *Lexer) stateFn {
	// read until newline
	return lexBase
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
