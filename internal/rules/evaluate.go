package rules

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrMissingFact         = errors.New("fact is not found")
	ErrFactValueMismatch   = errors.New("fact value does not match")
	ErrUnsupportedOperator = errors.New("unsupported operator")
)

// asFloat64 coerces JSON-decoded numbers (float64, int, etc.) to float64.
func asFloat64(v any) (float64, error) {
	switch x := v.(type) {
	case float64:
		return x, nil
	case float32:
		return float64(x), nil
	case int:
		return float64(x), nil
	case int64:
		return float64(x), nil
	default:
		return 0, fmt.Errorf("expected numeric fact, got %T", v)
	}
}

func Evaluate(facts Facts, rules Rules) (bool, error) {
	for _, rule := range rules {
		factValue, ok := facts[rule.Field]
		if !ok {
			return false, fmt.Errorf("%w: fact %s is not found", ErrMissingFact, rule.Field)
		}

		switch rule.Operator {
		case "equals", "==":
			if fmt.Sprintf("%v", factValue) != rule.Value {
				return false, fmt.Errorf("%w: fact value %v does not equal %s", ErrFactValueMismatch, factValue, rule.Value)
			}
		case "greater_than", ">":
			fv, err := asFloat64(factValue)
			if err != nil {
				return false, fmt.Errorf("%w: %v", ErrFactValueMismatch, err)
			}
			rv, err := strconv.ParseFloat(rule.Value, 64)
			if err != nil {
				return false, fmt.Errorf("%w: invalid threshold %q", ErrFactValueMismatch, rule.Value)
			}
			if fv <= rv {
				return false, fmt.Errorf("%w: fact value %v is not greater than %s", ErrFactValueMismatch, factValue, rule.Value)
			}
		case "less_than", "<":
			fv, err := asFloat64(factValue)
			if err != nil {
				return false, fmt.Errorf("%w: %v", ErrFactValueMismatch, err)
			}
			rv, err := strconv.ParseFloat(rule.Value, 64)
			if err != nil {
				return false, fmt.Errorf("%w: invalid threshold %q", ErrFactValueMismatch, rule.Value)
			}
			if fv >= rv {
				return false, fmt.Errorf("%w: fact value %v is not less than %s", ErrFactValueMismatch, factValue, rule.Value)
			}
		case "greater_than_or_equal_to", ">=":
			fv, err := asFloat64(factValue)
			if err != nil {
				return false, fmt.Errorf("%w: %v", ErrFactValueMismatch, err)
			}
			rv, err := strconv.ParseFloat(rule.Value, 64)
			if err != nil {
				return false, fmt.Errorf("%w: invalid threshold %q", ErrFactValueMismatch, rule.Value)
			}
			if fv < rv {
				return false, fmt.Errorf("%w: fact value %v is not greater than or equal to %s", ErrFactValueMismatch, factValue, rule.Value)
			}
		case "less_than_or_equal_to", "<=":
			fv, err := asFloat64(factValue)
			if err != nil {
				return false, fmt.Errorf("%w: %v", ErrFactValueMismatch, err)
			}
			rv, err := strconv.ParseFloat(rule.Value, 64)
			if err != nil {
				return false, fmt.Errorf("%w: invalid threshold %q", ErrFactValueMismatch, rule.Value)
			}
			if fv > rv {
				return false, fmt.Errorf("%w: fact value %v is not less than or equal to %s", ErrFactValueMismatch, factValue, rule.Value)
			}
		default:
			return false, fmt.Errorf("%w: %s", ErrUnsupportedOperator, rule.Operator)
		}
	}
	return true, nil
}
