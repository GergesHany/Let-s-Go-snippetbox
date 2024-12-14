package main 

import (
	"net/http"
)

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