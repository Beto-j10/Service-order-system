package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"harper/api/login"
	db "harper/database"
	"harper/lib/graphql"
)

func main() {

	db.Init()

	port := ":8080"
	Router := gin.Default()

	//Config CORS
	Router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"},
		// AllowHeaders: []string{"Content-Type", "Authorization"},
		AllowMethods: []string{"POST", "GET"},
	}))

	Router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	login.Login(Router)
	graphql.Service(Router)

	log.Println("Harper started and listening")
	log.Fatal(Router.Run(port))

}
