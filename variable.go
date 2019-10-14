package timy

// VariableType can be anything
type VariableType string

// Variable is a simple variable struct, with this you can create a simple fuzzy algorithm
type Variable struct {
	ID   string
	Name string
	Type VariableType
}
