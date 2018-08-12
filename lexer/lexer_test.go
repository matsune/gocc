package lexer

import (
	. "gocc/token"
	"testing"
)

var lexerTests = []struct {
	source string
	kinds  []TokenKind
}{
	{
		`a "aaa111" 12 'c'`,
		[]TokenKind{
			IDENT, STRING_CONST, INT_CONST, CHAR_CONST, EOF, // - TODO: FLOAT_CONST
		},
	},
	{
		`a int void char float long short do while if else for auto return switch case default continue break goto const extern register signed unsigned sizeof static struct typedef union volatile`,
		[]TokenKind{
			IDENT, INT, VOID, CHAR, FLOAT, LONG, SHORT,
			DO, WHILE, IF, ELSE, FOR, AUTO, RETURN, SWITCH, CASE, DEFAULT, CONTINUE, BREAK, GOTO,
			CONST, EXTERN, REGISTER, SIGNED, UNSIGNED, SIZEOF, STATIC, STRUCT, TYPEDEF, UNION, VOLATILE, EOF,
		},
	},
	{
		`... >>= <<= += -= *= /= %=
		 &= ^= |= >> << ++ -- -> && || <= >= == != ;
		 { } , : = ( ) [ ] . & ! ~ - + * / % < > ^ | ?`,
		[]TokenKind{
			ELLIPSIS, RIGHT_ASSIGN, LEFT_ASSIGN, ADD_ASSIGN, SUB_ASSIGN, MUL_ASSIGN, DIV_ASSIGN, REM_ASSIGN,
			AND_ASSIGN, XOR_ASSIGN, OR_ASSIGN, RSHIFT, LSHIFT, INC, DEC, ARROW, LAND, LOR, LE, GE, EQ, NE, SEMICOLON,
			LBRACE, RBRACE, COMMA, COLON, ASSIGN, LPAREN, RPAREN, LBRACK, RBRACK, PERIOD, AND, NOT, TILDE, SUB, ADD,
			MUL, DIV, REM, LT, GT, XOR, OR, QUE, EOF,
		},
	},
	{
		`// comment
		/*
		comment
		*/`,
		[]TokenKind{
			EOF,
		},
	},
	{
		`printf("%d\n",10)`,
		[]TokenKind{
			IDENT, LPAREN, STRING_CONST, COMMA, INT_CONST, RPAREN, EOF,
		},
	},
}

func TestLexer(t *testing.T) {
	for n, tt := range lexerTests {
		t.Logf("[test %d]", n+1)

		l := NewLexer([]byte(tt.source))
		token := l.Next()
		i := 0

		for token.Kind != EOF {
			if i > len(tt.kinds)-1 {
				t.Errorf("[test %d] i = %d, but count of tt.kinds is %d", n+1, i, len(tt.kinds))
				return
			}

			if token.Kind != tt.kinds[i] {
				t.Errorf("[test %d] expected kind is "+tt.kinds[i].String()+", but got "+token.Kind.String(), n+1)
			} else {
				t.Logf("succeed %s", token.Str)
			}

			i++
			token = l.Next()
		}

		if i != len(tt.kinds)-1 {
			t.Errorf("[test %d] lexer got %d tokens, expected %d tokens", n+1, i+1, len(tt.kinds))
			return
		}
	}
}
