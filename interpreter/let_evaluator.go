package main

// Binding represents a pairing of a variable and a value
type Binding struct {
	varName string
	value   string
}

////////////////////////////////////////////////////////////////

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

// PushBinding -
func (e *Evaluator) PushBinding(b Binding) {
	e.bindings = PushBinding(e.bindings, b)
}

// Lookup -
func (e *Evaluator) Lookup(varName string) string {
	return Lookup(e.bindings, varName)
}

// Evaluate -
func (e *Evaluator) Evaluate() string {
	return e.evaluate(e.astRoot)
}

func (e *Evaluator) evaluate(localParent *AstNode) string {
	return ""
}
