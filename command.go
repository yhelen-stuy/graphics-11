package main

type Command interface {
	CommandType() string
}

type PushCommand struct{}

func (c PushCommand) CommandType() string { return "push" }

type PopCommand struct{}

func (c PopCommand) CommandType() string { return "pop" }

type ScaleCommand struct {
	x    float64
	y    float64
	z    float64
	knob string
}

func (c ScaleCommand) CommandType() string { return "scale" }

type MoveCommand struct {
	x    float64
	y    float64
	z    float64
	knob string
}

func (c MoveCommand) CommandType() string { return "move" }

type RotateCommand struct {
	axis  string
	angle float64
	knob  string
}

func (c RotateCommand) CommandType() string { return "rotate" }

type BoxCommand struct {
	x      float64
	y      float64
	z      float64
	height float64
	width  float64
	depth  float64
}

func (c BoxCommand) CommandType() string { return "box" }

type SphereCommand struct {
	center []float64
	radius float64
}

func (c SphereCommand) CommandType() string { return "sphere" }

type TorusCommand struct {
	center []float64
	r1     float64
	r2     float64
}

func (c TorusCommand) CommandType() string { return "torus" }

type LineCommand struct {
	p1 []float64
	p2 []float64
}

func (c LineCommand) CommandType() string { return "line" }

type SaveCommand struct {
	filename string
}

func (c SaveCommand) CommandType() string { return "save" }

type DisplayCommand struct{}

func (c DisplayCommand) CommandType() string { return "display" }
