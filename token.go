package main

type TokenType int

// TODO: also comment this code lol
//maybe i should move this
const (
	T_ILLEGAL TokenType = iota
	T_ERROR
	T_COMMENT    // 2
	T_INT        // 3
	T_FLOAT      // 4
	T_STRING     // 5
	T_NEWLINE    // 6
	T_IDENTIFIER // 7
	T_EOF        // 8

	LIGHT
	CONSTANTS
	CAMERA
	AMBIENT
	TORUS
	SPHERE
	BOX
	LINE
	MOVE
	SCALE
	ROTATE
	BASENAME
	FRAMES
	VARY

	PUSH
	POP

	SAVE
	DISPLAY
)

type Token struct {
	ttype TokenType
	val   string
}

var ops = map[string]TokenType{
	"torus":    TORUS,
	"sphere":   SPHERE,
	"box":      BOX,
	"line":     LINE,
	"move":     MOVE,
	"scale":    SCALE,
	"rotate":   ROTATE,
	"basename": BASENAME,
	"frames":   FRAMES,
	"vary":     VARY,
	"push":     PUSH,
	"pop":      POP,
	"save":     SAVE,
	"display":  DISPLAY,

	// UNIMPLEMENTED
	"light":     LIGHT,
	"ambient":   AMBIENT,
	"camera":    CAMERA,
	"constants": CONSTANTS,
}

func FindOp(op string) TokenType {
	if ttype, isOp := ops[op]; isOp {
		return ttype
	}
	// invalid operation
	return T_ILLEGAL
}
