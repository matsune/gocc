package main

import "fmt"

type Gen struct {
	s string
}

func (gen *Gen) emitMain() {
	gen.s += ".global _main\n_main:\n"
}

func (gen *Gen) prologue() {
	gen.s += "\tpushq %rbp\n\tmovq %rsp, %rbp\n"
}

func (gen *Gen) epilogue() {
	gen.s += "\tpopq %rbp\nret\n"
}

func (gen *Gen) expr(e Expr) {
	fmt.Println("expr ", e)
	switch v := e.(type) {
	case BinaryExpr:
		gen.binary(v, 0)
	case IntVal:
		gen.movl(v, RAX)
	}
}

func (gen *Gen) binary(e BinaryExpr, i int) {
	switch y := e.Y.(type) {
	case BinaryExpr:
		gen.binary(y, i)
	case IntVal:
		gen.movl(y, RAX)
	}

	gen.push(RAX)

	var op string
	if e.Op.Kind == ADD {
		op = "addq"
	} else if e.Op.Kind == SUB {
		op = "subq"
	} else if e.Op.Kind == MUL {
		op = "imulq"
	} else if e.Op.Kind == DIV {
		op = "idivq"
	} else {
		panic("binaryexpr")
	}

	switch x := e.X.(type) {
	case BinaryExpr:
		gen.binary(x, i)
	case IntVal:
		gen.movl(x, RAX)
	}

	gen.pop(RBX)
	if op == "idivq" {
		gen.s += fmt.Sprintf("\tcltd\n\tidivq %s\n", string(RBX))
	} else {
		gen.emit(op, string(RBX), string(RAX))
	}
}

type Reg string

const (
	RAX Reg = "%rax"
	RBX     = "%rbx"
)

func (gen *Gen) movl(e IntVal, r Reg) {
	gen.emit("movq", "$"+string(e.Token.Str), string(r))
}

func (gen *Gen) emit(op, src, dist string) {
	gen.s += fmt.Sprintf("\t%s\t%s, %s\n", op, src, dist)
}

func (gen *Gen) push(r Reg) {
	gen.s += fmt.Sprintf("\tpush\t%s\n", string(r))
}

func (gen *Gen) pop(r Reg) {
	gen.s += fmt.Sprintf("\tpop\t%s\n", string(r))
}
