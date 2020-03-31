package main

import (
	"os"
	"strconv"
)

////////////////////////////////////////////////////////////////

// PopBinding -
func PopBinding(bindings *[]Binding) Binding {
	b := *bindings
	binding := b[0]
	b = b[1:]
	*bindings = b

	return binding
}

// Lookup -
func Lookup(bindings []Binding, varName string) string {
	val := ""
	for i := range bindings {
		if bindings[i].varName == varName {
			// Keep attempting to set a value because we want the
			// last viable variable within the environment (ie right most)
			val = bindings[i].value
		}
	}

	return val
}

// StrToBool -
func StrToBool(s string) bool {
	if s == "true" {
		return true
	}

	return false
}

////////////////////////////////////////////////////////////////

// Evaluator -
type Evaluator struct {
	astRoot *AstNode
}

// CreateEvaluator -
func CreateEvaluator(root *AstNode) Evaluator {
	return Evaluator{root}
}

// Evaluate -
func (e *Evaluator) Evaluate() string {
	return e.evaluate(e.astRoot, []Binding{})
}

func (e *Evaluator) evaluate(localParent *AstNode, bindings []Binding) string {
	localParent.environment = bindings
	switch localParent.tokenType {
	case LetKeyword:
		varName := localParent.children[0].tokenValue
		expOneVal := e.evaluate(localParent.children[1], bindings)

		bindings = append(bindings, Binding{varName, expOneVal})

		return e.evaluate(localParent.children[2], bindings)

	case MinusKeyword:
		expOneVal, err := strconv.Atoi(e.evaluate(localParent.children[0], bindings))
		if err != nil {
			os.Exit(1)
		}

		expTwoVal, err := strconv.Atoi(e.evaluate(localParent.children[1], bindings))
		if err != nil {
			os.Exit(1)
		}

		return strconv.Itoa(expOneVal - expTwoVal)

	case IszeroKeyword:
		expVal, err := strconv.Atoi(e.evaluate(localParent.children[0], bindings))
		if err != nil {
			os.Exit(1)
		}

		return strconv.FormatBool(expVal == 0)

	case IfKeyword:
		expValBool := StrToBool(e.evaluate(localParent.children[0], bindings))
		if expValBool {
			return e.evaluate(localParent.children[1], bindings)
		}

		return e.evaluate(localParent.children[2], bindings)

	case Ident:
		return Lookup(bindings, localParent.tokenValue)
	case IntLit:
		return localParent.tokenValue
	}

	return ""
}
