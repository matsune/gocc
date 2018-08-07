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

var ARG_COUNT = 6

func argsRegister(i int, t CType) Register {
	switch i {
	case 0:
		return registerDI(t)
	case 1:
		return registerSI(t)
	case 2:
		return registerD(t)
	case 3:
		return registerC(t)
	case 4:
		return registerR8(t)
	case 5:
		return registerR9(t)
	default:
		panic("max i is 5")
	}
}

type Operand interface {
	Str() string
}

func (r Register) Str() string { return r.String() }
func (i IntVal) Str() string   { return "$" + string(i.Token.Str) }
func (c CharVal) Str() string  { return "$" + fmt.Sprintf("%d", c.Token.Str[0]) }
func (i Ident) Str() string    { return i.Token.String() }

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

func (gen *Gen) emit(c Opcode, ops ...Operand) {
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
	gen.emit(MOVQ, RSP, RBP)
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
	gen.pos += n.Type.Bytes()
	gen.add(n.Name, gen.pos)
	gen.emitf("\t%s\t$%d, %s\n", SUBQ, n.Type.Bytes(), RSP)
	gen.emitf("\t%s\t%s, %d(%s)\n", MOVL, EAX, -gen.pos, RBP)
}

func (gen *Gen) argDef(a FuncArg) {
	gen.pos += a.Type.Bytes()
	gen.add(a.Name.String(), gen.pos)
	gen.emitf("\t%s\t$%d, %s\n", SUBQ, a.Type.Bytes(), RSP)
}

func (gen *Gen) funcDef(v FuncDef) {
	gen.pos = 0
	gen.m = Map{}

	if v.Name == "main" {
		gen.global("_main")
		gen.emitFuncDef("_main")
	} else {
		gen.global(v.Name)
		gen.emitFuncDef(v.Name)
	}
	gen.prologue()

	for i, arg := range v.Args {
		gen.argDef(arg)
		if pos, ok := gen.lookup(arg.Name.String()); ok {
			if i < ARG_COUNT {
				gen.emitf("\t%s\t%s, %d(%s)\n", mov(arg.Type), argsRegister(i, arg.Type), -pos, RBP)
			} else {
				gen.emitf("\t%s\t%d(%s), %s\n", MOVL, (i-ARG_COUNT+1)*8+8, RBP, EAX)
				gen.emitf("\t%s\t%s, %d(%s)\n", MOVL, EAX, -pos, RBP)
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
		(count > -1 && v.Block.Nodes[count].Kind() != AST_ReturnStmt) {
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
			gen.emitf("\t%s \t%d(%s), %s\n", MOVL, -pos, RBP, EAX)
		} else {
			panic("ident is not defined")
		}
	case IntVal:
		gen.emit(mov(C_int), v, registerA(C_int))
	case CharVal:
		gen.emit(mov(C_char), v, registerA(C_char))
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
	gen.emit(MOVL, EAX, EBX)

	gen.emit(POP, RAX)

	var c Opcode
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
			gen.emit(MOVL, EDX, EAX)
		}
	} else {
		gen.emit(c, EBX, EAX)
	}
}

func (gen *Gen) funcCall(e FuncCall) {
	for i := len(e.Args) - 1; i >= 0; i-- {
		gen.expr(e.Args[i])
		if i > ARG_COUNT-1 {
			gen.emit(PUSH, RAX)
		} else {
			// FIXME
			gen.emit(MOVL, EAX, argsRegister(i, C_int))
		}
	}
	gen.emit(CALL, e.Ident)
}
