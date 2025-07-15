package server

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"main.go/handlers"
	"main.go/middlewares"
	"main.go/models"
	"net/http"
)

type Server struct {
	Router *mux.Router
	Addr   string
}

func SetUpRoutes(addr string) *Server {
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1/users").Subrouter()

	//public routes
	api.HandleFunc("/login", handlers.UserLogin).Methods("POST")
	api.HandleFunc("/register", handlers.UserRegister).Methods("POST")
	api.HandleFunc("employeeManager/login", handlers.EmployeeManagerLogin).Methods("POST")
	api.HandleFunc("/admin/login", handlers.AdminLogin).Methods("POST")
	api.HandleFunc("/assetManager/login", handlers.AssetManagerLogin).Methods("POST")

	// admin routes
	Protected := api.PathPrefix("/").Subrouter()
	Protected.Use(middlewares.AuthMiddleware)
	Protected.HandleFunc("/admin/updateRole", handlers.UpdateUserRole).Methods("PUT")

	//asset Manager routes
	sharedRoutes := Protected.PathPrefix("/shared").Subrouter()
	sharedRoutes.Use(middlewares.RequireRole(string(models.Admin), string(models.AssetManager)))
	sharedRoutes.HandleFunc("/createAsset", handlers.CreateAsset).Methods("POST")
	sharedRoutes.HandleFunc("/assignAsset", handlers.AssignAsset).Methods("POST")
	sharedRoutes.HandleFunc("/getAsset", handlers.GetAsset).Methods("GET")
	sharedRoutes.HandleFunc("/retrieveAsset", handlers.RetrieveEmployeeAsset).Methods("PUT")
	sharedRoutes.HandleFunc("/updateAsset", handlers.UpdateUserAsset).Methods("PUT")
	sharedRoutes.HandleFunc("/sendToService", handlers.SendToService).Methods("POST")
	sharedRoutes.HandleFunc("/assetTimeline", handlers.GetAssetTimeline).Methods("GET")

	//employee Manager routes
	EmployeeRoutes := Protected.PathPrefix("/employeeManager").Subrouter()
	EmployeeRoutes.Use(middlewares.RequireRole(string(models.Admin), string(models.EmployeeManager)))
	EmployeeRoutes.HandleFunc("/updateEmployeeInfo", handlers.UpdateEmployeeInfo).Methods("PUT")
	EmployeeRoutes.HandleFunc("/getEmployeeInfo", handlers.GetEmployeeInfo).Methods("GET")
	EmployeeRoutes.HandleFunc("/getUserAssetTimeline", handlers.UserAssetTimeline).Methods("GET")
	EmployeeRoutes.HandleFunc("/deleteUser", handlers.DeleteEmployeeInfo).Methods("PUT")
	return &Server{Router: r, Addr: addr}
}
func (s *Server) Start() {
	err := http.ListenAndServe(s.Addr, s.Router)
	if err != nil {
		logrus.Fatal(err)
		return
	}
}
