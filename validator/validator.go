package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	Errors map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddFieldError(key, message string) {
	if v.Errors == nil {
		v.Errors = make(map[string]string)
	}
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func (v *Validator) NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func (v *Validator) MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func (v *Validator) ValidEmail(email string) bool {
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailPattern.MatchString(email)
}

func (v *Validator) ValidPassword(password string) {
	v.CheckField(len(password) > 0, "password", "Password Should not be empty")
	if len(password) > 0 {
		emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%\-$@]{5,20}$`)
		v.CheckField(emailPattern.MatchString(password), "password", "Password Should contain Alphanumeric char and Special char (._%-$@) only between 5 to 20 char")
		// v.CheckField(len(password) < 20, "password", "Password Should be less than 20")
	}

}
