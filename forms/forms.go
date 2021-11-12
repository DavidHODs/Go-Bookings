package forms

import (
	"fmt"
	"net/url"

	"github.com/asaskevich/govalidator"
)

// Form creates a custom form struct, embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// New initializes a new form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Has checks if form field is not empty
func (f *Form) Has(field string) bool {
	x1 := f.Get(field)

	if x1 == "" {
		f.Errors.Add(field, "This field cannot be blank")
		return false
	}
	
	return true
}

// Valid returns true if there are no errors i.e form is Valid
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// MinLength checks for string minimum length
func(f *Form) MinLength(field string, length int) bool {
	x := f.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be atleast 3 characters long"))
		return false
	}

	return true
}

// IsEmail checks if email supplied is valid
func (f *Form) IsEmail(field string) bool {
	x := f.Get(field)
	if !govalidator.IsEmail(x) {
		f.Errors.Add(field, "Invalid Email Address")
		return false
	}

	return true
}