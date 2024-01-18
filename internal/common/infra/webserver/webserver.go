package webserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]map[string]http.HandlerFunc // [path][method]handler
	WebServerPort string
	Middlewares   []func(h http.Handler) http.Handler
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]map[string]http.HandlerFunc),
		WebServerPort: serverPort,
	}
}

func (s *WebServer) Start() {
	s.loadMiddlewares()
	s.loadHandlers()

	server := &http.Server{
		Addr:         ":" + s.WebServerPort,
		Handler:      s.Router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	fmt.Println("Server running on port " + s.WebServerPort)
	if err := server.ListenAndServe(); err != nil {
		panic(err.Error())
	}
}

func (s *WebServer) AddMiddleware(middleware func(h http.Handler) http.Handler) {
	s.Middlewares = append(s.Middlewares, middleware)
}

func (s *WebServer) Get(path string, handler http.HandlerFunc) {
	s.addHandler(path, "GET", handler)
}

func (s *WebServer) Post(path string, handler http.HandlerFunc) {
	s.addHandler(path, "POST", handler)
}

func (s *WebServer) Options(path string, handler http.HandlerFunc) {
	s.addHandler(path, "OPTIONS", handler)
}

func (s *WebServer) Put(path string, handler http.HandlerFunc) {
	s.addHandler(path, "PUT", handler)
}

func (s *WebServer) Delete(path string, handler http.HandlerFunc) {
	s.addHandler(path, "DELETE", handler)
}

func (s *WebServer) Patch(path string, handler http.HandlerFunc) {
	s.addHandler(path, "PATCH", handler)
}

func (s *WebServer) addHandler(path string, method string, handler http.HandlerFunc) {
	if _, ok := s.Handlers[path]; !ok {
		s.Handlers[path] = make(map[string]http.HandlerFunc)
	}
	s.Handlers[path][method] = handler
}

func (s *WebServer) loadMiddlewares() {
	s.Router.Use(middleware.Logger)
	s.Router.Use(cors.Handler(cors.Options{}))
	s.Router.Use(middleware.StripSlashes)
	s.Router.Use(middleware.Recoverer)
	for _, middleware := range s.Middlewares {
		s.Router.Use(middleware)
	}
}

func (s *WebServer) loadHandlers() {
	methods := map[string]interface{}{
		"GET":     s.Router.Get,
		"POST":    s.Router.Post,
		"PUT":     s.Router.Put,
		"DELETE":  s.Router.Delete,
		"OPTIONS": s.Router.Options,
		"PATCH":   s.Router.Patch,
	}
	for path, handler := range s.Handlers {
		for method, handler := range handler {
			methods[method].(func(string, http.HandlerFunc))(path, handler)
		}
	}
}
