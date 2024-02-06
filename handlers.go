package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/afifallimaruf/go-booksell/forms"
	"github.com/afifallimaruf/go-booksell/models"
	"github.com/afifallimaruf/go-booksell/models/mysql"
)

func (app *application) indexHandler(w http.ResponseWriter, r *http.Request) {
	books, err := mysql.GetBooks()
	if err != nil {
		app.errorLog.Fatal(err)
	}

	flash := app.session.PopString(r, "flash")

	app.render(w, r, &templateData{
		Flash:           flash,
		Books:           books,
		IsAuthenticated: app.isAuthenticated(r),
	}, "views/html/index.html", "views/html/base.html")
}

func (app *application) booksHandler(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, &templateData{
		IsAuthenticated: app.isAuthenticated(r),
	}, "views/html/books.html", "views/html/base.html")
}

func (app *application) addChart(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, &templateData{
		IsAuthenticated: app.isAuthenticated(r),
	}, "views/html/chart.html", "views/html/base.html")
}

func (app *application) aboutHandler(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, &templateData{
		IsAuthenticated: app.isAuthenticated(r),
	}, "views/html/about.html", "views/html/base.html")
}

func (app *application) nonFictionHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/category/non-fiction" {
		app.render(w, r, &templateData{}, "views/html/non_fiction.html", "views/html/base.html")
	}
}

func (app *application) signupPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, &templateData{
		Form: forms.New(nil),
	}, "views/html/signup.html", "views/html/base.html")
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.errorLog.Fatal(err)
	}

	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MaxLength("password", 5)
	form.MatchesPattern("email", forms.EmailRX)
	fmt.Println(form)

	if !form.Valid() {
		app.render(w, r, &templateData{
			Form: form,
		}, "views/html/signup.html", "views/html/base.html")

		return
	}

	if err := mysql.Insert(form.Get("name"), form.Get("email"), form.Get("password")); err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "Address is already in use")
			app.render(w, r, &templateData{
				Form: form,
			}, "views/html/signup.html", "views/html/base.html")
		} else {
			app.errorLog.Fatal(err)
		}

		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) loginPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, &templateData{
		Form: forms.New(nil),
	}, "views/html/login.html", "views/html/base.html")
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.errorLog.Fatal(err)
	}

	form := forms.New(r.PostForm)

	id, err := mysql.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or Password is incorrect")
			app.render(w, r, &templateData{
				Form: form,
			}, "views/html/login.html", "views/html/base.html")
		} else {
			app.errorLog.Fatal(err)
		}

		return
	}

	app.session.Put(r, "authenticatedUserID", id)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedUserID")

	app.session.Put(r, "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) addBooksForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, &templateData{
		IsAuthenticated: app.isAuthenticated(r),
	}, "views/html/add_books.html", "views/html/base.html")
}

func (app *application) addBooks(w http.ResponseWriter, r *http.Request) {
	// parse form files
	// statement r.ParseMultipartForm digunakan utuk parsing form data yang dikirm
	if err := r.ParseMultipartForm(1024); err != nil {
		app.errorLog.Fatal(http.StatusInternalServerError)
		app.errorLog.Fatal(err)
	}

	// FormFile digunakan untuk mengambil image yang di upload
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		app.errorLog.Fatal(err)
	}

	defer file.Close()

	// memanggil datas
	title := r.PostForm.Get("title")
	author := r.PostForm.Get("author")
	summary := r.PostForm.Get("summary")
	price := r.PostForm.Get("price")

	// lempar nama image ke fungsi untuk memasukan image kedalam directory
	imgName, err := app.addImage(file, fileHeader)
	if err != nil {
		app.errorLog.Fatal(err)
	}

	if ok := mysql.InsertBook(title, author, summary, price, imgName); !ok {
		app.errorLog.Fatal("ERROR")
	}

	app.session.Put(r, "flash", "Book Successfully Added!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
