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
	MOV
	IMUL
	IDIV
	CLTD
	XORL
	PUSH
	POP
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
	EDI
	ESI
	R8D
	R9D
	RAX
	RBP
	RSP
)

var argRegs = []Reg{EDI, ESI, EDX, ECX, R8D, R9D}

type Operand interface {
	Str() string
}

func (r Reg) Str() string     { return r.String() }
func (i IntVal) Str() string  { return "$" + string(i.Token.Str) }
func (c CharVal) Str() string { return "$" + fmt.Sprintf("%d", c.Token.Str[0]) }
func (i Ident) Str() string   { return i.Token.String() }

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
	gen.emit(PUSH, RBP)
	gen.emit(MOV, RSP, RBP)
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
	gen.emitf("\tsub \t$%d, %%rsp\n", n.Type.Size())
	gen.emitf("\tmov \t%%eax, %d(%%rbp)\n", -gen.pos)
}

func (gen *Gen) argDef(a FuncArg) {
	gen.pos += a.Type.Size()
	gen.add(a.Name.String(), gen.pos)
	gen.emitf("\tsub \t$%d, %%rsp\n", a.Type.Size())
}

func (gen *Gen) funcDef(v FuncDef) {
	gen.pos = 0
	gen.m = Map{}

	if v.Name == "main" {
		gen.global("_main")
		gen.emitFuncDef("_main")
	} else {
		gen.emitFuncDef(v.Name)
	}
	gen.prologue()

	for i, arg := range v.Args {
		gen.argDef(arg)
		if pos, ok := gen.lookup(arg.Name.String()); ok {
			if i < len(argRegs) {
				gen.emitf("\tmov \t%s, %d(%%rbp)\n", argRegs[i], -pos)
			} else {
				gen.emitf("\tmov \t%d(%%rbp), %%eax\n", (i-len(argRegs)+1)*8+8)
				gen.emitf("\tmov \t%%eax, %d(%%rbp)\n", -pos)
			}
		} else {
			panic("ident is not defined")
		}

	}

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
			gen.emitf("\tmov \t%d(%%rbp), %%eax\n", -pos)
		} else {
			panic("ident is not defined")
		}
	case IntVal:
		gen.emit(MOV, v, EAX)
	case CharVal:
		gen.emit(MOV, v, EAX)
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
	gen.emit(PUSH, RAX)

	gen.expr(e.Y)
	gen.emit(MOV, EAX, EBX)

	gen.emit(POP, RAX)

	var c Code
	if e.Op.Kind == ADD {
		c = ADDL
	} else if e.Op.Kind == SUB {
		c = SUBL
	} else if e.Op.Kind == MUL {
		c = IMUL
	} else if e.Op.Kind == DIV || e.Op.Kind == REM {
		c = IDIV
	} else {
		panic("unimplemented")
	}

	if c == IDIV {
		gen.emit(CLTD)
		gen.emit(IDIV, EBX)
		if e.Op.Kind == REM {
			gen.emit(MOV, EDX, EAX)
		}
	} else {
		gen.emit(c, EBX, EAX)
	}
}

func (gen *Gen) funcCall(e FuncCall) {
	for i := len(e.Args) - 1; i >= 0; i-- {
		gen.expr(e.Args[i])
		if i > len(argRegs)-1 {
			gen.emit(PUSH, RAX)
		} else {
			gen.emit(MOV, EAX, argRegs[i])
		}
	}
	gen.emit(CALL, e.Ident)
}

func (c Code) String() string {
	switch c {
	case ADDL:
		return "add "
	case SUBL:
		return "sub "
	case MOV:
		return "mov "
	case IMUL:
		return "imul"
	case IDIV:
		return "idiv"
	case CLTD:
		return "cltd"
	case XORL:
		return "xorl"
	case PUSH:
		return "push"
	case POP:
		return "pop "
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
	case EDI:
		return "%edi"
	case ESI:
		return "%esi"
	case R8D:
		return "%r8d"
	case R9D:
		return "%r9d"
	case RAX:
		return "%rax"
	case RBP:
		return "%rbp"
	case RSP:
		return "%rsp"
	default:
		panic("undefined reg")
	}
}
