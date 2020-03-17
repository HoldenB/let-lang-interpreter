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
	filename := "example.let"

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
		if lexer.currToken == Eof {
			lexer.AppendToData(LexemeData{Eof, "EOF"})
			break
		}
		lexer.Lex()
	}

	for _, data := range lexer.dataOutput {
		fmt.Printf("Token: %d | Lexeme: %s\n", data.token, data.lexeme)
	}
}
