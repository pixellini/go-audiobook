package coqui

import "fmt"

// Validator defines the interface for types that can validate themselves.
// Types implementing this interface can check their own validity and provide
// string representation, commonly used for enum-like types.
type Validator interface {
	fmt.Stringer // Embeds String() string method
	IsValid() bool
}

// Compile-time interface checks to ensure our types implement Validator.
// These variables verify at compile time that Model, Language, and Device
// all properly implement the Validator interface.
var (
	_ Validator = (*Model)(nil)
	_ Validator = (*Language)(nil)
	_ Validator = (*Device)(nil)
)

// ValidateAll validates multiple Validator instances in sequence.
// Returns the first validation error encountered, or nil if all are valid.
func ValidateAll(validators ...Validator) error {
	for _, v := range validators {
		if !v.IsValid() {
			return fmt.Errorf("invalid value: %s", v.String())
		}
	}
	return nil
}

// MustValidate validates multiple Validator instances and panics if any are invalid.
// Use this for validation that should never fail in correct program flow.
func MustValidate(validators ...Validator) {
	if err := ValidateAll(validators...); err != nil {
		panic(err)
	}
}

// IsValidValidator safely checks if a value implements Validator and is valid.
// Returns false for nil values or invalid validators, making it nil-safe.
func IsValidValidator(v Validator) bool {
	return v != nil && v.IsValid()
}

// ValidateWithContext validates a Validator with additional context information.
// Provides more descriptive error messages by including the context of what failed.
func ValidateWithContext(v Validator, context string) error {
	if !v.IsValid() {
		return fmt.Errorf("invalid %s: %s", context, v.String())
	}
	return nil
}
