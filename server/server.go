package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"main.go/handlers"
	"net/http"
)

type Server struct {
	Router *mux.Router
	Addr   string
}

func SetUpRoutes(addr string) *Server {
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/make", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"message": "Hello World!"})
	})
	api.HandleFunc("/login", handlers.UserLogin).Methods("POST")
	return &Server{Router: r, Addr: addr}
}
func (s *Server) Start() {
	err := http.ListenAndServe(s.Addr, s.Router)
	if err != nil {
		logrus.Fatal(err)
		return
	}
}
