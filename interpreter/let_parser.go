package main

import (
	"fmt"
	"os"
)

// Binding represents a pairing of a variable and a value
type Binding struct {
	varName string
	value   string
}

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
	// Environment (list of bindings)
	environment []Binding
}

func (node *AstNode) printASTbasic(indentLevel int) {
	indentStr := ""
	for i := 0; i < indentLevel; i++ {
		indentStr += " "
	}

	fmt.Printf("%s", indentStr)
	fmt.Printf("%s (%s) \n", node.tokenValue, printTokenNameVerbose(node.tokenType))

	if !node.isLeaf && len(node.children) > 0 {
		for i := 0; i < len(node.children); i++ {
			node.children[i].printASTbasic(indentLevel + 1)
		}
	}
}

func (node *AstNode) printAST(indentLevel int) {
	indentStr := ""
	for i := 0; i < indentLevel; i++ {
		indentStr += "    "
	}

	fmt.Printf("%s", indentStr)
	fmt.Printf("%s", printTokenName(node.tokenType))
	for i, b := range node.environment {
		if i == 0 {
			fmt.Printf(" -> Env [")
		}

		if i > 0 && i < len(node.environment) {
			fmt.Printf(", ")
		}

		fmt.Printf("%s", printBinding(b))

		if i == len(node.environment)-1 {
			fmt.Printf("]")
		}
	}

	paren := needsParen(node.tokenType)
	if paren {
		fmt.Print(" (\n")
	}

	if !node.isLeaf && len(node.children) > 0 {
		for i := 0; i < len(node.children); i++ {
			node.children[i].printAST(indentLevel + 1)

			// Special case for minus keywords
			if node.tokenType == MinusKeyword && i == 0 {
				fmt.Printf(",")
			}

			// Special case for let keywords
			if node.tokenType == LetKeyword && i == 1 {
				fmt.Printf(",")
			}

			fmt.Printf("\n")
		}
	}

	indentStrIdent := ""
	for i := 0; i < len(printTokenName(node.tokenType))-2; i++ {
		indentStrIdent += " "
	}

	// Idents need double quotes
	if node.tokenType == Ident {
		fmt.Printf("%s", indentStr+indentStrIdent)
		fmt.Printf("\"%s\"\n", node.tokenValue)
	}

	// Const cannot have double quotes
	if node.tokenType == IntLit {
		fmt.Printf("%s", indentStr+indentStrIdent)
		fmt.Printf("%s\n", node.tokenValue)
	}

	fmt.Printf("%s", indentStr)
	if paren {
		fmt.Print(")")
	}
}

func needsParen(token TokenType) bool {
	if token == Ident ||
		token == IntLit ||
		token == MinusKeyword ||
		token == IszeroKeyword ||
		token == IfKeyword ||
		token == LetKeyword {
		return true
	}

	return false
}

func printTokenName(token TokenType) string {
	switch token {
	case Ident:
		return "VarExp"
	case IntLit:
		return "ConstExp"
	case LeftParen:
		return ""
	case RightParen:
		return ""
	case Comma:
		return ""
	case EqualSign:
		return ""
	case MinusKeyword:
		return "DiffExp"
	case IszeroKeyword:
		return "IsZeroExp"
	case IfKeyword:
		return "IfExp"
	case ThenKeyword:
		return ""
	case ElseKeyword:
		return ""
	case LetKeyword:
		return "LetExp"
	case InKeyword:
		return ""
	case UnknownType:
		return ""
	default:
		return ""
	}
}

