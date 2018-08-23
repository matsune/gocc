package gen

import (
	"fmt"
	"gocc/ast"
	"gocc/token"
	"reflect"
)

// [key: var name, value: offset from ebp]
type Column struct {
	pos int
	ty  ast.CType
}

type Map map[string]Column

type Gen struct {
	Str string
	pos int
	m   Map
}

func NewGen() *Gen {
	return &Gen{Str: "", pos: 0, m: Map{}}
}

var labelCount = 0

var ARG_COUNT = 6

func argsRegister(i int, t ast.CType) Register {
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

func argsRegisterPtr(i int) Register {
	switch i {
	case 0:
		return RDI
	case 1:
		return RSI
	case 2:
		return RDX
	case 3:
		return RCX
	case 4:
		return R8
	case 5:
		return R9
	default:
		panic("max i is 5")
	}
}

type Operand interface {
	Str() string
}

func (r Register) Str() string { return r.String() }

func (gen *Gen) add(n string, p int, ty ast.CType) {
	gen.m[n] = Column{p, ty}
}

func (gen *Gen) lookup(n string) (Column, bool) {
	for k, v := range gen.m {
		if k == n {
			return v, true
		}
	}
	return Column{}, false
}

func (gen *Gen) emit(c Opcode, ops ...Operand) {
	gen.Str += "\t" + c.String()
	for i, v := range ops {
		if i != 0 {
			gen.Str += ","
		}
		gen.Str += "\t" + v.Str()
	}
	gen.Str += "\n"
}

func (gen *Gen) emitf(format string, a ...interface{}) {
	gen.Str += fmt.Sprintf(format, a...)
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
	gen.Str += ".global _" + n + "\n"
	gen.Str += "_" + n + ":\n"
}

func (gen *Gen) Generate(n ast.Node) {
	switch v := n.(type) {
	case ast.VarDef:
		gen.varDef(v)
	case ast.ArrayDef:
		gen.arrayDef(v)
	case ast.FuncDef:
		gen.funcDef(v)
	case ast.Expr:
		gen.expr(v)
	case ast.Stmt:
		gen.stmt(v)
	default:
		panic("unimplemented")
	}
}

func (gen *Gen) varDef(n ast.VarDef) {
	if n.Init != nil {
		gen.expr(*n.Init)
	}
	gen.pos += n.Type.Bytes()
	gen.add(n.Token.String(), gen.pos, n.Type)
	gen.emitf("\t%s\t$%d, %s\n", SUBQ, n.Type.Bytes(), RSP)
	gen.emitf("\t%s\t%s, %d(%s)\n", mov(n.Type), registerA(n.Type), -gen.pos, RBP)
}

func (gen *Gen) arrayDef(a ast.ArrayDef) {
	if a.Subscript == nil {
		// e.g.) int a[] = {0, 1}
		s := len(a.Init.List) * a.Type.Bytes()

		gen.pos += s
		gen.add(a.Token.String(), gen.pos, a.Type)
		gen.emitf("\t%s\t$%d, %s\n", SUBQ, s, RSP)

		for idx, v := range a.Init.List {
			gen.expr(v)
			gen.emitf("\t%s\t%s, %d(%s)\n", mov(a.Type), registerA(a.Type), a.Type.Bytes()*idx-gen.pos, RBP)
		}
	} else {
		// e.g.) int a[5]
		i, ok := (*a.Subscript).(ast.IntVal)
		if !ok {
			panic("subscript is not intVal")
		}
		s := i.Num * a.Type.Bytes()

		gen.pos += s
		gen.add(a.Token.String(), gen.pos, a.Type)
		gen.emitf("\t%s\t$%d, %s\n", SUBQ, s, RSP)

		if a.Init != nil {
			for idx, v := range a.Init.List {
				gen.expr(v)
				gen.emitf("\t%s\t%s, %d(%s)\n", mov(a.Type), registerA(a.Type), a.Type.Bytes()*idx-gen.pos, RBP)
			}
		}
	}
}

func (gen *Gen) argDef(a ast.FuncArg) {
	gen.pos += a.Type.Bytes()
	gen.add(a.Name.String(), gen.pos, a.Type)
	gen.emitf("\t%s\t$%d, %s\n", SUBQ, a.Type.Bytes(), RSP)
}

func (gen *Gen) funcDef(v ast.FuncDef) {
	gen.pos = 0
	gen.m = Map{}

	gen.emitFuncDef(v.Name)
	gen.prologue()

	for i, arg := range v.Args {
		gen.argDef(arg)
		if col, ok := gen.lookup(arg.Name.String()); ok {
			if i < ARG_COUNT {
				gen.emitf("\t%s\t%s, %d(%s)\n", mov(arg.Type), argsRegister(i, arg.Type), -col.pos, RBP)
			} else {
				gen.emitf("\t%s\t%d(%s), %s\n", MOVL, (i-ARG_COUNT+1)*8+8, RBP, EAX)
				gen.emitf("\t%s\t%s, %d(%s)\n", MOVL, EAX, -col.pos, RBP)
			}
		} else {
			panic("ident is not defined")
		}
	}

	count := -1
	for i, node := range v.Block.Nodes {
		gen.Generate(node)
		count = i
	}

	if count == -1 ||
		(count > -1 && v.Block.Nodes[count].Kind() != ast.RETURN_STMT) {
		gen.emit(XORL, EAX, EAX)
	}

	gen.epilogue()
}

func (gen *Gen) expr(e ast.Expr) {
	switch v := e.(type) {
	case ast.BinaryExpr:
		gen.binary(v)
	case ast.Ident:
		if col, ok := gen.lookup(v.Token.String()); ok {
			if col.ty.Array {
				gen.emitf("\t%s \t%d(%s), %s\n", LEAQ, -col.pos, RBP, RAX)
			} else {
				gen.emitf("\t%s \t%d(%s), %s\n", mov(col.ty), -col.pos, RBP, registerA(col.ty))
			}
		} else {
			panic("ident is not defined")
		}
	case ast.IntVal:
		gen.emit(MOVL, v, EAX)
	case ast.CharVal:
		gen.emit(MOVB, v, AL)
	case ast.FuncCall:
		gen.funcCall(v)
	case ast.UnaryExpr:
		gen.unaryExpr(v)
	case ast.PtrVal:
		gen.pointerVal(v)
	case ast.AddressVal:
		gen.addressVal(v)
	case ast.AssignExpr:
		gen.assignExpr(v)
	case ast.SubscriptExpr:
		gen.subscriptExpr(v)
	case ast.IncExpr:
		if col, ok := gen.lookup(v.Ident.Token.String()); ok {
			if col.ty.Ptr {
				gen.emitf("\t%s\t$%d, %d(%s)\n", ADDQ, col.ty.Primitive.Bytes(), -col.pos, RBP)
			} else {
				gen.emitf("\t%s\t%s, %d(%s)\n", ADDL, "$1", -col.pos, RBP)
			}
		} else {
			panic("ident is not defined")
		}
	case ast.DecExpr:
		if col, ok := gen.lookup(v.Ident.Token.String()); ok {
			if col.ty.Ptr {
				gen.emitf("\t%s\t%s, %d(%s)\n", SUBQ, "$8", -col.pos, RBP)
			} else {
				gen.emitf("\t%s\t%s, %d(%s)\n", SUBL, "$1", -col.pos, RBP)
			}
		} else {
			panic("ident is not defined")
		}
	default:
		panic(fmt.Sprintf("unimplemented expr type: %s", reflect.TypeOf(e)))
	}
}

func (gen *Gen) stmt(e ast.Stmt) {
	switch v := e.(type) {
	case ast.ExprStmt:
		gen.expr(v.Expr)
	case ast.ReturnStmt:
		gen.expr(v.Expr)
	case ast.IfStmt:
		gen.ifStmt(v)
	case ast.ForStmt:
		gen.forStmt(v)
	}
}

func (gen *Gen) blockStmt(b ast.BlockStmt) {
	for _, n := range b.Nodes {
		gen.Generate(n)
	}
}

func (gen *Gen) ifStmt(v ast.IfStmt) {
	if e := v.Expr; e != nil { // if (...) { ... }
		switch e := (*v.Expr).(type) {
		case ast.BinaryExpr:
			gen.binary(e)

			gen.invertJump(e.Op.Kind, labelCount)

			gen.blockStmt(v.Block)
			if v.Else != nil {
				gen.emitf("\t%s\t.L%d\n", JMP, labelCount+1)
			}

			gen.emitf(".L%d:\n", labelCount)
			labelCount++

			if el := v.Else; el != nil {
				gen.ifStmt(*el)
			}
		case ast.IntVal:
			panic("unimplemented if ast.IntVal")
		default:
			panic("unimplemented if expr type")
		}
	} else { // else { ... }
		gen.blockStmt(v.Block)
		gen.emitf(".L%d:\n", labelCount)
	}
}

func (gen *Gen) forStmt(v ast.ForStmt) {
	if v.E1 != nil {
		gen.Generate(v.E1)
	}
	gen.emitf("\t%s\t.L%d\n", JMP, labelCount)
	gen.emitf(".L%d:\n", labelCount+1)
	gen.blockStmt(v.Block)
	if v.E3 != nil {
		gen.expr(*v.E3)
	}
	gen.emitf(".L%d:\n", labelCount)
	if v.E2 != nil {
		gen.expr(*v.E2)
		if b, ok := (*v.E2).(ast.BinaryExpr); ok {
			gen.jump(b.Op.Kind, labelCount+1)
		}
	}
	labelCount += 2
}

func (gen *Gen) jump(kind token.TokenKind, label int) {
	var op Opcode
	switch kind {
	case token.NE:
		op = JNE
	case token.EQ:
		op = JE
	case token.LT:
		op = JL
	case token.LE:
		op = JLE
	case token.GT:
		op = JG
	case token.GE:
		op = JGE
	default:
		panic(fmt.Sprintf("unimplemented jump token %s", kind))
	}
	gen.emitf("\t%s\t.L%d\n", op, label)
}

func (gen *Gen) invertJump(kind token.TokenKind, label int) {
	var op Opcode
	switch kind {
	case token.NE:
		op = JE
	case token.EQ:
		op = JNE
	case token.LT:
		op = JLE
	case token.LE:
		op = JL
	case token.GT:
		op = JGE
	case token.GE:
		op = JG
	default:
		panic(fmt.Sprintf("unimplemented jump token %s", kind))
	}
	gen.emitf("\t%s\t.L%d\n", op, label)
}

func (gen *Gen) binary(e ast.BinaryExpr) {
	gen.expr(e.X)
	gen.emit(PUSH, RAX)

	gen.expr(e.Y)
	gen.emit(MOVL, EAX, EBX)

	gen.emit(POP, RAX)

	switch e.Op.Kind {
	case token.ADD:
		gen.emit(ADDL, EBX, EAX)
	case token.SUB:
		gen.emit(SUBL, EBX, EAX)
	case token.MUL:
		gen.emit(IMUL, EBX, EAX)
	case token.DIV, token.REM:
		gen.emit(CLTD)
		gen.emit(IDIV, EBX)
		if e.Op.Kind == token.REM {
			gen.emit(MOVL, EDX, EAX)
		}
	case token.EQ, token.NE, token.LT, token.LE, token.GT, token.GE:
		gen.emit(CMPL, EBX, EAX)
	default:
		panic("unimplemented binary op")
	}
}

func (gen *Gen) funcCall(e ast.FuncCall) {
	for i := len(e.Args) - 1; i >= 0; i-- {
		gen.expr(e.Args[i])
		if i > ARG_COUNT-1 {
			gen.emit(PUSH, RAX)
		} else {
			gen.emit(MOVQ, RAX, argsRegisterPtr(i))
		}
	}
	gen.emit(CALL, e.Ident)
}

func (gen *Gen) unaryExpr(e ast.UnaryExpr) {
	panic("gen.unaryExpr")
}

func (gen *Gen) assignExpr(e ast.AssignExpr) {
	gen.expr(e.R)

	switch v := e.L.(type) {
	case ast.PtrVal:
		if col, ok := gen.lookup(v.Token.String()); ok {
			gen.emitf("\t%s\t%d(%s), %s\n", mov(col.ty), -col.pos, RBP, registerB(col.ty))
		} else {
			panic("ident is not defined")
		}
		gen.emitf("\t%s\t%s, (%s)\n", MOVQ, RAX, RBX)
	case ast.SubscriptExpr:
		if col, ok := gen.lookup(v.Token.String()); ok {
			i, ok := v.Expr.(ast.IntVal)
			if !ok {
				panic("subscript expr is not intVal")
			}
			gen.emitf("\t%s\t%s, %d(%s)\n", mov(col.ty), registerA(col.ty), (i.Num*col.ty.Bytes() - col.pos), RBP)
		} else {
			panic("ident is not defined")
		}
	case ast.Ident:
		if col, ok := gen.lookup(v.Token.String()); ok {
			gen.emitf("\t%s\t%s, %d(%s)\n", mov(col.ty), registerA(col.ty), -col.pos, RBP)
		} else {
			panic("ident is not defined")
		}
	default:
		panic("unimplemented assignExpr L type")
	}
}

func (gen *Gen) subscriptExpr(e ast.SubscriptExpr) {
	if col, ok := gen.lookup(e.Token.String()); ok {
		i, ok := e.Expr.(ast.IntVal)
		if !ok {
			panic("subscript should be intVal")
		}

		gen.emitf("\t%s\t%d(%s), %s\n", mov(col.ty), (i.Num*col.ty.Bytes() - col.pos), RBP, registerA(col.ty))
	} else {
		panic("ident is not defined")
	}
}

func (gen *Gen) pointerVal(e ast.PtrVal) {
	if col, ok := gen.lookup(e.Token.String()); ok {
		gen.emitf("\t%s\t%d(%s), %s\n", mov(col.ty), -col.pos, RBP, registerB(col.ty))
		gen.emitf("\t%s\t(%s), %s\n", mov(col.ty), registerB(col.ty), registerA(col.ty))
	} else {
		panic("ident is not defined")
	}
}

func (gen *Gen) addressVal(e ast.AddressVal) {
	if col, ok := gen.lookup(e.Token.String()); ok {
		gen.emitf("\t%s\t%d(%s), %s\n", LEAQ, -col.pos, RBP, RAX)
	} else {
		panic("ident is not defined")
	}
}
