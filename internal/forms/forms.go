package forms

import (
	"net/http"
	"net/url")


type Form struct {
	url.Values
	Errors errors
}

// Valid returns boolean if there are errors
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// New initialize a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Has(field string, r *http.Request) bool {
	formField := r.Form.Get(field)
	if formField == "" {
		f.Errors.Add(field, "This field cannot be blank")
		return false
	}
	return true
}