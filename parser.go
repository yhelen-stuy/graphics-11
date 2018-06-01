package main

import (
	// "bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	// "strings"
)

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
	for i := range c {
		fmt.Println(c[i].CommandType())
	}
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
				axis, err := p.nextRequired([]TokenType{T_STRING})
				if err != nil {
					return nil, err
				}
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
				basename, err := p.nextRequired([]TokenType{T_STRING})
				if err != nil {
					return nil, err
				}
				p.basename = basename.val
			case SAVE:
				c := SaveCommand{}
				filename, err := p.nextRequired([]TokenType{T_STRING})
				if err != nil {
					return nil, err
				}
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

func (p *Parser) nextRequired(ttypes []TokenType) (Token, error) {
	t := p.next()
	for _, ttype := range ttypes {
		if t.ttype == ttype {
			return t, nil
		}
	}
	fmt.Println(ttypes)
	panic(fmt.Errorf("Unexpected type received: got %d with value: %s", t.ttype, t.val))
}

func (p *Parser) nextInt() int {
	t, _ := p.nextRequired([]TokenType{T_INT})
	i, _ := strconv.Atoi(t.val)
	return i
}

func (p *Parser) nextFloat() float64 {
	t, _ := p.nextRequired([]TokenType{T_INT, T_FLOAT})
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
