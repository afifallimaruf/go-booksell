package main

import (
	"github.com/afifallimaruf/go-booksell/forms"
	"github.com/afifallimaruf/go-booksell/models"
)

type templateData struct {
	Flash           string
	IsAuthenticated bool
	Form            *forms.Form
	Books           []*models.Books
}
