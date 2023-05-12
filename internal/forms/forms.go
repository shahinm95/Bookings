package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Form creates a custm form struct , embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// New initialize a form struct
func New (data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required check all given field to not be empty
func (f *Form) Required (fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field can not be blank")
		}
	}
}


// Has check if form field is in post and not empty
func (f *Form) Has (field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		f.Errors.Add(field , "This field can not be empty")
		return false
	}
	return true
}

// Valid return true if there are no errors , otherwise false
func (f *Form) Valid () bool {
	return len(f.Errors) == 0 
}

// MinLength check for minimum string length
func(f *Form) MinLength (field string , length int, r *http.Request) bool {
	x := r.Form.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("%s must at least %d characters long", field, length))
		return false
	}
	return true
}