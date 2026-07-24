package query

// Operator defines a supported SQL comparison operator.
type Operator string

const (
	OpEqual    Operator = "="
	OpNotEqual Operator = "!="
	OpLike     Operator = "like"
	OpIn       Operator = "in"
	OpIsNull   Operator = "is_null"
	OpNotNull  Operator = "is_not_null"
	OpRaw      Operator = "raw"
)
