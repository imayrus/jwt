package main

import (
	"github.com/imayrus/jwt-api/routes"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// var (
// 	router    *mux.Router
// 	secretkey string = "jwtsecretkey"
// )

func main() {

	routes.CreateRouter()
	routes.InitializeRoute()
	routes.ServerStart()
}
