package main

import (
	"flag"
	"fmt"
	"gocc/parser"
	"io/ioutil"
	"os"
)

func main() {
	outName := flag.String("o", "main.s", "output")
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Println("Please pass c file.")
		fmt.Println("gocc [<Option>] *.c")
		os.Exit(1)
	}

	source, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		panic(err)
	}

	l := parser.NewLexer(source)

	outFile, err := os.Create(*outName)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	var c CodeGen
	c.emitMain()
	c.prologue()
	c.s += fmt.Sprintf("\tmovl $%s, %%eax\n", l.Next().Str)
	c.epilogue()

	if _, err := outFile.WriteString(c.s); err != nil {
		panic(err)
	}
}

type CodeGen struct {
	s string
}

func (c *CodeGen) emitMain() {
	c.s += ".global _main\n_main:\n"
}

func (c *CodeGen) prologue() {
	c.s += "\tpushq %rbp\n\tmovq %rsp, %rbp\n"
}

func (c *CodeGen) epilogue() {
	c.s += "\tpopq %rbp\nret\n"
}
