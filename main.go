package main

import (
	"log"

	"example.com/hello/models"
	"example.com/hello/routes"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {
	db := InitDB()
	// Automatically migrate your schema
	err := db.AutoMigrate(&models.Submissions{})
	if err != nil {
		log.Fatalln(err)
	}

	router := gin.Default()
	router.Use(cors.New(cors.Config{
        AllowOrigins: []string{"http://localhost:3000"},
        AllowMethods: []string{"POST", "PUT", "PATCH", "DELETE"},
        AllowHeaders: []string{"Content-Type,access-control-allow-origin, access-control-allow-headers"},
    }))
	routes.InitRoutes(router,db)
	router.Run()
}
