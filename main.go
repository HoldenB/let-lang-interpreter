package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

////////////////////////////////////////////////////////////////

func main() {
	filename := "example_3.let"

	filebuffer, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	inputdata := string(filebuffer)
	data := bufio.NewScanner(strings.NewReader(inputdata))
	data.Split(bufio.ScanRunes)

	lexer := NewLexer(data)
	for {
		if lexer.currTokenType == EOF {
			lexer.EnqueueToken(Token{EOF, "EOF"})
			break
		}
		lexer.Lex()
	}

	for _, data := range lexer.tokenQueue {
		fmt.Printf("Token: %d | Lexeme: %s\n", data.tokenType, data.tokenValue)
	}

	println()
	root := ParseTokenStream(lexer.tokenQueue)
	println()
	PrintTreeBasic(root)
	println()
	PrintTree(root)
}
