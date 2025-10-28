package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"reader/internal/config"
	database "reader/internal/db"
	"reader/internal/dictionary"
	"reader/internal/nlp"
	"reader/internal/reader"
	readerHandler "reader/internal/reader/handler"
	"reader/internal/users"

	"reader/static"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
)

func main() {
	envFile := flag.String("env", "./.env", "path to environment file")
	flag.Parse()

	cfg := config.Get(*envFile)
	db := database.Connect(cfg.Db)

	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)

	mux := http.NewServeMux()

	staticFileServer := http.FileServer(http.FS(static.StaticFiles))
	mux.Handle("GET /static/", http.StripPrefix("/static", staticFileServer))

	usersModel := users.NewModel(db)
	usersService := users.NewService(cfg.Oauth, sessionManager, &usersModel)
	usersService.Register(mux)
	injectUserMiddleware := users.NewInjectUserMiddleware(&usersModel, sessionManager)

	nlpClient := nlp.NewClient(cfg.NlpUrl)
	dictionaryClient := dictionary.NewClient()

	textModel := reader.NewTextModel(db)
	wordModel := reader.NewWordModel(db)
	readerService := reader.NewService(&textModel, &wordModel, &nlpClient, &dictionaryClient)
	readerServiceHandler := readerHandler.New(&readerService)
	readerServiceHandler.Register(mux)

	server := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      sessionManager.LoadAndSave(injectUserMiddleware.Wrap(mux)),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	if cfg.Secure {
		fmt.Printf("Listening on https://127.0.0.1:%d\n", cfg.Port)
		err := server.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			panic(err)
		}
		return
	}

	fmt.Printf("Listening on http://127.0.0.1:%d\n", cfg.Port)
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
