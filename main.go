package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	outName := flag.String("o", "asm/gocc.s", "output")
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Println("Please pass c file.")
		fmt.Println("gocc [<Option>] <filename>")
		os.Exit(1)
	}

	source, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		panic(err)
	}

	outFile, err := os.Create(*outName)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	p := NewParser(source)
	p.next()

	var c CodeGen
	c.emitMain()
	c.prologue()

	for !p.match(EOF) {
		e := p.expr()
		switch v := e.(type) {
		case BinaryExpr:
			c.s += fmt.Sprintf("\tmovl $%s, %%eax\n", v.X.(IntVal).Token.Str)
			var op string
			if v.Op.Kind == ADD {
				op = "addl"
			} else if v.Op.Kind == SUB {
				op = "subl"
			} else {
				panic("binary op")
			}
			c.s += fmt.Sprintf("\t%s $%s, %%eax\n", op, v.Y.(IntVal).Token.Str)
		case IntVal:
			c.s += fmt.Sprintf("\tmovl $%s, %%eax\n", v.Token.Str)
		}
	}

	c.epilogue()

	if _, err := outFile.WriteString(c.s); err != nil {
		panic(err)
	}
}
