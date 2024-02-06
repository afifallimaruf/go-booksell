package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

type Form struct {
	url.Values
	Errors errors
}

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9]))")

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (form *Form) Required(fields ...string) {
	for _, field := range fields {
		value := form.Get(field)
		if strings.TrimSpace(value) == "" {
			form.Errors.Add(field, "This field cannot be blank")
		}
	}
}

func (form *Form) MaxLength(field string, max int) {
	key := form.Get(field)
	if field == "" {
		return
	}

	if utf8.RuneCountInString(key) > max {
		form.Errors.Add(field, fmt.Sprintf("This field is too long (max is %d characters)", max))
	}
}

func (form *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := form.Get(field)
	if value == "" {
		return
	}

	if !pattern.MatchString(value) {
		form.Errors.Add(field, "This field is invalid")
	}
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
