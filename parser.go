package main

import (
	"fmt"
	"os"
)

// AstNode represents a Node in an abstract syntax tree
type AstNode struct {
	// Pointer to parent node if exists. If null, then root.
	parent *AstNode
	// Child nodes
	children []*AstNode
	// Token type, used to distinguish btw vars and integers.
	tokenType TokenType
	// Represent the contents as a string, even if it’s an int.
	tokenValue string
	// Is it a terminal symbol (leaf node)?
	isLeaf bool
}

func (node *AstNode) printAST(indentLevel int) {
	outString := ""
	for i := 0; i < indentLevel; i++ {
		outString += " "
	}

	fmt.Printf("%s", outString)
	fmt.Printf("%s \n", node.tokenValue)

	if node.isLeaf == false && len(node.children) > 0 {
		for i := 0; i < len(node.children); i++ {
			node.children[i].printAST(indentLevel + 1)
		}
	}
}

////////////////////////////////////////////////////////////////

// PopToken will pop a token off the front of the queue and modify the
// queue by removing the first element and replacing the old slice
func PopToken(tokenQueue *[]Token) Token {
	tq := *tokenQueue
	token := tq[0]
	tq = tq[1:]
	*tokenQueue = tq
	return token
}

// PrintTree -
func PrintTree(root *AstNode) {
	root.printAST(0)
}

// ParseTokenStream -
func ParseTokenStream(tokenQueue []Token) *AstNode {
	p := Parser{tokenQueue, &AstNode{}}
	_ = p.parseExp()

	return p.root
}

////////////////////////////////////////////////////////////////

// Parser -
type Parser struct {
	tokenQueue []Token
	root       *AstNode
}

func (p *Parser) advanceToken() {
	_ = PopToken(&p.tokenQueue)
}

func (p *Parser) peekNextToken() (TokenType, string) {
	if len(p.tokenQueue) <= 0 {
		return UnknownType, ""
	}

	return p.tokenQueue[0].tokenType, p.tokenQueue[0].tokenValue
}

func (p *Parser) checkExpectedToken(t TokenType, errString string) {
	nextType, _ := p.peekNextToken()

	// Advance input stream if we have our expected token
	if t != nextType {
		println(errString)
		// Dirty but it'll work for now :(
		os.Exit(1)
	}

	p.advanceToken()
}

func (p *Parser) initLocalParentNode(parent *AstNode) {
	// Initialize
	parent.isLeaf = true
	parent.children = make([]*AstNode, 0, 5)
	parent.tokenType, parent.tokenValue = p.peekNextToken()

	// Advance the token stream
	p.advanceToken()
}

func (p *Parser) parseExp() AstNode {
	parentNode := AstNode{}

	// Should only be set on our first call
	if p.root.children == nil {
		p.root = &parentNode
	}

	p.initLocalParentNode(&parentNode)

	println(parentNode.tokenValue)
	switch parentNode.tokenType {
	case MinusKeyword:
		println("minus keyword found")
		parentNode.isLeaf = false
		p.checkExpectedToken(LeftParen, "unexpected token, expected left paren")
		println("left paren found, advancing input")

		leftChild := p.parseExp()
		fmt.Printf("left child found: %s\n", leftChild.tokenValue)

		p.checkExpectedToken(Comma, "unexpected token, expected comma")

		println("comma found, advancing input")

		rightChild := p.parseExp()
		fmt.Printf("right child found: %s\n", rightChild.tokenValue)

		p.checkExpectedToken(RightParen, "unexpected token, expected right paren")

		println("right paren found, advancing input")

		// We have valid children and must set them on the parent
		leftChild.parent = &parentNode
		parentNode.children = append(parentNode.children, &leftChild)

		rightChild.parent = &parentNode
		parentNode.children = append(parentNode.children, &rightChild)
	}

	return parentNode
}
