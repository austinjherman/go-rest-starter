package main

import (
	"aherman/src/container"
	oauthControllers "aherman/src/controllers/oauth"
	userControllers "aherman/src/controllers/user"
	"aherman/src/database"
	errorFacades "aherman/src/facades/error"
	tokenFacades "aherman/src/facades/token"
	userFacades "aherman/src/facades/user"
	"aherman/src/middleware"
	tokenModels "aherman/src/models/token"
	userModels "aherman/src/models/user"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/ttacon/chalk"
)

func main() {

	// load .env file
	// todo: error handling
  err := godotenv.Load(".env")
  if err != nil {
    log.Printf("%swarn%s no environment file found\n", chalk.Yellow, chalk.Reset)
  } else {
		dir, _ := os.Getwd()
		log.Printf("%sinfo%s loading env file %s/.env\n", chalk.Blue, chalk.Reset, dir)
	}

	// todo: error handling
	db := database.Init()

	// create a container for dependency injection
	dependencies := &container.Container{
		Current: &container.CurrentContainer{
			Token: &tokenModels.Token{},
			User: &userModels.User{},
		},
		Facades: &container.Facades{
			Error: &errorFacades.Error{},
			Token: &tokenFacades.Token{
				DB: db,
			},
			User: &userFacades.User{
				DB: db,
			},
		},
	}

	r := gin.Default()

	// Abort 204 options request
	r.Use(
		func(c *gin.Context) {
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
		},
		func(c *gin.Context) {
			dependencies.Facades.Error.ServerWriter = gin.DefaultErrorWriter
		},
	)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/oauth/token", oauthControllers.Router(dependencies))
	r.GET("/oauth/refresh", oauthControllers.RefreshForAccess(dependencies))

	r.GET("/logout", middleware.Token(dependencies), userControllers.Logout(dependencies))
	r.GET("/logout-all", middleware.Token(dependencies), userControllers.LogoutAll(dependencies))

	r.POST("/user", middleware.Invitation(dependencies), userControllers.Create(dependencies))
	r.GET("/user", middleware.Token(dependencies), userControllers.Read(dependencies))
	r.PUT("/user", middleware.Token(dependencies), userControllers.Update(dependencies))
	r.DELETE("/user", middleware.Token(dependencies), userControllers.Delete(dependencies))

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}