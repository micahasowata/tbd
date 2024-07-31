package validator

<<<<<<< HEAD
import "strings"

=======
>>>>>>> caa53d69f37ffc231928ff06ffb56c8fc0fc8915
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

func (v *Validator) Valid() bool {
	return len(v.errs) == 0
}

func (v *Validator) Errors() map[string]string {
	return v.errs
}

func (v *Validator) AddError(field, message string) {
	_, exist := v.errs[field]
	if !exist {
		v.errs[field] = message
	}
}

func (v *Validator) RequiredString(s, field, message string) {
<<<<<<< HEAD
	empty := len(strings.TrimSpace(s)) == 0
=======
	empty := len(s) == 0
>>>>>>> caa53d69f37ffc231928ff06ffb56c8fc0fc8915
	if empty {
		v.AddError(field, message)
	}
}

func (v *Validator) MinString(s string, length int, field string, message string) {
<<<<<<< HEAD
	lesser := len(strings.TrimSpace(s)) < length
=======
	lesser := len(s) < length
>>>>>>> caa53d69f37ffc231928ff06ffb56c8fc0fc8915
	if lesser {
		v.AddError(field, message)
	}
}
