package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	buf, _ := ioutil.ReadFile("test.mdl")
	s := string(buf)
	l, _ := Lex(s)
	for {
		s := l.NextToken()
		if s.ttype != 0 {
			fmt.Println(s)
		}
	}
	// time.Sleep(time.Second * 10)
	// image := MakeImage(500, 500)
	// t := MakeMatrix(4, 4)
	// t.Ident()
	// e := MakeMatrix(4, 0)
	// p := MakeMatrix(4, 0)
	// ParseFile("galleryscript", t, p, e, image)
	// ParseFile("script", t, p, e, image)
}
