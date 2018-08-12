package main

import (
	"flag"
	"fmt"
	"gocc/gen"
	"gocc/parser"
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

	p := parser.NewParser(source)
	gen := gen.NewGen()

	for !p.IsEnd() {
		n := p.Parse()
		gen.Generate(n)
	}

	if _, err := outFile.WriteString(gen.Str); err != nil {
		panic(err)
	}
}