func printTokenNameVerbose(token TokenType) string {
	switch token {
	case Ident:
		return "Ident"
	case IntLit:
		return "IntLit"
	case LeftParen:
		return "LeftParen"
	case RightParen:
		return "RightParen"
	case Comma:
		return "Comma"
	case EqualSign:
		return "EqualSign"
	case MinusKeyword:
		return "MinusKeyword"
	case IszeroKeyword:
		return "iszeroKeyword"
	case IfKeyword:
		return "IfKeyword"
	case ThenKeyword:
		return "ThenKeyword"
	case ElseKeyword:
		return "ElseKeyword"
	case LetKeyword:
		return "LetKeyword"
	case InKeyword:
		return "InKeyword"
	case UnknownType:
		return "UnknownType"
	default:
		return ""
	}
}

func printBinding(b Binding) string {
	return "(" + b.varName + ", " + b.value + ")"
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

// PrintTreeBasic -
func PrintTreeBasic(root *AstNode) {
	root.printASTbasic(0)
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

func (p *Parser) checkExpectedToken(t TokenType, advanceInput bool, errString string) {
	nextType, _ := p.peekNextToken()

	// Advance input stream if we have our expected token
	if t != nextType {
		println(errString)
		os.Exit(1)
	}

	if advanceInput {
		p.advanceToken()
	}
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

	switch parentNode.tokenType {
	case MinusKeyword:
		parentNode.isLeaf = false

		p.checkExpectedToken(LeftParen, true, "unexpected token, expected left paren")
		leftChild := p.parseExp()
		p.checkExpectedToken(Comma, true, "unexpected token, expected comma")
		rightChild := p.parseExp()
		p.checkExpectedToken(RightParen, true, "unexpected token, expected right paren")

		// We have valid children and must set them on the parent
		leftChild.parent = &parentNode
		parentNode.children = append(parentNode.children, &leftChild)

		rightChild.parent = &parentNode
		parentNode.children = append(parentNode.children, &rightChild)

	case IszeroKeyword:
		parentNode.isLeaf = false
		p.checkExpectedToken(LeftParen, true, "unexpected token, expected left paren")
		childExp := p.parseExp()
		p.checkExpectedToken(RightParen, true, "unexpected token, expected right paren")

		// Valid child expression
		childExp.parent = &parentNode
		parentNode.children = append(parentNode.children, &childExp)

	case IfKeyword:
		// For now we'll follow the grammar in the sense that if we
		// have an "if" keyword, we expect 3 total expressions in the form
		// of: if exp then exp else exp
		parentNode.isLeaf = false

		// Currently our only predicate keyword is "iszero"
		// We do not want to advance input when checking for the expected token,
		// because we need to evaluate the predicate keyword
		p.checkExpectedToken(IszeroKeyword, false, "unexpected token, expected iszero keyword")
		predicateExp := p.parseExp()
		p.checkExpectedToken(ThenKeyword, true, "unexpected token, expected then keyword")
		caseFalseExp := p.parseExp()
		p.checkExpectedToken(ElseKeyword, true, "unexpected token, expected else keyword")
		caseTrueExp := p.parseExp()

		// Valid if then else statement
		predicateExp.parent = &parentNode
		parentNode.children = append(parentNode.children, &predicateExp)

		caseFalseExp.parent = &parentNode
		parentNode.children = append(parentNode.children, &caseFalseExp)

		caseTrueExp.parent = &parentNode
		parentNode.children = append(parentNode.children, &caseTrueExp)

	case LetKeyword:
		parentNode.isLeaf = false

		p.checkExpectedToken(Ident, false, "unexpected token, expected identifier")
		identifier := p.parseExp()
		p.checkExpectedToken(EqualSign, true, "unexpected token, expected assignment")
		childExpOne := p.parseExp()
		p.checkExpectedToken(InKeyword, true, "unexpected token, expected in keyword")
		childExpTwo := p.parseExp()

		// Valid let assignment
		identifier.parent = &parentNode
		parentNode.children = append(parentNode.children, &identifier)

		childExpOne.parent = &parentNode
		parentNode.children = append(parentNode.children, &childExpOne)

		childExpTwo.parent = &parentNode
		parentNode.children = append(parentNode.children, &childExpTwo)
	}

	return parentNode
}
