package query

// Logic combines predicates and nested groups.
type Logic string

const (
	And Logic = "and"
	Or  Logic = "or"
)

// Predicate is a single SQL condition.
type Predicate struct {
	Column   string
	Operator Operator
	Value    any
	SQL      string
	Args     []any
}

// Group is a nested AND/OR condition tree.
type Group struct {
	Logic      Logic
	Predicates []Predicate
	Groups     []Group
}
