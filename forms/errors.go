package forms

type errors map[string][]string

// Add allows the addition of error message for a given form field during form validation
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get is for retrieving first error message associated with field during form validation
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}

	return es[0]
}
