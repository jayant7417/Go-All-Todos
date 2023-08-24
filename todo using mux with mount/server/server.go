package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/mux"
	"jayant/handler"
	"jayant/utils"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	chi.Router
	server *http.Server
}

const (
	readTimeout       = 5 * time.Minute
	readHeaderTimeout = 30 * time.Second
	writeTimeout      = 5 * time.Minute
)

func SetupRouter() http.Handler {
	router := mux.NewRouter()
	mount(router, "/todo", Router())
	return router
}

func Router() http.Handler {
	router := mux.NewRouter()
	mount(router, "/health", Health())
	mount(router, "/public", Public())
	mount(router, "/user", User())
	mount(router, "/task", Task())
	return router
}

func Health() http.Handler {
	r := mux.NewRouter()
	r.Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"status": "server is running",
		})
	})
	return r
}

func Public() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/register", handler.CreateUser).Methods("POST")
	r.HandleFunc("/login", handler.LoginUser).Methods("POST")
	return r
}

func User() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/", handler.UserInfo).Methods("GET")
	r.HandleFunc("/", handler.DeleteUser).Methods("DELETE")
	r.HandleFunc("/logout", handler.Logout).Methods("DELETE")
	r.HandleFunc("/", handler.UpdateUser).Methods("PUT")
	return r
}

func Task() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/", handler.CreateTask).Methods("POST")
	r.HandleFunc("/", handler.AllTask).Methods("GET")
	r.HandleFunc("/{taskId:[0-9]+}", handler.DeleteTask).Methods("DELETE")
	r.HandleFunc("/", handler.UpdateTask).Methods("PUT")
	r.HandleFunc("/{taskId:[0-9]+/complete", handler.Complete).Methods("PUT")
	return r
}

//func SetupRouter() *Server {
//	router := chi.NewRouter()
//	router.Route("/todo", func(v1 chi.Router) {
//		v1.Get("/health", func(w http.ResponseWriter, r *http.Request) {
//			utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
//				"status": "server is running",
//			})
//		})
//		v1.Route("/public", func(p chi.Router) {
//			p.Post("/register", handler.CreateUser)
//			p.Post("/login", handler.LoginUser)
//		})
//		v1.Route("/user", func(u chi.Router) {
//			u.Use(middlewares.AuthMiddleware)
//			u.Get("/", handler.UserInfo)
//			u.Delete("/", handler.DeleteUser)
//			u.Delete("/logout", handler.Logout)
//			u.Put("/", handler.UpdateUser)
//		})
//
//		v1.Route("/task", func(t chi.Router) {
//			t.Use(middlewares.AuthMiddleware)
//			t.Post("/", handler.CreateTask)
//			t.Get("/", handler.AllTask)
//			t.Delete("/{taskId}", handler.DeleteTask)
//			t.Put("/", handler.UpdateTask)
//			t.Put("/{taskId}/complete", handler.Complete)
//		})
//	})
//	return &Server{
//		Router: router,
//	}
//}

func mount(r *mux.Router, path string, handler http.Handler) {
	r.PathPrefix(path).Handler(
		http.StripPrefix(
			strings.TrimSuffix(path, "/"),
			handler,
		),
	)
}

func (svc *Server) Run(port string) error {
	svc.server = &http.Server{
		Addr:              port,
		Handler:           svc.Router,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
	}
	return svc.server.ListenAndServe()
}

func (svc *Server) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return svc.server.Shutdown(ctx)
}
