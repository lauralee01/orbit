package rules

import (
	"testing"
)

func TestEvaluate(t *testing.T) {
	facts := Facts{
		"name": "John",
		"age":  30,
	}
	rules := Rules{
		{Field: "name", Operator: "==", Value: "John"},
		{Field: "age", Operator: "equals", Value: "30"},
	}

	got, err := Evaluate(facts, rules)
	if err != nil {
		t.Errorf("Evaluate() error = %v", err)
	}
	if !got {
		t.Errorf("Evaluate() = %v, want %v", got, true)
	}
}
