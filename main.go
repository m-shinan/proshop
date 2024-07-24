package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/m-shinan/project-shop/database"
	"github.com/m-shinan/project-shop/routes"
)

var (
	R *gin.Engine
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	database.DbConnect()

	if err := database.CreateTables(); err != nil {
		log.Fatal("Error creating tables:", err)
	}

	R = gin.Default()

	R.Static("/uploads", "./uploads")

	R.LoadHTMLGlob("templates/*")

}
func main() {

	routes.AdminRoutes(R)
	routes.UserRoutes(R)

	if err := R.Run(":8080"); err != nil {
		log.Fatal("Failed to run server:", err)
	}

}
