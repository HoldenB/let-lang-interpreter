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
		if _, err := os.Stat(reader.Text()); !os.IsNotExist(err) {
			file = reader.Text()
			fmt.Println("\n==============================================")
			fmt.Println("Executing " + file)
			fmt.Println("==============================================")
			fmt.Println("\n...\n...\n...")
			break
		} else {
			fmt.Println("File not in directory")
		}
	}

	filepath := file

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
	fmt.Println("==============================================")
	fmt.Println("TOKEN QUEUE")
	fmt.Println("==============================================\n")
	for _, data := range lexer.tokenQueue {
		fmt.Printf("Token: %d | Lexeme: %s\n", data.tokenType, data.tokenValue)
	}

	root := ParseTokenStream(lexer.tokenQueue)
	println()

	fmt.Println("==============================================")
	fmt.Println("ABSTRACT SYNTAX TREE (WITHOUT ENVIRONMENT)")
	fmt.Println("==============================================\n")
	PrintTree(root)
	println()
	eval := CreateEvaluator(root)

	fmt.Println("\n==============================================")
	fmt.Println("EVALUATION")
	fmt.Println("==============================================")
	fmt.Printf("\nExpression evaluated to: %s\n", eval.Evaluate())
	println()

	fmt.Println("==============================================")
	fmt.Println("ABSTRACT SYNTAX TREE (WITH ENVIRONMENT)")
	fmt.Println("==============================================\n")
	PrintTree(eval.astRoot)
	println()
	fmt.Println("\n==============================================")
	fmt.Println("Done.")
}
