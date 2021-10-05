package main

import (
	"aherman/src/container"
	oauthControllers "aherman/src/controllers/oauth"
	userControllers "aherman/src/controllers/user"
	"aherman/src/database"
	tokenFacades "aherman/src/facades/token"
	userFacades "aherman/src/facades/user"
	"aherman/src/middleware"
	tokenModels "aherman/src/models/token"
	userModels "aherman/src/models/user"
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
		Current: &container.CurrentContainer{
			Token: &tokenModels.Token{},
			User: &userModels.User{},
		},
		Token: &tokenFacades.Token{DB: db},
		User: &userFacades.User{DB: db},
	}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/oauth/token", oauthControllers.Router(dependencies))

	r.GET("/logout", middleware.Token(dependencies), userControllers.Logout(dependencies))

	r.POST("/user", middleware.Invitation(), userControllers.Create(dependencies))
	r.GET("/user", middleware.Token(dependencies), userControllers.Read(dependencies))
	r.PUT("/user", middleware.Token(dependencies), userControllers.Update(dependencies))
	r.DELETE("/user", middleware.Token(dependencies), userControllers.Delete(dependencies))

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}