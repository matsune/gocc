package main

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
