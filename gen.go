package main

type Gen struct {
	s string
}

type Code int

const (
	ADDL Code = iota
	SUBL
	MOVL
	IMULL
	IDIVL
	CLTD
	MOVQ
	PUSHQ
	POPQ
	RET
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

type Operand interface {
	Str() string
}

func (r Reg) Str() string    { return r.String() }
func (i IntVal) Str() string { return "$" + string(i.Token.Str) }

func (gen *Gen) emit(c Code, ops ...Operand) {
	gen.s += "\t" + c.String() + "\t"
	for i, v := range ops {
		if i != 0 {
			gen.s += ", "
		}
		gen.s += v.Str()
	}
	gen.s += "\n"
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

	var c Code
	if e.Op.Kind == ADD {
		c = ADDL
	} else if e.Op.Kind == SUB {
		c = SUBL
	} else if e.Op.Kind == MUL {
		c = IMULL
	} else if e.Op.Kind == DIV {
		c = IDIVL
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

	if c == IDIVL {
		gen.emit(CLTD)
		gen.emit(IDIVL, EBX)
	} else {
		gen.emit(c, EBX, EAX)
	}
}

func (c Code) String() string {
	switch c {
	case ADDL:
		return "addl"
	case SUBL:
		return "subl"
	case MOVL:
		return "movl"
	case IMULL:
		return "imull"
	case IDIVL:
		return "idivl"
	case CLTD:
		return "cltd"
	case MOVQ:
		return "movq"
	case PUSHQ:
		return "pushq"
	case POPQ:
		return "popq"
	case RET:
		return "Ret"
	default:
		return ""
	}
}

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
