package parser

import . "gocc/token"

var (
	typeKeys = []TokenKind{
		INT,
		VOID,
		CHAR,
		FLOAT,
		LONG,
		SHORT,
		DOUBLE,
		STRUCT,
		UNION,
		SIGNED,
		UNSIGNED,
		STATIC,
		AUTO,
		EXTERN,
		REGISTER,
		CONST,
		VOLATILE,
	}

	unaryOps = []TokenKind{
		AND,
		MUL,
		ADD,
		SUB,
		TILDE,
		NOT,
	}

	assignOps = []TokenKind{
		ASSIGN,
		MUL_ASSIGN,
		DIV_ASSIGN,
		REM_ASSIGN,
		ADD_ASSIGN,
		SUB_ASSIGN,
		LEFT_ASSIGN,
		RIGHT_ASSIGN,
		AND_ASSIGN,
		OR_ASSIGN,
		XOR_ASSIGN,
	}

	storageSpecifiers = []TokenKind{
		AUTO,
		REGISTER,
		STATIC,
		EXTERN,
		TYPEDEF,
	}

	typeSpecifiers = []TokenKind{
		VOID,
		CHAR,
		SHORT,
		INT,
		SHORT,
		LONG,
		FLOAT,
		DOUBLE,
		SIGNED,
		UNSIGNED,
	}
)
