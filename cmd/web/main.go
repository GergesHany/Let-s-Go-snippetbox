package main

import (
	"database/sql" 
	"log"
	"net/http"
	"flag" 
	"os"

    "snippetbox.alexedwards.net/internal/models"

	_ "github.com/go-sql-driver/mysql" // sql driver(mysql)

	"html/template"
)

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger

	snippets *models.SnippetModel
	templateCache map[string]*template.Template
}

func main()  {

	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		snippets: &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
	
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// db.Ping() checks if the database connection is correctly set up
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}