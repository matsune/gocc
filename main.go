package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	out := flag.String("o", "main.s", "output")
	flag.Parse()

	if len(flag.Args()) != 1 {
		panic("pass 1 value")
	}

	f, err := os.Create(*out)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Fprintf(f, `.global _main
_main:
	pushq %%rbp
	movq	%%rsp, %%rbp
	movl	$%s,	%%eax
	popq	%%rbp
	ret
`, flag.Arg(0))
}
