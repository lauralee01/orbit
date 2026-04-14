package rules

// Rule is a struct that represents a rule for the rules engine. it sets the conditions for the rule.
type Rule struct {
	Field    string
	Operator string
	Value    string
}

type Rules []Rule
