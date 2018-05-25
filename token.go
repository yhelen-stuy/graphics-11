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
	WEB

	PUSH
	POP

	SAVE
	DISPLAY
	GENERATERAYFILES
	SHADING
	SHADINGTYPE
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
	"light":             LIGHT,
	"ambient":           AMBIENT,
	"camera":            CAMERA,
	"constants":         CONSTANTS,
	"save_coord_system": SAVECOORDS,
	"mesh":              MESH,
	"set":               SET,
	"tween":             TWEEN,
	"generate_rayfiles": GENERATERAYFILES,
	"shading":           SHADING,
	"focal":             FOCAL,
	"setknobs":          SETKNOBS,
}

func FindOp(op string) TokenType {
	if ttype, isOp := ops[op]; isOp {
		return ttype
	}
	// invalid operation
	return T_ILLEGAL
}
