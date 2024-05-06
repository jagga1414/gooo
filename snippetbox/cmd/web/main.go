package main

import (
	// "log"
	"log/slog"
	"net/http"
	"os"

	// "strconv"
	"flag"
	// "fmt"
)

type application struct {
	logger *slog.Logger 
}

func main() {

	addr := flag.String("addr", ":4000", "set your web host")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := &application{
		logger: logger, 
	}

	app.logger.Info("starting server on","port", *addr)
	err := http.ListenAndServe(*addr, app.routes())
	app.logger.Error(err.Error())
	os.Exit(1)
}
