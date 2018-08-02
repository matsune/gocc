package parser

type TokenKind int

const (
	IDENT TokenKind = iota
	// constants
	INT_CONST
	FLOAT_CONST
	STRING_CONST
	CHAR_CONST

	// type
	INT
	VOID
	CHAR
	FLOAT
	LONG
	SHORT
	DOUBLE

	// keyword
	DO
	WHILE
	IF
	ELSE
	FOR
	AUTO
	RETURN
	SWITCH
	CASE
	DEFAULT
	CONTINUE
	BREAK
	GOTO
	CONST
	EXTERN
	REGISTER
	SIGNED
	UNSIGNED
	SIZEOF
	STATIC
	STRUCT
	TYPEDEF
	UNION
	VOLATILE
	ENUM

	// operator
	ADD   // +
	SUB   // -
	MUL   // *
	DIV   // /
	REM   // %
	AND   // &
	OR    // |
	QUE   // ?
	NOT   // !
	XOR   // ^
	TILDE // ~

	ADD_ASSIGN   // +=
	SUB_ASSIGN   // -=
	MUL_ASSIGN   // *=
	DIV_ASSIGN   // /=
	REM_ASSIGN   // %=
	RIGHT_ASSIGN // >>=
	LEFT_ASSIGN  // <<=
	AND_ASSIGN   // &=
	OR_ASSIGN    // |=
	XOR_ASSIGN   // ^=

	LSHIFT // <<
	RSHIFT // >>
	ARROW  // ->
	LAND   // &&
	LOR    // ||
	INC    // ++
	DEC    // --
	EQ     // ==
	LT     // <
	GT     // >
	ASSIGN // =
	NE     // !=
	LE     // <=
	GE     // >=

	LPAREN   // (
	LBRACK   // [
	LBRACE   // {
	COMMA    // ,
	PERIOD   // .
	ELLIPSIS // ...

	RPAREN    // )
	RBRACK    // ]
	RBRACE    // }
	SEMICOLON // ;
	COLON     // :

	// special tokens
	EOF
	COMMENT // /* or //
	UNKNOWN
)

func (k TokenKind) String() string {
	return [...]string{
		IDENT: "IDENT",

		INT_CONST:    "INT_CONST",
		FLOAT_CONST:  "FLOAT_CONST",
		STRING_CONST: "STRING_CONST",
		CHAR_CONST:   "CHAR_CONST",

		INT:    "INT",
		VOID:   "VOID",
		CHAR:   "CHAR",
		FLOAT:  "FLOAT",
		LONG:   "LONG",
		SHORT:  "SHORT",
		DOUBLE: "DOUBLE",

		DO:       "DO",
		WHILE:    "WHILE",
		IF:       "IF",
		ELSE:     "ELSE",
		FOR:      "FOR",
		AUTO:     "AUTO",
		RETURN:   "RETURN",
		SWITCH:   "SWITCH",
		CASE:     "CASE",
		DEFAULT:  "DEFAULT",
		CONTINUE: "CONTINUE",
		BREAK:    "BREAK",
		GOTO:     "GOTO",
		CONST:    "CONST",
		EXTERN:   "EXTERN",
		REGISTER: "REGISTER",
		SIGNED:   "SIGNED",
		UNSIGNED: "UNSIGNED",
		SIZEOF:   "SIZEOF",
		STATIC:   "STATIC",
		STRUCT:   "STRUCT",
		TYPEDEF:  "TYPEDEF",
		UNION:    "UNION",
		VOLATILE: "VOLATILE",
		ENUM:     "ENUM",

		ADD:   "ADD",
		SUB:   "SUB",
		MUL:   "MUL",
		DIV:   "DIV",
		REM:   "REM",
		AND:   "AND",
		OR:    "OR",
		QUE:   "QUE",
		XOR:   "XOR",
		TILDE: "TILDE",

		ADD_ASSIGN:   "ADD_ASSIGN",
		SUB_ASSIGN:   "SUB_ASSIGN",
		MUL_ASSIGN:   "MUL_ASSIGN",
		DIV_ASSIGN:   "DIV_ASSIGN",
		REM_ASSIGN:   "REM_ASSIGN",
		RIGHT_ASSIGN: "RIGHT_ASSIGN",
		LEFT_ASSIGN:  "LEFT_ASSIGN",
		AND_ASSIGN:   "AND_ASSIGN",
		OR_ASSIGN:    "OR_ASSIGN",
		XOR_ASSIGN:   "XOR_ASSIGN",

		LSHIFT: "LSHIFT",
		RSHIFT: "RSHIFT",
		ARROW:  "ARROW",
		LAND:   "LAND",
		LOR:    "LOR",
		INC:    "INC",
		DEC:    "DEC",
		EQ:     "EQ",
		LT:     "LT",
		GT:     "GT",
		ASSIGN: "ASSIGN",
		NOT:    "NOT",
		NE:     "NE",
		LE:     "LE",
		GE:     "GE",

		LPAREN: "LPAREN",
		LBRACK: "LBRACK",
		LBRACE: "LBRACE",
		COMMA:  "COMMA",
		PERIOD: "PERIOD",

		RPAREN:    "RPAREN",
		RBRACK:    "RBRACK",
		RBRACE:    "RBRACE",
		SEMICOLON: "SEMICOLON",
		COLON:     "COLON",
		ELLIPSIS:  "ELLIPSIS",

		EOF:     "EOF",
		COMMENT: "COMMENT",
		UNKNOWN: "UNKNOWN",
	}[k]
}

var TypeKeys = map[string]TokenKind{
	"int":      INT,
	"void":     VOID,
	"char":     CHAR,
	"float":    FLOAT,
	"long":     LONG,
	"short":    SHORT,
	"double":   DOUBLE,
	"struct":   STRUCT,
	"union":    UNION,
	"signed":   SIGNED,
	"unsigned": UNSIGNED,
	"static":   STATIC,
	"auto":     AUTO,
	"extern":   EXTERN,
	"register": REGISTER,
	"const":    CONST,
	"volatile": VOLATILE,
}

var Keywords = map[string]TokenKind{
	"do":       DO,
	"while":    WHILE,
	"if":       IF,
	"else":     ELSE,
	"for":      FOR,
	"return":   RETURN,
	"switch":   SWITCH,
	"case":     CASE,
	"default":  DEFAULT,
	"continue": CONTINUE,
	"break":    BREAK,
	"goto":     GOTO,
	"sizeof":   SIZEOF,
	"typedef":  TYPEDEF,
}

var SingleTokens = map[byte]TokenKind{
	'?': QUE,
	'(': LPAREN,
	')': RPAREN,
	'{': LBRACE,
	'}': RBRACE,
	'[': LBRACK,
	']': RBRACK,
	';': SEMICOLON,
	',': COMMA,
	':': COLON,
	'~': TILDE,
}

type Token struct {
	Kind TokenKind
	Str  []byte
	Pos  Position
}

func NewToken() *Token {
	return &Token{Kind: UNKNOWN}
}

func (t Token) String() string {
	return string(t.Str)
}
