package main

import (
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"
)

func (app *application) render(w http.ResponseWriter, r *http.Request, data *templateData, file ...string) {
	tmp, err := template.ParseFiles(file...)
	if err != nil {
		app.errorLog.Fatal(err)
	}

	err = tmp.Execute(w, data)
	if err != nil {
		app.errorLog.Fatal(err)
	}
}

func (app *application) addImage(file multipart.File, name *multipart.FileHeader) (string, error) {

	// directory saat ini
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// variable untuk menyimpan nama image beserta ekstensinya
	fileName := name.Filename

	// directory untuk menyimpan semua gambar
	fileLocation := filepath.Join(dir, "views/static/images", fileName)

	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}

	defer targetFile.Close()

	// isi directory yang telah dibuka d targetFile dengan file yang telah di upload
	_, err = io.Copy(targetFile, file)
	if err != nil {
		return "", nil
	}

	return fileName, nil

}

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9]))")

func validateSignup(name, email, pass string, min int, pattern *regexp.Regexp) map[string]string {
	errors := make(map[string]string)

	// validasi form blank
	if strings.TrimSpace(name) == "" {
		errors["name"] = "This field cannot be blank"
	}

	if strings.TrimSpace(email) == "" {
		errors["email"] = "This field cannot be blank"
	} else if !pattern.MatchString(email) {
		errors["email"] = "This field is invalid"
	}

	if strings.TrimSpace(pass) == "" {
		errors["password"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(pass) < min {
		err := fmt.Sprintf("This field is too short (minimum is %d characters)", min)
		errors["password"] = err
	}

	return errors
}

func (app *application) isAuthenticated(r *http.Request) bool {
	return app.session.Exists(r, "authenticatedUserID")
}
