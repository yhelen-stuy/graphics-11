package main

import "fmt"

type Command interface {
	CommandString() string
}

type PushCommand struct{}

func (c PushCommand) CommandString() string { return "push" }

type PopCommand struct{}

func (c PopCommand) CommandString() string { return "pop" }

type ScaleCommand struct {
	x    float64
	y    float64
	z    float64
	knob string
}

func (c ScaleCommand) CommandString() string {
	return fmt.Sprintf("scale %.1f %.1f %.1f %s", c.x, c.y, c.z, c.knob)
}

type MoveCommand struct {
	x    float64
	y    float64
	z    float64
	knob string
}

func (c MoveCommand) CommandString() string {
	return fmt.Sprintf("move %.1f %.1f %.1f %s", c.x, c.y, c.z, c.knob)
}

type RotateCommand struct {
	axis  string
	angle float64
	knob  string
}

func (c RotateCommand) CommandString() string {
	return fmt.Sprintf("rotate %s %.1f %s", c.axis, c.angle, c.knob)
}

type BoxCommand struct {
	cons   string
	x      float64
	y      float64
	z      float64
	height float64
	width  float64
	depth  float64
}

func (c BoxCommand) CommandString() string {
	return fmt.Sprintf("box %s %.1f %.1f %.1f %.1f %.1f %.1f", c.cons, c.x, c.y, c.z, c.height, c.width, c.depth)
}

type SphereCommand struct {
	cons   string
	center []float64
	radius float64
}

func (c SphereCommand) CommandString() string {
	return fmt.Sprintf("sphere %s %.1f %.1f %.1f %.1f", c.cons, c.center[0], c.center[1], c.center[2], c.radius)
}

type TorusCommand struct {
	cons   string
	center []float64
	r1     float64
	r2     float64
}

func (c TorusCommand) CommandString() string {
	return fmt.Sprintf("torus %s %.1f %.1f %.1f %.1f %.1f", c.cons, c.center[0], c.center[1], c.center[2], c.r1, c.r2)
}

type LineCommand struct {
	cons string
	p1   []float64
	p2   []float64
}

func (c LineCommand) CommandString() string {
	return fmt.Sprintf("line %s %.1f %.1f %.1f %.1f %.1f %.1f", c.cons, c.p1[0], c.p1[1], c.p1[2], c.p2[0], c.p2[1], c.p2[2])
}

type SaveCommand struct {
	filename string
}

func (c SaveCommand) CommandString() string {
	return fmt.Sprintf("save %s", c.filename)
}

type DisplayCommand struct{}

func (c DisplayCommand) CommandString() string { return "display" }
