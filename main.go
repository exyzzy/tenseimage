package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/exyzzy/tenseimage/match"
)

type config struct {
	port string
	env  string
}

type application struct {
	config config
	logger *log.Logger
}

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/match", app.matchHandler)
	mux.HandleFunc("/", app.homeHandler)
	return mux
}

func main() {
	var cfg config

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	flag.StringVar(&cfg.port, "port", port, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		config: cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err := srv.ListenAndServe()
	logger.Fatal(err)
}

func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet: //just print available api calls
		fmt.Fprintf(w, "environment: %s\n", app.config.env)
		fmt.Fprintln(w, "you can use:")
		fmt.Fprintf(w, "\tcurl -X OPTIONS %s/match -i\n", r.Host)
		fmt.Fprintf(w, "\tcurl -d '{\"Url\":\"https://www.ndow.org/wp-content/uploads/2021/10/neovison_vison-992x679.jpg\"}' -H \"Content-Type: application/json\" -X POST %s/match\n", r.Host)

	case http.MethodOptions:
		w.Header().Set("Allow", "GET, OPTIONS")
		w.WriteHeader(http.StatusNoContent)

	default:
		w.Header().Set("Allow", "GET, OPTIONS")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

type TenseImage struct {
	Url string
}

func (app *application) matchHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		var tenseImage TenseImage
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&tenseImage); err != nil {
			http.Error(w, "invalid request payload", http.StatusBadRequest)
		}
		defer r.Body.Close()
		fmt.Fprintln(w, tenseImage)

		dir := "./model"
		best := match.Match(&dir, &tenseImage.Url, true)
		fmt.Fprintln(w, best)

	case http.MethodOptions:
		w.Header().Set("Allow", "POST, OPTIONS")
		w.WriteHeader(http.StatusNoContent)

	default:
		w.Header().Set("Allow", "POST, OPTIONS")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
