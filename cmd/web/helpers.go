package main
import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
    
   // A stack trace is a report of the active stack frames at 
   // a certain point in time during the execution of a program.
   trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack()) 

   // The Output method of a log.Logger instance is used to write a log entry
   // to the log file. The first argument is the log level, which is set to 2 to indicate that this 
   // message should be written as an error message. The second argument is the message to write.
   app.errorLog.Output(2, trace) 
   http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
   http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
   app.clientError(w, http.StatusNotFound)
}