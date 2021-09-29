package main

import (
	oauthControllers "aherman/src/controllers/oauth"
	userControllers "aherman/src/controllers/user"
	"aherman/src/database"
	"aherman/src/middleware"
	clientModels "aherman/src/models/client"
	oauthModels "aherman/src/models/oauth"
	userModels "aherman/src/models/user"
	"aherman/src/types/container"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// load .env file
	// todo: error handling
  err := godotenv.Load(".env")
  if err != nil {
    log.Fatalf("Error loading .env file")
  }

	// todo: error handling
	db := database.Init()

	// create a container for dependency injection
	dependencies := &container.Container{
		Client: &clientModels.Env{DB: db},
		OAuth: &oauthModels.Env{DB: db},
		User: &userModels.Env{DB: db},
	}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/oauth/token", oauthControllers.PasswordGrant(dependencies))

	r.POST("/users/create", middleware.Invitation(), userControllers.Create(dependencies))

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}