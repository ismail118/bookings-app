package forms

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/http"
	"net/url"
	"strings"
)

type Form struct {
	Data   url.Values
	Errors errors
}

// Valid return true if there is error, false otherwise
func (f *Form) Valid() bool {
	if len(f.Errors) == 0 {
		return true
	}

	return false
}

// New initialize a form struct
func New(data url.Values) *Form {
	return &Form{
		Data:   data,
		Errors: errors(make(map[string][]string)),
	}
}

// Required validate given multiple fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Data.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// Has validate given field
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		f.Errors.Add(field, "This field cannot be blank")
		return false
	}

	return true
}

// MinLength check min length for given fields
func (f *Form) MinLength(field string, minLength int, r *http.Request) bool {
	x := r.Form.Get(field)
	if len(x) < minLength {
		f.Errors.Add(field, fmt.Sprintf("Minimal %d character for this field", minLength))
		return false
	}

	return true
}

func (f *Form) IsEmail(field string) bool {
	if !govalidator.IsEmail(f.Data.Get(field)) {
		f.Errors.Add(field, "Invalid email")
		return false
	}

	return true
}
