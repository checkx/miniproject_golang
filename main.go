package main

import (
	"rest-echo-gorm/config"
	"rest-echo-gorm/routes"
)

func main() {
	// Initialize Database
	config.InitialMigration()

	// Initialize Routes
	e := routes.New()

	e.Logger.Fatal(e.Start("localhost:9000"))
}
