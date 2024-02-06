package forms

type errors map[string][]string

func (e errors) Add(field string, message string) {
	e[field] = append(e[field], message)
}

func (e errors) Get(field string) string {
	el := e[field]

	if len(el) == 0 {
		return ""
	}

	return el[0]
}
