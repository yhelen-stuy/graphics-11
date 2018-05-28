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
	lexer   *Lexer
	stack   *Stack
	poly    *Matrix
	edge    *Matrix
	trans   *Matrix
	image   *Image
	history []Token
}

func MakeParser() *Parser {
	trans := MakeMatrix(4, 4)
	trans.Ident()
	return &Parser{
		stack:   MakeStack(),
		poly:    MakeMatrix(4, 0),
		edge:    MakeMatrix(4, 0),
		trans:   trans,
		image:   MakeImage(500, 500),
		history: make([]Token, 1),
	}
}

func ParseFile(filename string) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.New("Couldn't read file")
	}
	str := string(buf)
	p := MakeParser()
	p.parseString(str)
	return nil
}

func (p *Parser) parseString(input string) error {
	p.lexer = Lex(input)
	for {
		t := p.next()
		switch t.ttype {
		case T_EOF:
			fmt.Println("received eof")
			return nil
		case T_IDENTIFIER:
			switch tt := FindOp(t.val); tt {
			case PUSH:
				trans := p.stack.Peek()
				if trans != nil {
					p.stack.Push(trans.Copy())
				}
			case POP:
				p.stack.Pop()
			case SCALE:
				scale := MakeScale(p.nextFloat(), p.nextFloat(), p.nextFloat())
				p.trans = p.stack.Pop()
				p.trans, _ = scale.Mult(p.trans)
				p.stack.Push(p.trans.Copy())
			case MOVE:
				translate := MakeTranslate(p.nextFloat(), p.nextFloat(), p.nextFloat())
				p.trans = p.stack.Pop()
				p.trans, _ = translate.Mult(p.trans)
				p.stack.Push(p.trans.Copy())
			case ROTATE:
				axis, err := p.nextRequired([]TokenType{T_STRING})
				if err != nil {
					panic(err)
				}
				switch axis.val {
				case "x":
					p.trans = p.stack.Pop()
					rot := MakeRotX(p.nextFloat())
					p.trans, _ = rot.Mult(p.trans)
					p.stack.Push(p.trans.Copy())
				case "y":
					p.trans = p.stack.Pop()
					rot := MakeRotY(p.nextFloat())
					p.trans, _ = rot.Mult(p.trans)
					p.stack.Push(p.trans.Copy())
				case "z":
					p.trans = p.stack.Pop()
					rot := MakeRotZ(p.nextFloat())
					p.trans, _ = rot.Mult(p.trans)
					p.stack.Push(p.trans.Copy())
				default:
					// TODO: Error handling
					fmt.Println("Rotate fail")
					continue
				}
			case BOX:
				p.poly.AddBox(p.nextFloat(), p.nextFloat(), p.nextFloat(), p.nextFloat(), p.nextFloat(), p.nextFloat())
				p.poly, _ = p.poly.Mult(p.stack.Peek())
				p.image.DrawPolygons(p.poly, Color{r: 0, b: 255, g: 0})
				p.poly = MakeMatrix(4, 0)
			case SPHERE:
				p.poly.AddSphere(p.nextFloat(), p.nextFloat(), p.nextFloat(), p.nextFloat())
				p.poly, _ = p.poly.Mult(p.stack.Peek())
				p.image.DrawPolygons(p.poly, Color{r: 0, b: 255, g: 0})
				p.poly = MakeMatrix(4, 0)
			case TORUS:
				p.poly.AddTorus(p.nextFloat(), p.nextFloat(), p.nextFloat(), p.nextFloat(), p.nextFloat())
				p.poly, _ = p.poly.Mult(p.stack.Peek())
				p.image.DrawPolygons(p.poly, Color{r: 0, b: 255, g: 0})
				p.poly = MakeMatrix(4, 0)
			case LINE:
				p.edge.AddEdge(p.nextFloat(), p.nextFloat(), p.nextFloat(), p.nextFloat(), p.nextFloat(), p.nextFloat())
				p.edge, _ = p.edge.Mult(p.stack.Peek())
				p.image.DrawLines(p.edge, Color{r: 255, b: 0, g: 0})
				p.edge = MakeMatrix(4, 0)
			case SAVE:
				f, err := p.nextRequired([]TokenType{T_STRING})
				if err != nil {
					panic(err)
				}
				p.image.SavePPM("temp")
				p.image.ConvertPNG("temp", f.val)
			case DISPLAY:
				p.image.Display()
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
