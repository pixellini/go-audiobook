package coqui

import "fmt"

// Validator defines the interface for types that can validate themselves
// and provide string representation. This is commonly used for enum-like types.
type Validator interface {
	fmt.Stringer // Embeds String() string method
	IsValid() bool
}

// Compile-time interface checks to ensure our types implement Validator
var (
	_ Validator = (*Model)(nil)
	_ Validator = (*Language)(nil)
	_ Validator = (*Device)(nil)
)

// ValidateAll validates multiple Validator instances and returns the first error encountered
func ValidateAll(validators ...Validator) error {
	for _, v := range validators {
		if !v.IsValid() {
			return fmt.Errorf("invalid value: %s", v.String())
		}
	}
	return nil
}

// MustValidate panics if any of the provided validators are invalid
func MustValidate(validators ...Validator) {
	if err := ValidateAll(validators...); err != nil {
		panic(err)
	}
}

// IsValidValidator is a convenience function that safely checks if a value
// implements Validator and is valid, returning false for nil or invalid values
func IsValidValidator(v Validator) bool {
	return v != nil && v.IsValid()
}

// ValidateWithContext validates a Validator with additional context information
func ValidateWithContext(v Validator, context string) error {
	if !v.IsValid() {
		return fmt.Errorf("invalid %s: %s", context, v.String())
	}
	return nil
}
