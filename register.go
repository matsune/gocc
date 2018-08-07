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
		return "%AL"
	case AH:
		return "%AH"
	case AX:
		return "%AX"
	case EAX:
		return "%EAX"
	case RAX:
		return "%RAX"

	case BL:
		return "%BL"
	case BH:
		return "%BH"
	case BX:
		return "%BX"
	case EBX:
		return "%EBX"
	case RBX:
		return "%RBX"

	case CL:
		return "%CL"
	case CH:
		return "%CH"
	case CX:
		return "%CX"
	case ECX:
		return "%ECX"
	case RCX:
		return "%RCX"

	case DL:
		return "%DL"
	case DH:
		return "%DH"
	case DX:
		return "%DX"
	case EDX:
		return "%EDX"
	case RDX:
		return "%RDX"

	case SIL:
		return "%SIL"
	case SI:
		return "%SI"
	case ESI:
		return "%ESI"
	case RSI:
		return "%RSI"

	case DIL:
		return "%DIL"
	case DI:
		return "%DI"
	case EDI:
		return "%EDI"
	case RDI:
		return "%RDI"

	case R8B:
		return "%R8B"
	case R8W:
		return "%R8W"
	case R8D:
		return "%R8D"
	case R8:
		return "%R8"

	case R9B:
		return "%R9B"
	case R9W:
		return "%R9W"
	case R9D:
		return "%R9D"
	case R9:
		return "%R9"

	case RBP:
		return "%RBP"
	case RSP:
		return "%RSP"

	default:
		panic("undefined Register")
	}
}
