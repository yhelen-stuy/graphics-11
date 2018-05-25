package main

import (
	// "bufio"
	"errors"
	"fmt"
	"io/ioutil"
	// "strconv"
	// "strings"
)

type Parser struct {
	l *Lexer
	s *Stack
}

func MakeParser() *Parser {
	return &Parser{
		s: MakeStack(),
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
	p.l = Lex(input)
	for {
		t := p.l.NextToken()
		switch t.ttype {
		case T_IDENTIFIER:
			switch tt := FindOp(t.val); tt {
			default:
				fmt.Println(tt)
			}
		}
	}
}

/*
	scanner := bufio.NewScanner(f)
	s := MakeStack()
	for scanner.Scan() {
		switch c := strings.TrimSpace(scanner.Text()); c {
		case "line":
			args := getArgs(scanner)
			if err := checkArgCount(args, 6); err != nil {
				fmt.Println(err)
				continue
			}
			fargs := numerize(args)
			e.AddEdge(fargs[0], fargs[1], fargs[2], fargs[3], fargs[4], fargs[5])
			e, _ = e.Mult(s.Peek())
			image.DrawLines(e, Color{r: 255, b: 0, g: 0})
			e = MakeMatrix(4, 0)

		case "pop":
			s.Pop()

		case "push":
			t := s.Peek()
			if t != nil {
				s.Push(t.Copy())
			}

		case "ident":
			t.Ident()

		case "scale":
			args := getArgs(scanner)
			if err := checkArgCount(args, 3); err != nil {
				fmt.Println(err)
				continue
			}
			fargs := numerize(args)
			scale := MakeScale(fargs[0], fargs[1], fargs[2])
			t = s.Pop()
			t, _ = scale.Mult(t)
			s.Push(t.Copy())

		case "move":
			args := getArgs(scanner)
			if err := checkArgCount(args, 3); err != nil {
				fmt.Println(err)
				continue
			}
			fargs := numerize(args)
			translate := MakeTranslate(fargs[0], fargs[1], fargs[2])
			t = s.Pop()
			t, _ = translate.Mult(t)
			s.Push(t.Copy())

		case "rotate":
			args := getArgs(scanner)
			if err := checkArgCount(args, 2); err != nil {
				fmt.Println(err)
				continue
			}
			// TODO: Error handling
			deg, _ := strconv.ParseFloat(args[1], 64)
			switch args[0] {
			case "x":
				t = s.Pop()
				rot := MakeRotX(deg)
				t, _ = rot.Mult(t)
				s.Push(t.Copy())
			case "y":
				t = s.Pop()
				rot := MakeRotY(deg)
				t, _ = rot.Mult(t)
				s.Push(t.Copy())
			case "z":
				t = s.Pop()
				rot := MakeRotZ(deg)
				t, _ = rot.Mult(t)
				s.Push(t.Copy())
			default:
				// TODO: Error handling
				fmt.Println("Rotate fail")
				continue
			}

		case "circle":
			args := getArgs(scanner)
			if err := checkArgCount(args, 4); err != nil {
				fmt.Println(err)
				continue
			}
			fargs := numerize(args)
			e.AddCircle(fargs[0], fargs[1], fargs[2], fargs[3])
			e, _ = e.Mult(s.Peek())
			image.DrawLines(e, Color{r: 255, b: 0, g: 0})
			e = MakeMatrix(4, 0)

		case "hermite":
			args := getArgs(scanner)
			if err := checkArgCount(args, 8); err != nil {
				fmt.Println(err)
				continue
			}
			fargs := numerize(args)
			err := e.AddHermite(fargs[0], fargs[1], fargs[2], fargs[3], fargs[4], fargs[5], fargs[6], fargs[7], 0.01)
			if err != nil {
				fmt.Println(err)
				continue
			}
			e, _ = e.Mult(s.Peek())
			image.DrawLines(e, Color{r: 255, b: 0, g: 0})
			e = MakeMatrix(4, 0)

		case "bezier":
			args := getArgs(scanner)
			if err := checkArgCount(args, 8); err != nil {
				fmt.Println(err)
				continue
			}
			fargs := numerize(args)
			err := e.AddBezier(fargs[0], fargs[1], fargs[2], fargs[3], fargs[4], fargs[5], fargs[6], fargs[7], 0.01)
			if err != nil {
				fmt.Println(err)
				continue
			}
			e, _ = e.Mult(s.Peek())
			image.DrawLines(e, Color{r: 255, b: 0, g: 0})
			e = MakeMatrix(4, 0)

		case "box":
			args := getArgs(scanner)
			if err := checkArgCount(args, 6); err != nil {
				fmt.Println(err)
				continue
			}
			fargs := numerize(args)
			p.AddBox(fargs[0], fargs[1], fargs[2], fargs[3], fargs[4], fargs[5])
			p, _ = p.Mult(s.Peek())
			image.DrawPolygons(p, Color{r: 0, b: 255, g: 0})
			p = MakeMatrix(4, 0)
			fmt.Println("drawing box")

		case "sphere":
			args := getArgs(scanner)
			if err := checkArgCount(args, 4); err != nil {
				fmt.Println(err)
				continue
			}
			fargs := numerize(args)
			p.AddSphere(fargs[0], fargs[1], fargs[2], fargs[3])
			p, _ = p.Mult(s.Peek())
			image.DrawPolygons(p, Color{r: 0, b: 255, g: 0})
			p = MakeMatrix(4, 0)

		case "torus":
			args := getArgs(scanner)
			if err := checkArgCount(args, 5); err != nil {
				fmt.Println(err)
				continue
			}
			fargs := numerize(args)
			p.AddTorus(fargs[0], fargs[1], fargs[2], fargs[3], fargs[4])
			p, _ = p.Mult(s.Peek())
			image.DrawPolygons(p, Color{r: 0, b: 255, g: 0})
			p = MakeMatrix(4, 0)

		case "clear":
			image.Clear()

		case "display":
			image.Display()

		case "save":
			args := getArgs(scanner)
			if err := checkArgCount(args, 1); err != nil {
				fmt.Println(err)
				continue
			}
			image.SavePPM(args[0])

		case "quit":
			break

		default:
			if c[0] != '#' {
				fmt.Printf("Error: Couldn't recognize command %s\n", c)
			}
			continue
		}
		// fmt.Println(t)
	}
	return nil
*/
