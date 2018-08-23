package ast

import (
	"fmt"
	"gocc/token"
)

type Kind int

const (
	VAR_DEF Kind = iota
	ARRAY_DEF
	FUNC_DEF
	FUNC_ARG
	IDENT
	// expr
	BINARY_EXPR
	COND_EXPR
	UNARY_EXPR
	ASSIGN_EXPR
	SUBSCRIPT_EXPR
	INC_EXPR
	DEC_EXPR
	FUNC_CALL
	INT_VAL
	CHAR_VAL
	PTR_VAL
	ADDRESS_VAL
	ARRAY_INIT
	// stmt
	BLOCK_STMT
	RETURN_STMT
	EXPR_STMT
	IF_STMT
	FOR_STMT
)

type PrimitiveType int

const (
	C_int PrimitiveType = iota
	C_void
	C_char
	C_float
	C_long
	C_short
	C_double
)

func (p PrimitiveType) Bytes() int {
	switch p {
	case C_int:
		return 4
	case C_char:
		return 1
	default:
		panic("unimplemented type size")
	}
}

func (t CType) Bytes() int {
	if t.Ptr {
		return 8
	}
	return t.Primitive.Bytes()
}

type CType struct {
	Primitive PrimitiveType
	Ptr       bool
	Array     bool
}

func (p PrimitiveType) String() string {
	switch p {
	case C_int:
		return "int"
	case C_void:
		return "void"
	case C_char:
		return "char"
	case C_float:
		return "float"
	case C_long:
		return "long"
	case C_short:
		return "short"
	case C_double:
		return "double"
	default:
		panic("undefined Type")
	}
}

func (t CType) String() string {
	s := t.Primitive.String()
	if t.Ptr {
		s += " *"
	}
	return s
}

type (
	Node interface {
		Kind() Kind
	}
)

type (
	Ident struct {
		Token *token.Token
	}

	VarDef struct {
		Type  CType
		Token *token.Token
		Init  *Expr
	}

	ArrayDef struct {
		Type      CType
		Token     *token.Token
		Subscript *Expr
		Init      *ArrayInit
	}

	FuncDef struct {
		Type  CType
		Name  string
		Args  []FuncArg
		Block BlockStmt
	}

	FuncArg struct {
		Type CType
		Name *token.Token
	}
)

type (
	Expr interface {
		Node
		expr()
	}

	BinaryExpr struct {
		X  Expr
		Op *token.Token
		Y  Expr
	}

	CondExpr struct {
		Cond Expr
		L    Expr
		R    Expr
	}

	UnaryExpr struct {
		Op   *token.Token
		Expr Expr
	}

	AssignExpr struct {
		L  Expr
		Op *token.Token
		R  Expr
	}

	// a[0], b[10]
	SubscriptExpr struct {
		Token *token.Token
		Expr  Expr
	}

	IncExpr struct {
		Ident Ident
	}

	DecExpr struct {
		Ident Ident
	}

	IntVal struct {
		Num int
	}

	CharVal struct {
		Token *token.Token
	}

	FuncCall struct {
		Ident Ident
		Args  []Expr
	}

	PtrVal struct {
		Token *token.Token
	}
	AddressVal struct {
		Token *token.Token
	}

	ArrayInit struct {
		List []Expr
	}
)

type (
	Stmt interface {
		Node
		stmt()
	}

	BlockStmt struct {
		Nodes []Node
	}

	ReturnStmt struct {
		Expr Expr
	}

	ExprStmt struct {
		Expr Expr
	}

	IfStmt struct {
		Expr  *Expr
		Block BlockStmt
		Else  *IfStmt
	}

	ForStmt struct {
		E1    Node
		E2    *Expr
		E3    *Expr
		Block BlockStmt
	}
)

func (VarDef) Kind() Kind        { return VAR_DEF }
func (ArrayDef) Kind() Kind      { return ARRAY_DEF }
func (FuncDef) Kind() Kind       { return FUNC_DEF }
func (FuncArg) Kind() Kind       { return FUNC_ARG }
func (Ident) Kind() Kind         { return IDENT }
func (BinaryExpr) Kind() Kind    { return BINARY_EXPR }
func (CondExpr) Kind() Kind      { return COND_EXPR }
func (UnaryExpr) Kind() Kind     { return UNARY_EXPR }
func (AssignExpr) Kind() Kind    { return ASSIGN_EXPR }
func (SubscriptExpr) Kind() Kind { return SUBSCRIPT_EXPR }
func (IncExpr) Kind() Kind       { return INC_EXPR }
func (DecExpr) Kind() Kind       { return DEC_EXPR }
func (FuncCall) Kind() Kind      { return FUNC_CALL }
func (IntVal) Kind() Kind        { return INT_VAL }
func (CharVal) Kind() Kind       { return CHAR_VAL }
func (PtrVal) Kind() Kind        { return PTR_VAL }
func (AddressVal) Kind() Kind    { return ADDRESS_VAL }
func (ArrayInit) Kind() Kind     { return ARRAY_INIT }
func (BlockStmt) Kind() Kind     { return BLOCK_STMT }
func (ReturnStmt) Kind() Kind    { return RETURN_STMT }
func (ExprStmt) Kind() Kind      { return EXPR_STMT }
func (IfStmt) Kind() Kind        { return IF_STMT }
func (ForStmt) Kind() Kind       { return FOR_STMT }

func (Ident) expr()         {}
func (BinaryExpr) expr()    {}
func (CondExpr) expr()      {}
func (UnaryExpr) expr()     {}
func (AssignExpr) expr()    {}
func (SubscriptExpr) expr() {}
func (IncExpr) expr()       {}
func (DecExpr) expr()       {}
func (FuncCall) expr()      {}
func (IntVal) expr()        {}
func (CharVal) expr()       {}
func (PtrVal) expr()        {}
func (AddressVal) expr()    {}
func (ArrayInit) expr()     {}

func (BlockStmt) stmt()  {}
func (ReturnStmt) stmt() {}
func (ExprStmt) stmt()   {}
func (IfStmt) stmt()     {}
func (ForStmt) stmt()    {}

func (i IntVal) Str() string  { return fmt.Sprintf("$%d", i.Num) }
func (c CharVal) Str() string { return "$" + fmt.Sprintf("%d", c.Token.Str[0]) }
func (i Ident) Str() string   { return "_" + i.Token.String() }
