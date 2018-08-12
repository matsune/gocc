package main

import (
	"flag"
	"fmt"
	"gocc/gen"
	"gocc/parser"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	o := flag.String("o", "", "outfile")
	s := flag.Bool("S", false, "output assembler file")
	c := flag.Bool("c", false, "generate object file")
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Println("gocc [<Option>] <filename>")
		os.Exit(1)
	}

	cFile := flag.Arg(0)
	sName := "tmp_gocc.s"

	if len(*o) > 0 && *s {
		sName = *o
	}

	if len(*o) < 1 {
		if ok := strings.HasSuffix(cFile, ".c"); !ok {
			fmt.Println("file does not have suffix .c")
			os.Exit(1)
		}
		_, name := filepath.Split(cFile)

		if *s {
			sName = strings.TrimSuffix(name, ".c") + ".s"
		} else {
			if *c {
				*o = strings.TrimSuffix(name, ".c") + ".o"
			} else {
				*o = "a.out"
			}
		}
	}

	source, err := ioutil.ReadFile(cFile)
	if err != nil {
		panic(err)
	}

	sFile, err := os.Create(sName)
	if err != nil {
		panic(err)
	}
	defer sFile.Close()

	p := parser.NewParser(source)
	gen := gen.NewGen()

	for !p.IsEnd() {
		n := p.Parse()
		gen.Generate(n)
	}

	if _, err := sFile.WriteString(gen.Str); err != nil {
		panic(err)
	}

	if !*s {
		if *c {
			err = exec.Command("as", "-o", *o, sName).Run()
		} else {
			err = exec.Command("gcc", "-o", *o, sName).Run()
		}
		if err != nil {
			panic(err)
		}
		err = exec.Command("rm", sName).Run()
		if err != nil {
			panic(err)
		}
	}
}
