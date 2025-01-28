package main 

import (
    "context"
	"net/http"
	"fmt"
	"github.com/justinas/nosurf"
)

func (app *application) authenticate(next http.Handler) http.Handler {
   return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	  // "authenticatedUserID" value is in the session
	  id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
      if id == 0 {
		  next.ServeHTTP(w, r)
		  return
	  }

	  // check if the user exists
	  exists, err := app.users.Exists(id)
	  if err != nil {
		  app.serverError(w, err)
		  return
	  }

      // If the user exists, add it to the request context
	  if exists {
		ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
		r = r.WithContext(ctx)
	  } 

	  next.ServeHTTP(w, r)
   })
}

func noSurf(next http.Handler) http.Handler {
	// Create a new CSRF handler using nosurf, passing in the next handler in the chain
	csrfHandler := nosurf.New(next)

	// Set the Secure flag on the CSRF cookie
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path: "/",
		Secure: true,
    })
	return csrfHandler
}

// secureHeaders -> servemux -> application handler
// flow of control actually looks like: secureHeaders → servemux → application handler → servemux → secureHeaders

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		w.Header().Set("Content-Security-Policy", 
		               "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

        w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
	
		next.ServeHTTP(w, r)
	})
}

// logRequest ↔ secureHeaders ↔ servemux ↔ application handler

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

// recoverPanic ↔ logRequest ↔ secureHeaders ↔ servemux ↔ application handler
func (app * application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			app.sessionManager.Put(r.Context(), "redirectPathAfterLogin", r.URL.Path)
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		// "Cache-Control: no-store" header so that pages
        // require authentication are not stored in the users browser cache 
		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}