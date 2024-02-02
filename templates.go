package main

import (
	"net/url"

	"github.com/afifallimaruf/go-booksell/models"
)

type templateData struct {
	Flash           string
	FormData        url.Values
	IsAuthenticated bool
	FormErrors      map[string]string
	Books           []*models.Books
}
