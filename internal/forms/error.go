package forms

type errors map[string][]string

// Add adds error message for give error field
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get return first field error message
func (e errors) Get(field string) string {
	errs := e[field]
	if len(errs) == 0 {
		return ""
	}

	return errs[0]
}
