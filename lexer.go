package main

import (
	"fmt"
	"utf8"
)

// ************************
//  Credits for design of
//  this lexer:
//  Golang Talk:
//  Lexical Scanning in Go
//  by Rob Pike
// ************************

type TokenType int

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
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.items())
}

func Lex(name, input string) (*lexer, chan Token) {
	l := &Lexer{
		name:   name,
		input:  input,
		tokens: make(chan Token),
	}
	go l.run()
	return l, l.tokens
}

func (t Token) String() {
	fmt.Sprintf("(%d: %s)", t.ttype, t.val)
}

func (l *Lexer) next() (rune int) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = width
	l.pos += width
	return r
}
