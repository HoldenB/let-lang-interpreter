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
	reader := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter .let file to parse. Must be in this directory.")
	file := ""
	for reader.Scan() {
		if _, err := os.Stat("interpreter/" + reader.Text()); !os.IsNotExist(err) {
			file = reader.Text()
			break
		} else {
			fmt.Println("File not in directory")
		}
	}

	filepath := "interpreter/" + file

	filebuffer, err := ioutil.ReadFile(filepath)
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

	println()

	for _, data := range lexer.tokenQueue {
		fmt.Printf("Token: %d | Lexeme: %s\n", data.tokenType, data.tokenValue)
	}

	root := ParseTokenStream(lexer.tokenQueue)
	println()
	PrintTree(root)
	println()
	eval := CreateEvaluator(root)
	fmt.Printf("\nExpression evaluated to: %s\n", eval.Evaluate())
	println()
	PrintTree(eval.astRoot)
	println()
}
