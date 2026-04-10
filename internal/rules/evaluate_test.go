package rules

import (
	"fmt"
	"testing"
)

func TestEvaluate(t *testing.T) {
	facts := Facts{
		"name": "John",
		"age": 30,
	}
	rules := Rules{
		{Field: "name", Operator: "==", Value: "John"},
		{Field: "age", Operator: "equals", Value: "30"},
	}

	got, err := Evaluate(facts, rules)
    fmt.Println(got, err)
}