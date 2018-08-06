package main

import "fmt"

// [key: var name, value: offset from ebp]
type Map map[string]int

type Gen struct {
	s   string
	pos int
	m   Map
}

func NewGen() *Gen {
	return &Gen{s: "", pos: 0, m: Map{}}
}

type Code int

const (
	ADDL Code = iota
	SUBL
	MOVL
	IMULL
	IDIVL
	CLTD
	XORL
	PUSHL
	POPL
	CALL
	LEAVE
	RET
)

type Reg int

const (
	EAX Reg = iota
	EBX
	ECX
	EDX
	EBP
	ESP
)

type Operand interface {
	Str() string
}

func (r Reg) Str() string    { return r.String() }
func (i IntVal) Str() string { return "$" + string(i.Token.Str) }
func (i Ident) Str() string  { return i.Token.String() }

// register offset of variable
func (gen *Gen) add(n string, p int) {
	gen.m[n] = p
}

func (gen *Gen) lookup(n string) (int, bool) {
	for k, v := range gen.m {
		if k == n {
			return v, true
		}
	}
	return 0, false
}

func (gen *Gen) emit(c Code, ops ...Operand) {
	gen.s += "\t" + c.String()
	for i, v := range ops {
		if i != 0 {
			gen.s += ","
		}
		gen.s += "\t" + v.Str()
	}
	gen.s += "\n"
}

func (gen *Gen) emitf(format string, a ...interface{}) {
	gen.s += fmt.Sprintf(format, a...)
}

func (gen *Gen) global(n string) {
	gen.s += ".global " + n + "\n"
}

func (gen *Gen) prologue() {
	gen.emit(PUSHL, EBP)
	gen.emit(MOVL, ESP, EBP)
}

func (gen *Gen) epilogue() {
	gen.emit(LEAVE)
	gen.emit(RET)
}

func (gen *Gen) emitFuncDef(n string) {
	gen.s += n + ":\n"
}

func (gen *Gen) generate(n Node) {
	switch v := n.(type) {
	case VarDef:
		gen.varDef(v)
	case FuncDef:
		gen.funcDef(v)
	case Expr:
		gen.expr(v)
	case Stmt:
		gen.stmt(v)
	default:
		panic("unimplemented")
	}
}

func (gen *Gen) varDef(n VarDef) {
	if n.Init != nil {
		gen.expr(*n.Init)
	}
	gen.pos += n.Type.Size()
	gen.add(n.Name, gen.pos)
	gen.emitf("\tsubl\t$%d, %%esp\n", n.Type.Size())
	gen.emitf("\tmovl\t%%eax, %d(%%ebp)\n", -gen.pos)
}

func (gen *Gen) funcDef(v FuncDef) {
	if v.Name == "main" {
		gen.global("_main")
		gen.emitFuncDef("_main")
	} else {
		gen.emitFuncDef(v.Name)
	}
	gen.prologue()

	count := -1
	for i, node := range v.Block.Nodes {
		gen.generate(node)
		count = i
	}

	if count == -1 ||
		(count > -1 && v.Block.Nodes[count].Kind() != RETURN_STMT) {
		gen.emit(XORL, EAX, EAX)
	}

	gen.epilogue()
}

func (gen *Gen) expr(e Expr) {
	switch v := e.(type) {
	case BinaryExpr:
		gen.binary(v)
	case Ident:
		if pos, ok := gen.lookup(v.Token.String()); ok {
			gen.emitf("\tmovl\t%d(%%ebp), %%eax\n", -pos)
		} else {
			panic("ident is not defined")
		}
	case IntVal:
		gen.emit(MOVL, v, EAX)
	case FuncCall:
		gen.funcCall(v)
	}
}

func (gen *Gen) stmt(e Stmt) {
	switch v := e.(type) {
	case ExprStmt:
		gen.expr(v.Expr)
	case ReturnStmt:
		gen.expr(v.Expr)
	}
}

func (gen *Gen) binary(e BinaryExpr) {
	gen.expr(e.X)
	gen.emit(PUSHL, EAX)

	gen.expr(e.Y)
	gen.emit(MOVL, EAX, EBX)

	gen.emit(POPL, EAX)

	var c Code
	if e.Op.Kind == ADD {
		c = ADDL
	} else if e.Op.Kind == SUB {
		c = SUBL
	} else if e.Op.Kind == MUL {
		c = IMULL
	} else if e.Op.Kind == DIV || e.Op.Kind == REM {
		c = IDIVL
	} else {
		panic("unimplemented")
	}

	if c == IDIVL {
		gen.emit(CLTD)
		gen.emit(IDIVL, EBX)
		if e.Op.Kind == REM {
			gen.emit(MOVL, EDX, EAX)
		}
	} else {
		gen.emit(c, EBX, EAX)
	}
}

func (gen *Gen) funcCall(e FuncCall) {
	gen.emit(CALL, e.Ident)
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
	case XORL:
		return "xorl"
	case PUSHL:
		return "pushl"
	case POPL:
		return "popl"
	case CALL:
		return "call"
	case LEAVE:
		return "leave"
	case RET:
		return "ret"
	default:
		panic("undefined code")
	}
}

func (r Reg) String() string {
	switch r {
	case EAX:
		return "%eax"
	case EBX:
		return "%ebx"
	case ECX:
		return "%ecx"
	case EDX:
		return "%edx"
	case EBP:
		return "%ebp"
	case ESP:
		return "%esp"
	default:
		panic("undefined reg")
	}
}
