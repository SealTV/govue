package server

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/middleware"

	"github.com/SealTV/govue/repository"
	"github.com/go-chi/chi"
)

type server struct {
	ar repository.AccountRepository
	r  *chi.Mux
}

// NewHTTPHandler - create new web request handler
func NewHTTPHandler(ar repository.AccountRepository) http.Handler {
	r := chi.NewRouter()
	s := server{
		r:  r,
		ar: ar,
	}

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "static")
	log.Println(filesDir)
	fileServer(r, "/static", http.Dir(filesDir))

	fs := http.FileServer(http.Dir("static/"))
	r.Handle("/static/", http.StripPrefix("/static/", fs))

	r.Get("/", s.getHelloWorld)
	r.Get("/account", s.getAccountHandler)
	r.NotFound(s.notFound)
	return &s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.r.ServeHTTP(w, r)
}

func (s *server) getHelloWorld(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	log.Println(w.Write([]byte("Hello world")))
}

func (s *server) getAccountHandler(w http.ResponseWriter, r *http.Request) {

}

func (s *server) notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Sorry, page no found :("))
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
