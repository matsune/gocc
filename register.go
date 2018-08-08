package main

type Register int

const (
	AL Register = iota
	AH
	AX
	EAX
	RAX

	BL
	BH
	BX
	EBX
	RBX

	CL
	CH
	CX
	ECX
	RCX

	DL
	DH
	DX
	EDX
	RDX

	SIL
	SI
	ESI
	RSI

	DIL
	DI
	EDI
	RDI

	R8B
	R8W
	R8D
	R8

	R9B
	R9W
	R9D
	R9

	RBP
	RSP
)

func registerA(t CType) Register {
	switch t.Bytes() {
	case 1:
		return AL
	case 2:
		return AX
	case 4:
		return EAX
	default:
		return RAX
	}
}

func registerB(t CType) Register {
	switch t.Bytes() {
	case 1:
		return BL
	case 2:
		return BX
	case 4:
		return EBX
	default:
		return RBX
	}
}

func registerC(t CType) Register {
	switch t.Bytes() {
	case 1:
		return CL
	case 2:
		return CX
	case 4:
		return ECX
	default:
		return RCX
	}
}

func registerD(t CType) Register {
	switch t.Bytes() {
	case 1:
		return DL
	case 2:
		return DX
	case 4:
		return EDX
	default:
		return RDX
	}
}

func registerDI(t CType) Register {
	switch t.Bytes() {
	case 1:
		return DIL
	case 2:
		return DI
	case 4:
		return EDI
	default:
		return RDI
	}
}

func registerSI(t CType) Register {
	switch t.Bytes() {
	case 1:
		return SIL
	case 2:
		return SI
	case 4:
		return ESI
	default:
		return RSI
	}
}

func registerR8(t CType) Register {
	switch t.Bytes() {
	case 1:
		return R8B
	case 2:
		return R8W
	case 4:
		return R8D
	default:
		return R8
	}
}

func registerR9(t CType) Register {
	switch t.Bytes() {
	case 1:
		return R9B
	case 2:
		return R9W
	case 4:
		return R9D
	default:
		return R9
	}
}

func (r Register) String() string {
	switch r {
	case AL:
		return "%al"
	case AH:
		return "%ah"
	case AX:
		return "%ax"
	case EAX:
		return "%eax"
	case RAX:
		return "%rax"

	case BL:
		return "%bl"
	case BH:
		return "%bh"
	case BX:
		return "%bx"
	case EBX:
		return "%ebx"
	case RBX:
		return "%rbx"

	case CL:
		return "%cl"
	case CH:
		return "%ch"
	case CX:
		return "%cx"
	case ECX:
		return "%ecx"
	case RCX:
		return "%rcx"

	case DL:
		return "%dl"
	case DH:
		return "%dh"
	case DX:
		return "%dx"
	case EDX:
		return "%edx"
	case RDX:
		return "%rdx"

	case SIL:
		return "%sil"
	case SI:
		return "%si"
	case ESI:
		return "%esi"
	case RSI:
		return "%rsi"

	case DIL:
		return "%dil"
	case DI:
		return "%di"
	case EDI:
		return "%edi"
	case RDI:
		return "%rdi"

	case R8B:
		return "%r8b"
	case R8W:
		return "%r8w"
	case R8D:
		return "%r8d"
	case R8:
		return "%r8"

	case R9B:
		return "%r9b"
	case R9W:
		return "%r9w"
	case R9D:
		return "%r9d"
	case R9:
		return "%r9"

	case RBP:
		return "%rbp"
	case RSP:
		return "%rsp"

	default:
		panic("undefined Register")
	}
}
