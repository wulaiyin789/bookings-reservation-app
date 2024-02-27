package models

import "github.com/tsawler/bookings-app/internal/forms"

// TemplateData holds data sent from handlers to templates
// IsAuth: Greater than 0, then login. Otherwise, unauthenticated
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
	Form      *forms.Form
	IsAuth    int
}
