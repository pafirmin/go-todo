package validator

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var (
	emailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type Validator struct {
	Errors map[string]string
}

type Validatable interface {
	Validate(*Validator)
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Exec(i Validatable) {
	i.Validate(v)
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

func (v *Validator) ValidDatetime(key, value string) {
	if _, err := time.Parse(time.RFC3339, value); err != nil {
		v.AddError(key, "must be valid RFC3339 date string")
	}
}

func (v *Validator) ValidDate(key, value string) {
	if _, err := time.Parse("2006-02-01", value); err != nil {
		v.AddError(key, "must be in format dd-MM-yyyy")
	}
}

func (v *Validator) ValidLength(key, value string, min, max int) {
	l := len(value)
	if l < min || l > max {
		v.AddError(key, fmt.Sprintf("must be between %d and %d characters", min, max))
	}
}

func (v *Validator) ValidEmail(key, value string) {
	if emailRX.MatchString(value) {
		return
	}

	v.AddError(key, "must be valid email address")
}

func (v *Validator) PermittedValue(key string, value string, permittedValues ...string) {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return
		}
	}

	msg := "must be one of: " + strings.Join(permittedValues, ", ")
	v.AddError(key, msg)
}
