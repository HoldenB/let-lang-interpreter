package main

import (
	"bufio"
	"unicode"
	"unicode/utf8"
)

////////////////////////////////////////////////////////////////

// CharClass enum
type CharClass int

// TokenType enum
type TokenType int

const (
	// EOF is an End Of File representation
	EOF = -1

	letter  CharClass = 0
	digit   CharClass = 1
	unknown CharClass = 99

	Ident         TokenType = 23
	IntLit        TokenType = 24
	LeftParen     TokenType = 25
	RightParen    TokenType = 26
	Comma         TokenType = 27
	EqualSign     TokenType = 28
	MinusKeyword  TokenType = 29
	IszeroKeyword TokenType = 30
	IfKeyword     TokenType = 31
	ThenKeyword   TokenType = 32
	ElseKeyword   TokenType = 33
	LetKeyword    TokenType = 34
	InKeyword     TokenType = 35
	UnknownType   TokenType = 36
)

// CharType represents a character and it's
// class
type CharType struct {
	char  string
	class CharClass
}

// Token is an object that contains
// a parsed token type and value
type Token struct {
	tokenType  TokenType
	tokenValue string
}

////////////////////////////////////////////////////////////////
// Helper functions

func lookup(char string) TokenType {
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

	return EOF
}

func lookupKeyword(s string) TokenType {
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
		return CharType{char, EOF}
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

// Lexer allows the parsing/lexing of byte data into tokens
type Lexer struct {
	currTokenType TokenType
	tokenQueue    []Token
	dataBuffer    *bufio.Scanner
}

// NewLexer creates a new lexer object given an initalized bufio.Scanner
func NewLexer(s *bufio.Scanner) *Lexer {
	return &Lexer{
		dataBuffer: s,
	}
}

// Lex provides a byte-stream waterfall parse. Lex will
// yield tokens based on parsing incoming characters from the
// configured byte-stream
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

		l.currTokenType = IntLit
		l.EnqueueToken(Token{l.currTokenType, lexeme})
		lexeme = ""
		if eofFound {
			return
		}
	}

	if currentChar.class == letter {
		lexeme += currentChar.char
		l.currTokenType = lookupKeyword(lexeme)
		if l.currTokenType != Ident {
			l.EnqueueToken(Token{l.currTokenType, lexeme})
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
			l.currTokenType = lookupKeyword(lexeme)
			if l.currTokenType != Ident {
				l.EnqueueToken(Token{l.currTokenType, lexeme})
				lexeme = ""
				return
			}

			currentChar = l.Next()
			if l.IsEOF(currentChar.char) {
				eofFound = true
				break
			}
		}

		l.currTokenType = lookupKeyword(lexeme)
		l.EnqueueToken(Token{l.currTokenType, lexeme})
		if eofFound {
			return
		}
	}

	if currentChar.class == unknown {
		l.currTokenType = lookup(currentChar.char)
		l.EnqueueToken(Token{l.currTokenType, currentChar.char})
	}
}

// Next will get the next non blank character type from the
// data buffer
func (l *Lexer) Next() CharType {
	return getCharType(getNextNonBlankChar(l.dataBuffer))
}

// EnqueueToken takes an input Token and appends it to the
// end of the TokenQueue
func (l *Lexer) EnqueueToken(t Token) {
	l.tokenQueue = append(l.tokenQueue, t)
}

// IsEOF takes a character input and returns true
// if the character is empty, false otherwise. When taking
// input from the data buffer, an empty character represents the
// end of the buffer
func (l *Lexer) IsEOF(char string) bool {
	if char == "" {
		l.currTokenType = EOF
		return true
	}

	return false
}
