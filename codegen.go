package main

import "fmt"

type CodeGen struct {
	s string
}

func (gen *CodeGen) emitMain() {
	gen.s += ".global _main\n_main:\n"
}

func (gen *CodeGen) prologue() {
	gen.s += "\tpushq %rbp\n\tmovq %rsp, %rbp\n"
}

func (gen *CodeGen) epilogue() {
	gen.s += "\tpopq %rbp\nret\n"
}

func (gen *CodeGen) expr(e Expr) {
	switch v := e.(type) {
	case BinaryExpr:
		gen.binary(v)
	case IntVal:
		gen.intVal(v)
	}
}

func (gen *CodeGen) binary(e BinaryExpr) {
	gen.expr(e.X)

	var op string
	if e.Op.Kind == ADD {
		op = "addl"
	} else if e.Op.Kind == SUB {
		op = "subl"
	} else if e.Op.Kind == MUL {
		op = "imul"
	} else {
		panic("binaryexpr")
	}

	switch y := e.Y.(type) {
	case BinaryExpr:
		gen.binary(y)
	case IntVal:
		gen.s += fmt.Sprintf("\t%s $%s, %%eax\n", op, y.Token.Str)
	}
}

func (gen *CodeGen) intVal(e IntVal) {
	gen.s += fmt.Sprintf("\tmovl $%s, %%eax\n", e.Token.Str)
}
