package validator

import "strings"

const (
	MinPasswordLength  = 8
	MinPasswordEntropy = 60

	Required = "must not be empty"
)

type Validator struct {
	errs map[string]string
}

func New() *Validator {
	return &Validator{
		errs: make(map[string]string),
	}
}

// Valid ensures that there is no error after validation
func (v *Validator) Valid() bool {
	return len(v.errs) == 0
}

// Errors returns all validation errors
func (v *Validator) Errors() map[string]string {
	return v.errs
}

// AddError adds new error field to the validation errors map
func (v *Validator) AddError(field, message string) {
	if v.errs == nil {
		v.errs = make(map[string]string)
	}

	_, exist := v.errs[field]
	if !exist {
		v.errs[field] = message
	}
}

// RequiredString ensures that a string is not empty
func (v *Validator) RequiredString(s, field, message string) {
	empty := len(strings.TrimSpace(s)) == 0
	if empty {
		v.AddError(field, message)
	}
}

// MinString ensure that does not contain less characters than min
func (v *Validator) MinString(s string, min int, field string, message string) {
	lesser := len(strings.TrimSpace(s)) < min
	if lesser {
		v.AddError(field, message)
	}
}
