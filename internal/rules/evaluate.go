package rules

import "fmt"

func Evaluate(facts Facts, rules Rules) (bool, error) {
	for _, rule := range rules {
		factValue, ok := facts[rule.Field]
		if !ok {
			return false, fmt.Errorf("fact %s is not found", rule.Field)
		}
		if rule.Operator == "equals" || rule.Operator == "==" {
			if fmt.Sprintf("%v", factValue) != rule.Value {
				return false, fmt.Errorf("fact value %v does not equal %s", factValue, rule.Value)

			}
		}
	}
	return true, nil
}
