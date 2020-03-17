package main

import (
	"bufio"
	"unicode"
	"unicode/utf8"
)

////////////////////////////////////////////////////////////////

type CharClass int
type Token int

const (
	Eof = -1

	letter  CharClass = 0
	digit   CharClass = 1
	unknown CharClass = 99

	Ident         Token = 23
	IntLit        Token = 24
	LeftParen     Token = 25
	RightParen    Token = 26
	Comma         Token = 27
	EqualSign     Token = 28
	MinusKeyword  Token = 29
	IszeroKeyword Token = 30
	IfKeyword     Token = 31
	ThenKeyword   Token = 32
	ElseKeyword   Token = 33
	LetKeyword    Token = 34
	InKeyword     Token = 35
)

type CharType struct {
	char  string
	class CharClass
}

type LexemeData struct {
	token  Token
	lexeme string
}

////////////////////////////////////////////////////////////////

func lookup(char string) Token {
	if char == "(" {
		return LeftParen
	}
	if char == ")" {
		return RightParen
	}
	if char == "," {
		return Comma
	}
	if char == "=" {
		return EqualSign
	}

	return Eof
}

func lookupKeyword(s string) Token {
	if s == "minus" {
		return MinusKeyword
	}
	if s == "iszero" {
		return IszeroKeyword
	}
	if s == "if" {
		return IfKeyword
	}
	if s == "then" {
		return ThenKeyword
	}
	if s == "else" {
		return ElseKeyword
	}
	if s == "let" {
		return LetKeyword
	}
	if s == "in" {
		return InKeyword
	}

	return Ident
}

func getCharType(char string) CharType {
	c, _ := utf8.DecodeRuneInString(char)

	if char == "" {
		return CharType{char, Eof}
	}
	if unicode.IsDigit(c) {
		return CharType{char, digit}
	}
	if unicode.IsLetter(c) {
		return CharType{char, letter}
	}

	return CharType{char, unknown}
}

func getNextNonBlankChar(s *bufio.Scanner) string {
	data := s.Scan()
	if !data {
		return ""
	}

	t := s.Text()
	r, _ := utf8.DecodeRuneInString(t)

	for {
		if !unicode.IsSpace(r) {
			break
		}

		data = s.Scan()
		if !data {
			return ""
		}

		t = s.Text()
		r, _ = utf8.DecodeRuneInString(t)
	}

	return t
}

////////////////////////////////////////////////////////////////

type Lexer struct {
	currToken  Token
	dataOutput []LexemeData
	dataBuffer *bufio.Scanner
}

func NewLexer(s *bufio.Scanner) *Lexer {
	return &Lexer{
		dataBuffer: s,
	}
}

func (l *Lexer) Lex() {
	lexeme := ""
	eofFound := false

	currentChar := l.Next()
	if l.IsEOF(currentChar.char) {
		return
	}

	if currentChar.class == digit {
		lexeme += currentChar.char
		currentChar = l.Next()
		if l.IsEOF(currentChar.char) {
		}

		for {
			if currentChar.class != digit {
				break
			}

			lexeme += currentChar.char
			currentChar = l.Next()
			if l.IsEOF(currentChar.char) {
				eofFound = true
				break
			}
		}

		l.currToken = IntLit
		l.AppendToData(LexemeData{l.currToken, lexeme})
		lexeme = ""
		if eofFound {
			return
		}
	}

	if currentChar.class == letter {
		lexeme += currentChar.char
		l.currToken = lookupKeyword(lexeme)
		if l.currToken != Ident {
			l.AppendToData(LexemeData{l.currToken, lexeme})
			lexeme = ""
			return
		}

		currentChar = l.Next()
		if l.IsEOF(currentChar.char) {
			return
		}

		for {
			if currentChar.class != letter {
				break
			}

			lexeme += currentChar.char
			l.currToken = lookupKeyword(lexeme)
			if l.currToken != Ident {
				l.AppendToData(LexemeData{l.currToken, lexeme})
				lexeme = ""
				return
			}

			currentChar = l.Next()
			if l.IsEOF(currentChar.char) {
				eofFound = true
				break
			}
		}

		l.currToken = lookupKeyword(lexeme)
		l.AppendToData(LexemeData{l.currToken, lexeme})
		if eofFound {
			return
		}
	}

	if currentChar.class == unknown {
		l.currToken = lookup(currentChar.char)
		l.AppendToData(LexemeData{l.currToken, currentChar.char})
	}
}

func (l *Lexer) Next() CharType {
	return getCharType(getNextNonBlankChar(l.dataBuffer))
}

func (l *Lexer) AppendToData(data LexemeData) {
	l.dataOutput = append(l.dataOutput, data)
}

func (l *Lexer) IsEOF(char string) bool {
	if char == "" {
		l.currToken = Eof
		return true
	}

	return false
}
