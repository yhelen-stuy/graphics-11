package main

import (
	// "bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	// "strings"
)

var knobs map[string][]float64

func init() {
	knobs = make(map[string][]float64)
}

type Parser struct {
	lexer    *Lexer
	stack    *Stack
	poly     *Matrix
	edge     *Matrix
	trans    *Matrix
	image    *Image
	history  []Token
	frames   int
	basename string
}

func MakeParser() *Parser {
	trans := MakeMatrix(4, 4)
	trans.Ident()
	return &Parser{
		stack:    MakeStack(),
		poly:     MakeMatrix(4, 0),
		edge:     MakeMatrix(4, 0),
		trans:    trans,
		image:    MakeImage(500, 500),
		history:  make([]Token, 1),
		frames:   -1,
		basename: "",
	}
}

func ParseFile(filename string) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.New("Couldn't read file")
	}
	str := string(buf)
	p := MakeParser()
	c, _ := p.parseString(str)
	p.runCommands(c)
	return nil
}

func (p *Parser) parseString(input string) ([]Command, error) {
	p.lexer = Lex(input)
	commands := make([]Command, 0)
	for {
		t := p.next()
		switch t.ttype {
		case T_EOF:
			return commands, nil
		case T_IDENTIFIER:
			switch tt := FindOp(t.val); tt {
			case PUSH:
				c := PushCommand{}
				commands = append(commands, c)
			case POP:
				c := PopCommand{}
				commands = append(commands, c)
			case SCALE:
				c := ScaleCommand{
					x: p.nextFloat(),
					y: p.nextFloat(),
					z: p.nextFloat(),
				}
				commands = append(commands, c)
			case MOVE:
				c := MoveCommand{
					x: p.nextFloat(),
					y: p.nextFloat(),
					z: p.nextFloat(),
				}
				commands = append(commands, c)
			case ROTATE:
				c := RotateCommand{}
				axis := p.nextRequired([]TokenType{T_STRING})
				c.axis = axis.val
				c.angle = p.nextFloat()
				commands = append(commands, c)
			case BOX:
				c := BoxCommand{
					x:      p.nextFloat(),
					y:      p.nextFloat(),
					z:      p.nextFloat(),
					height: p.nextFloat(),
					width:  p.nextFloat(),
					depth:  p.nextFloat(),
				}
				commands = append(commands, c)
			case SPHERE:
				c := SphereCommand{
					center: []float64{p.nextFloat(), p.nextFloat(), p.nextFloat()},
					radius: p.nextFloat(),
				}
				commands = append(commands, c)
			case TORUS:
				c := TorusCommand{
					center: []float64{p.nextFloat(), p.nextFloat(), p.nextFloat()},
					r1:     p.nextFloat(),
					r2:     p.nextFloat(),
				}
				commands = append(commands, c)
			case LINE:
				c := LineCommand{
					p1: []float64{p.nextFloat(), p.nextFloat(), p.nextFloat()},
					p2: []float64{p.nextFloat(), p.nextFloat(), p.nextFloat()},
				}
				commands = append(commands, c)
			case VARY:
				if p.frames <= 1 {
					return nil, errors.New("Frames not set")
				}
				name := p.nextRequired([]TokenType{T_STRING}).val
				knob, isKnob := knobs[name]
				if !isKnob {
					knobs[name] = make([]float64, p.frames)
				}
				startFrame := p.nextInt()
				endFrame := p.nextInt()
				if startFrame < 0 || endFrame < startFrame || endFrame < 0 {
					return nil, errors.New("Invalid frame values in vary")
				}
				startVal := p.nextFloat()
				endVal := p.nextFloat()
				// slope
				m := (endVal - startVal) / float64(endFrame-startFrame+1)
				for i := startFrame; i < endFrame; i++ {
					knob[i] = startVal
					startVal += m
				}
				fmt.Println(knobs)
			case FRAMES:
				if p.frames != -1 {
					fmt.Println("Warning: Setting frames multiple times")
				}
				p.frames = p.nextInt()
				if p.frames <= 0 {
					panic(errors.New("Frames must be positive"))
				}
			case BASENAME:
				if p.basename != "" {
					fmt.Println("Warning: Setting basename multiple times")
				}
				basename := p.nextRequired([]TokenType{T_STRING})
				p.basename = basename.val
			case SAVE:
				c := SaveCommand{}
				filename := p.nextRequired([]TokenType{T_STRING})
				c.filename = filename.val
				commands = append(commands, c)
			case DISPLAY:
				c := DisplayCommand{}
				commands = append(commands, c)
			default:
				fmt.Println(tt, t.val)
			}
		default:
			continue
		}
	}
}

