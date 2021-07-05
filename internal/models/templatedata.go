package models

import "github.com/johnreybacal/go-book/internal/forms"

type TemplateData struct {
	StringMap map[string]string
	IntMapmap map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
	Form      *forms.Form
}
