package model

// Translate gherkin structures to ADM structures. Uses 'ModelLink'
// to capture relationships between model structures. These references
// are set as part of applying ADM rules.

// References can only point to attack, defense or policy
type TypeConstraint interface {
	*Attack | *Defense | []*Attack | []*Defense
}

// Holds the Step definition and/or references the target entity referred
// by this step. References are resolved using ADM rules. When referenced,
// the 'Item' field holds the reference.
type ModelLink[T TypeConstraint] struct {
	Step *Step
	Item T
}
