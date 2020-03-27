package main

import (
	"fmt"
	"os"
	"strconv"
)

// Binding represents a pairing of a variable and a value
type Binding struct {
	varName string
	value   string
}

////////////////////////////////////////////////////////////////

// PopBinding -
func PopBinding(bindings *[]Binding) Binding {
	b := *bindings
	binding := b[0]
	b = b[1:]
	*bindings = b

	return binding
}

// PushBinding -
func PushBinding(bindings []Binding, b Binding) []Binding {
	newBindings := []Binding{b}
	bindings = append(newBindings, bindings...)

	return bindings
}

// Lookup -
func Lookup(bindings []Binding, varName string) string {
	for i := range bindings {
		if bindings[i].varName == varName {
			return bindings[i].value
		}
	}

	return ""
}

////////////////////////////////////////////////////////////////

// Evaluator -
type Evaluator struct {
	astRoot  *AstNode
	bindings []Binding
}

// CreateEvaluator -
func CreateEvaluator(root *AstNode) Evaluator {
	return Evaluator{root, []Binding{}}
}

// PushBinding -
func (e *Evaluator) PushBinding(b Binding) {
	e.bindings = PushBinding(e.bindings, b)
}

// PopBinding -
func (e *Evaluator) PopBinding() Binding {
	return PopBinding(&e.bindings)
}

// Lookup -
func (e *Evaluator) Lookup(varName string) string {
	return Lookup(e.bindings, varName)
}

// Evaluate -
func (e *Evaluator) Evaluate() string {
	return e.evaluate(e.astRoot, e.bindings)
}

func (e *Evaluator) evaluate(localParent *AstNode, bindings []Binding) string {
	fmt.Printf("Evaluating token type: %s\n", printTokenNameVerbose(localParent.tokenType))
	switch localParent.tokenType {
	case LetKeyword:
		fmt.Printf("Let keyword found\n")
		varName := localParent.children[0].tokenValue
		expOneVal := e.evaluate(localParent.children[1], bindings)

		e.PushBinding(Binding{varName, expOneVal})

		return e.evaluate(localParent.children[2], bindings)
	case MinusKeyword:
		fmt.Printf("Minus keyword found\n")
		expOneVal, err := strconv.Atoi(e.evaluate(localParent.children[0], bindings))
		if err != nil {
			// Dirty but it'll work for now :(
			os.Exit(1)
		}
		expTwoVal, err := strconv.Atoi(e.evaluate(localParent.children[1], bindings))
		if err != nil {
			os.Exit(1)
		}

		return strconv.Itoa(expOneVal - expTwoVal)
	case IszeroKeyword:
		fmt.Printf("Iszero keyword found\n")
		return ""
	case IfKeyword:
		fmt.Printf("If keyword found\n")
		return ""
	case Ident:
		fmt.Printf("Ident found\n")
		return e.Lookup(localParent.tokenValue)
	case IntLit:
		fmt.Printf("IntLit found\n")
		return localParent.tokenValue
	}

	return ""
}
