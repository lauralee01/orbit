package rules

type Rule struct {
	Field string
	Operator string
	Value string
}

type Rules []Rule
