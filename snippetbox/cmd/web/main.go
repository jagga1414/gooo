package main

import (
	// "log"
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"snippetbox.jagdish.net/internal/models"
	// "fmt"
)

type application struct {
	logger *slog.Logger 
	snippets *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder *form.Decoder
}

func main() {

	addr := flag.String("addr", ":4000", "set your web host")
	dsn := flag.String("dsn", "web:123456@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	db, err := openDB(*dsn)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	templateCache, err := newTemplateCache() 
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1) 
	}
	form := form.NewDecoder()
	app := &application{
		logger: logger, 
		snippets: &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder: form,
	}
	defer db.Close()
	app.logger.Info("starting server on","port", *addr)
	err = http.ListenAndServe(*addr, app.routes())
	app.logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn) 
	if err != nil {
		return nil, err
	}
	err = db.Ping() 
	if err != nil {
		db.Close()
		return nil, err 
	}
	return db, nil
}