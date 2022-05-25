package validator

import "regexp"

var (
	emailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)


type Validator struct {
	Errors map[string]string
}

type Validatable interface {
	Validate (*Validator)
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Validate (i Validatable) bool {
	i.Validate(v)

	return v.Valid()
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}

	return false
}

func IsEmail(value string) bool {
	return emailRX.MatchString(value)
}
