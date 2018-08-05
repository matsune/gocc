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

	gen := NewGen()
	gen.emitMain()
	gen.prologue()

	for !p.isEnd() {
		n := p.parse()

		switch v := n.(type) {
		case Expr:
			gen.expr(v)
		case VarDef:
			gen.varDef(v)
		default:
			panic("unimplemented")
		}
	}

	gen.epilogue()

	if _, err := outFile.WriteString(gen.s); err != nil {
		panic(err)
	}
}
