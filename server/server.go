package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"main.go/handlers"
	"main.go/middlewares"
	"net/http"
)

type Server struct {
	Router *mux.Router
	Addr   string
}

func SetUpRoutes(addr string) *Server {
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1/users").Subrouter()
	api.HandleFunc("/make", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"message": "Hello World!"})
	})
	api.HandleFunc("/login", handlers.UserLogin).Methods("POST")
	api.HandleFunc("/register", handlers.UserRegister).Methods("POST")

	// admin routes
	api.HandleFunc("/admin/login", handlers.AdminLogin).Methods("POST")
	Protected := api.PathPrefix("/").Subrouter()
	Protected.Use(middlewares.AuthMiddleware)
	Protected.HandleFunc("/admin/updateRole", handlers.UpdateUserRole).Methods("PUT")

	//asset Manager routes
	api.HandleFunc("/assetManager/login", handlers.AssetManagerLogin).Methods("POST")
	sharedRoutes := Protected.PathPrefix("/shared").Subrouter()
	sharedRoutes.Use(middlewares.RequireRole("admin", "asset_manager"))
	sharedRoutes.HandleFunc("/createAsset", handlers.CreateAsset).Methods("POST")
	sharedRoutes.HandleFunc("/assignAsset", handlers.AssignAsset).Methods("POST")

	sharedRoutes.HandleFunc("/getAsset", handlers.GetAsset).Methods("GET")

	//employee Manager routes
	api.HandleFunc("/employeeManager/login", handlers.EmployeeManagerLogin).Methods("POST")

	return &Server{Router: r, Addr: addr}
}
func (s *Server) Start() {
	err := http.ListenAndServe(s.Addr, s.Router)
	if err != nil {
		logrus.Fatal(err)
		return
	}
}
