package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv" // for converting string to int
	"strings"
	"unicode/utf8"

	"snippetbox.alexedwards.net/internal/models"

	"github.com/julienschmidt/httprouter"

)

type snippetCreateForm struct {
	Title string 
	Content string
	Expires int
	FieldErrors map[string]string
}

func (app *application) home(w http.ResponseWriter, r *http.Request){

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r) 
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	
   params := httprouter.ParamsFromContext(r.Context());

   id, err := strconv.Atoi(params.ByName("id"))
   if err != nil || id < 1 {
	  app.notFound(w)
	  return
   }

   snippets, err := app.snippets.Get(id)
   if err != nil {
	   if errors.Is(err, models.ErrNoRecord) {
   		   app.notFound(w)
	   } else {
		   app.serverError(w, err)
	   }
	   return
   }

   data := app.newTemplateData(r)
   data.Snippet = snippets

   app.render(w, http.StatusOK, "view.tmpl.html", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
   data := app.newTemplateData(r)
   
   data.Form = snippetCreateForm{
	 Expires: 365,
   }

   app.render(w, http.StatusOK, "create.tmpl.html", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request){

    err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := snippetCreateForm{
		Title: title,
		Content: content,
		Expires: expires,
		FieldErrors: make(map[string]string),
	}

	// Validate the title field.
	if strings.TrimSpace(title) == "" {
		form.FieldErrors["title"] = "Title cannot be blank"
	}else if utf8.RuneCountInString(title) > 100 {
		form.FieldErrors["title"] = "Title cannot be longer than 100 characters"
	}

    // Validate the content field.
	if strings.TrimSpace(content) == "" {
		form.FieldErrors["content"] = "Content cannot be blank"
	}

	// validate the expires field.
	if expires != 1 && expires != 7 && expires != 365 {
		form.FieldErrors["expires"] = "Please select a valid expiration time"
	}

	// If there are any errors, redisplay the form.
	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		// StatusUnprocessableEntity: 422
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}