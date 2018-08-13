package lexer

import . "gocc/token"

var (
	typeKeys = map[string]TokenKind{
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

	keywords = map[string]TokenKind{
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

	singleTokens = map[byte]TokenKind{
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
)