func (p *Parser) runCommands(commands []Command) {
	for _, com := range commands {
		switch com.(type) {
		case PushCommand:
			p.trans = p.stack.Peek()
			if p.trans != nil {
				p.stack.Push(p.trans.Copy())
			}
		case PopCommand:
			p.stack.Pop()
		case ScaleCommand:
			c := com.(ScaleCommand)
			scale := MakeScale(c.x, c.y, c.z)
			p.trans = p.stack.Pop()
			p.trans, _ = scale.Mult(p.trans)
			p.stack.Push(p.trans.Copy())
		case MoveCommand:
			c := com.(MoveCommand)
			translate := MakeTranslate(c.x, c.y, c.z)
			p.trans = p.stack.Pop()
			p.trans, _ = translate.Mult(p.trans)
			p.stack.Push(p.trans.Copy())
		case RotateCommand:
			c := com.(RotateCommand)
			switch c.axis {
			case "x":
				p.trans = p.stack.Pop()
				rot := MakeRotX(c.angle)
				p.trans, _ = rot.Mult(p.trans)
				p.stack.Push(p.trans.Copy())
			case "y":
				p.trans = p.stack.Pop()
				rot := MakeRotY(c.angle)
				p.trans, _ = rot.Mult(p.trans)
				p.stack.Push(p.trans.Copy())
			case "z":
				p.trans = p.stack.Pop()
				rot := MakeRotZ(c.angle)
				p.trans, _ = rot.Mult(p.trans)
				p.stack.Push(p.trans.Copy())
			default:
				fmt.Println("Rotate fail")
				continue
			}
		case BoxCommand:
			c := com.(BoxCommand)
			p.poly.AddBox(c.x, c.y, c.z, c.height, c.width, c.depth)
			p.poly, _ = p.poly.Mult(p.stack.Peek())
			p.image.DrawPolygons(p.poly, Color{r: 0, b: 255, g: 0})
			p.poly = MakeMatrix(4, 0)
		case SphereCommand:
			c := com.(SphereCommand)
			p.poly.AddSphere(c.center[0], c.center[1], c.center[2], c.radius)
			p.poly, _ = p.poly.Mult(p.stack.Peek())
			p.image.DrawPolygons(p.poly, Color{r: 0, b: 255, g: 0})
			p.poly = MakeMatrix(4, 0)
		case TorusCommand:
			c := com.(TorusCommand)
			p.poly.AddTorus(c.center[0], c.center[1], c.center[2], c.r1, c.r2)
			p.poly, _ = p.poly.Mult(p.stack.Peek())
			p.image.DrawPolygons(p.poly, Color{r: 0, b: 255, g: 0})
			p.poly = MakeMatrix(4, 0)
		case LineCommand:
			c := com.(LineCommand)
			p.poly.AddEdge(c.p1[0], c.p1[1], c.p1[2], c.p2[0], c.p1[1], c.p2[2])
			p.edge, _ = p.edge.Mult(p.stack.Peek())
			p.image.DrawLines(p.edge, Color{r: 255, b: 0, g: 0})
			p.edge = MakeMatrix(4, 0)
		case SaveCommand:
			c := com.(SaveCommand)
			p.image.SavePPM("temp")
			p.image.ConvertPNG("temp", c.filename)
		case DisplayCommand:
			p.image.Display()
		}
	}
}

func (p *Parser) nextRequired(ttypes []TokenType) Token {
	t := p.next()
	for _, ttype := range ttypes {
		if t.ttype == ttype {
			return t
		}
	}
	fmt.Println(ttypes)
	panic(fmt.Errorf("Unexpected type received: got %d with value: %s", t.ttype, t.val))
}

func (p *Parser) nextInt() int {
	t := p.nextRequired([]TokenType{T_INT})
	i, _ := strconv.Atoi(t.val)
	return i
}

func (p *Parser) nextFloat() float64 {
	t := p.nextRequired([]TokenType{T_INT, T_FLOAT})
	i, _ := strconv.ParseFloat(t.val, 64)
	return i
}

// returns next token
func (p *Parser) next() Token {
	length := len(p.history)
	if length == 0 {
		return p.lexer.NextToken()
	}
	t := p.history[length-1]
	p.history = p.history[:length-1]
	return t
}

func (p *Parser) backup(t Token) {
	p.history = append(p.history, t)
}

func (p *Parser) peek() Token {
	t := p.next()
	p.backup(t)
	return t
}
