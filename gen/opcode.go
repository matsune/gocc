package gen

import "gocc/ast"

type Opcode int

const (
	MOVB Opcode = iota
	MOVW
	MOVL
	MOVQ
	ADDL
	SUBL
	SUBQ
	IMUL
	IDIV
	CLTD
	XORL
	PUSH
	POP
	LEAQ
	CALL
	LEAVE
	RET
)

func mov(t ast.CType) Opcode {
	switch t.Bytes() {
	case 1:
		return MOVB
	case 2:
		return MOVW
	case 4:
		return MOVL
	default:
		return MOVQ
	}
}

func (c Opcode) String() string {
	switch c {
	case MOVB:
		return "movb"
	case MOVW:
		return "movw"
	case MOVL:
		return "movl"
	case MOVQ:
		return "movq"
	case ADDL:
		return "addl"
	case SUBL:
		return "subl"
	case SUBQ:
		return "subq"
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
	case LEAQ:
		return "leaq"
	case LEAVE:
		return "leave"
	case RET:
		return "ret"
	default:
		panic("undefined code")
	}
}