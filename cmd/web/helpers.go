package main
import (
	"fmt"
	"net/http"
	"runtime/debug"
   "bytes"
   "time"
   "errors"

   "github.com/go-playground/form/v4" 
   "github.com/justinas/nosurf"
)

func (app *application) isAuthenticated(r *http.Request) bool {
   isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
   // If the isAuthenticatedContextKey key does not exist in the request context
   if !ok {
     return false
   }
   return isAuthenticated
}


func (app *application) newTemplateData(r *http.Request) *templateData {
   csrfToken := nosurf.Token(r)
   // app.infoLog.Printf("CSRF Token: %s", csrfToken) // Log the CSRF token

   return &templateData{
      CurrentYear:     time.Now().Year(),
      Flash:           app.sessionManager.PopString(r.Context(), "flash"),
      IsAuthenticated: app.isAuthenticated(r),
      CSRFToken:       csrfToken, // Add the CSRF token to the template data.
   }
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
   
   // dst: destination struct

   err := r.ParseForm()
   if err != nil {
      return err
   }
   
   // Use the decoder to decode the form data into the struct.
   err = app.formDecoder.Decode(dst, r.PostForm)

   if err != nil {
      var invalidDecoderError *form.InvalidDecoderError

      // error.As() is used to check if the error returned by Decode is an InvalidDecoderError.
      if errors.As(err, &invalidDecoderError) {
         panic(err) // This will be caught by the serverError helper.
      }
      return err
   }
   
   return nil
}

func (app *application) serverError(w http.ResponseWriter, err error) {
   // A stack trace is a report of the active stack frames at 
   // a certain point in time during the execution of a program.
   trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack()) 

   // The Output method of a log.Logger instance is used to write a log entry
   // to the log file. The first argument is the log level, which is set to 2 to indicate that this 
   // message should be written as an error message. The second argument is the message to write.
   app.errorLog.Output(2, trace) 


   if app.debug {
      http.Error(w, trace, http.StatusInternalServerError)
      return
   }

   http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
   http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
   app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
  
   ts, ok := app.templateCache[page]
   if !ok {
      err := fmt.Errorf("the template %s does not exist", page)
      app.serverError(w, err)
      return
   }

   // Initialize a new buffer.
   buf := new(bytes.Buffer)
   
   // Write the template to the buffer instead of directly to the http.ResponseWriter.
   err := ts.ExecuteTemplate(buf, "base", data)
   if err != nil {
      app.serverError(w, err)
      return
   }
   
   w.WriteHeader(status)
   buf.WriteTo(w) // Write the contents of the buffer to the http.ResponseWriter.
}
