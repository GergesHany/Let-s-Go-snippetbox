package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv" // for converting string to int

	"snippetbox.alexedwards.net/internal/models"

	"github.com/julienschmidt/httprouter"

)

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
  w.Write([]byte("Display the form for creating a new snippet..."))
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request){
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n-Kobayashi Issa"
	expires := 7

	// Insert the snippet into the database
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
	   app.serverError(w, err)
	   return
	}

	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}