package main

import "fmt"

type Gen struct {
	s string
}

type Inst string

const (
	ADDL  Inst = "addl"
	SUBL       = "subl"
	MOVL       = "movl"
	IMULL      = "imull"
	IDIVL      = "idivl"
	CLTD       = "cltd"
	MOVQ       = "movq"
	PUSHQ      = "pushq"
	POPQ       = "popq"
	RET        = "Ret"
)

type Reg int

const (
	RAX Reg = iota
	RBX
	EAX
	EBX
	RBP
	RSP
)

func (r Reg) String() string {
	switch r {
	case RAX:
		return "%rax"
	case RBX:
		return "%rbx"
	case EAX:
		return "%eax"
	case EBX:
		return "%ebx"
	case RBP:
		return "%rbp"
	case RSP:
		return "%rsp"
	default:
		return ""
	}
}

func (gen *Gen) emitMain() {
	gen.s += ".global _main\n_main:\n"
}

func (gen *Gen) prologue() {
	gen.emit(PUSHQ, RBP)
	gen.emit(MOVQ, RSP, RBP)
}

func (gen *Gen) epilogue() {
	gen.emit(POPQ, RBP)
	gen.emit(RET)
}

func (gen *Gen) expr(e Expr) {
	fmt.Println("expr ", e)
	switch v := e.(type) {
	case BinaryExpr:
		gen.binary(v, 0)
	case IntVal:
		gen.emit(MOVL, v, EAX)
	}
}

func (gen *Gen) binary(e BinaryExpr, i int) {
	switch y := e.Y.(type) {
	case BinaryExpr:
		gen.binary(y, i)
	case IntVal:
		gen.emit(MOVL, y, EAX)
	}

	gen.emit(PUSHQ, RAX)

	var op Inst
	if e.Op.Kind == ADD {
		op = ADDL
	} else if e.Op.Kind == SUB {
		op = SUBL
	} else if e.Op.Kind == MUL {
		op = IMULL
	} else if e.Op.Kind == DIV {
		op = IDIVL
	} else {
		panic("unimplemented")
	}

	switch x := e.X.(type) {
	case BinaryExpr:
		gen.binary(x, i)
	case IntVal:
		gen.emit(MOVL, x, EAX)
	}

	gen.emit(POPQ, RBX)

	if op == IDIVL {
		gen.emit(CLTD)
		gen.emit(IDIVL, EBX)
	} else {
		gen.emit(op, EBX, EAX)
	}
}

type Src interface {
	Str() string
}

func (r Reg) Str() string    { return r.String() }
func (i IntVal) Str() string { return "$" + string(i.Token.Str) }

func (gen *Gen) emit(op Inst, src ...Src) {
	gen.s += "\t" + string(op) + "\t"
	for i, v := range src {
		if i != 0 {
			gen.s += ", "
		}
		gen.s += v.Str()
	}
	gen.s += "\n"
}
