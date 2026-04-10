package rules
import "fmt"


func Evaluate(facts Facts, rules Rules) (bool, error) {
	for _, rule := range rules {
		fmt.Println("Evaluating rule:", rule.Field, rule.Operator, rule.Value)
		factValue, ok := facts[rule.Field]
		fmt.Println("Fact value:", factValue, "Exists:", ok)
		if !ok {
			return false, fmt.Errorf("fact %s is not found", rule.Field)
		}
		if rule.Operator == "equals" || rule.Operator == "==" {
			if(fmt.Sprintf("%v", factValue) != rule.Value) {
				return false, fmt.Errorf("fact value %v does not equal %s", factValue, rule.Value)

			}
		}
	}
	return true, nil
}



