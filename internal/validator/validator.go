package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// this is pointer to a "complied" regexp.Regexp type, parsing this during startup is more performant than reparsing everytime we need it

var EmailRX = regexp.MustCompile("^\\w+(?:\\.\\w+)*@\\w+(?:\\.[\\w-]+)*\\.[a-zA-Z]{2,}$")

// Validator struct that holds a map of the validation errors
type Validator struct {
	FieldErrors    map[string]string
	NonFieldErrors []string //validation errors not related to a particular field
}

// Valid method to check if there is any errors
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

// AddFieldError() adds an error message to the FieldErrors map (so long as no // entry already exists for the given key).
func (v *Validator) AddFieldError(key, message string) {
	// Note: We need to initialize the map first, if it isn't already // initialized.
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// helper function for adding non field errors
func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

// CheckField adds an error message to the FieldErrors map only if a // validation check is not 'ok'.
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// NotBlank() returns true if a value is not an empty string.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars() returns true if a value contains no more than n characters.
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// PermittedInt() returns true if a value is in a list of permitted integers.
func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}

// check is the length is greater than a given length
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// check if the given string matches a given regex pattern
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
