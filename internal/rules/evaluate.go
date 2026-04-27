package rules

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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

func equals(factValue any, ruleValue string) bool {
	// Numeric facts compare numerically when possible so 30 and "30" match.
	if fv, err := asFloat64(factValue); err == nil {
		if rv, parseErr := strconv.ParseFloat(ruleValue, 64); parseErr == nil {
			return fv == rv
		}
	}

	switch v := factValue.(type) {
	case string:
		return v == ruleValue
	case bool:
		rv, err := strconv.ParseBool(strings.ToLower(strings.TrimSpace(ruleValue)))
		return err == nil && v == rv
	default:
		return fmt.Sprint(factValue) == ruleValue
	}
}

// compareFloat performs a numeric comparison between a fact value and a rule value.
// It does three things:
//   1. Coerces the fact value (which may be float64, int, etc.) into a float64.
//   2. Parses the rule's string value into a float64.
//   3. Applies the provided comparison function (cmp) to determine if the rule passes.
//
// If the fact is not numeric, the rule value cannot be parsed, or the comparison fails,
// it returns an ErrFactValueMismatch with a descriptive message.

func compareFloat(factValue any, ruleValue, mismatchMsg string, cmp func(fv, rv float64) bool) error {
	fv, err := asFloat64(factValue)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFactValueMismatch, err)
	}
	rv, err := strconv.ParseFloat(ruleValue, 64)
	if err != nil {
		return fmt.Errorf("%w: invalid threshold %q", ErrFactValueMismatch, ruleValue)
	}
	if !cmp(fv, rv) {
		return fmt.Errorf("%w: fact value %v %s %s", ErrFactValueMismatch, factValue, mismatchMsg, ruleValue)
	}
	return nil
}

// Evaluate checks a set of facts against a list of rules.
// It processes each rule sequentially and returns:
//
//   - (true, nil) if *all* rules pass
//   - (false, error) as soon as any rule fails
//
// Rule evaluation is type-aware:
//   - "equals" supports numeric, string, and boolean comparisons
//   - numeric operators (>, <, >=, <=) require both sides to be numeric
//
// Errors are descriptive and include:
//   - ErrMissingFact when a required fact is not present
//   - ErrFactValueMismatch when a rule fails
//   - ErrUnsupportedOperator for unknown operators
//
// This function short-circuits on the first failure, which keeps evaluation fast
// and makes the returned error directly explain the reason for failure.

func Evaluate(facts Facts, rules Rules) (bool, error) {
	for _, rule := range rules {
		factValue, ok := facts[rule.Field]
		if !ok {
			return false, fmt.Errorf("%w: fact %s is not found", ErrMissingFact, rule.Field)
		}

		op := strings.ToLower(strings.TrimSpace(rule.Operator))

		switch op {
		case "equals", "==":
			if !equals(factValue, rule.Value) {
				return false, fmt.Errorf("%w: fact value %v does not equal %s", ErrFactValueMismatch, factValue, rule.Value)
			}
		case "greater_than", ">":
			if err := compareFloat(factValue, rule.Value, "is not greater than", func(fv, rv float64) bool { return fv > rv }); err != nil {
				return false, err
			}
		case "less_than", "<":
			if err := compareFloat(factValue, rule.Value, "is not less than", func(fv, rv float64) bool { return fv < rv }); err != nil {
				return false, err
			}
		case "greater_than_or_equal_to", ">=":
			if err := compareFloat(factValue, rule.Value, "is not greater than or equal to", func(fv, rv float64) bool { return fv >= rv }); err != nil {
				return false, err
			}
		case "less_than_or_equal_to", "<=":
			if err := compareFloat(factValue, rule.Value, "is not less than or equal to", func(fv, rv float64) bool { return fv <= rv }); err != nil {
				return false, err
			}
		default:
			return false, fmt.Errorf("%w: %s", ErrUnsupportedOperator, rule.Operator)
		}
	}
	return true, nil
}
