package client

import (
	"aherman/src/http/response"
	clientModels "aherman/src/models/client"
	"aherman/src/types/container"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

// Create a user and persist to database
func Create(app *container.Container) gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			client *clientModels.Client = &clientModels.Client{}
			clientCreatable *clientModels.Createable = &clientModels.Createable{}
			clientPublic *clientModels.Public = &clientModels.Public{}
		)

		// bind input to client variable.
		// if we can't, there was a validation error.
		if err := c.ShouldBindBodyWith(clientCreatable, binding.JSON); err != nil {
			res := response.ErrValidation
			c.Error(res)
			c.JSON(response.Error(res))
			return
		}

		client.ID = uuid.New()
		client.Name = clientCreatable.Name

		if result := app.Client.DB.Create(client); result.Error != nil {
			res := response.ErrDatabase
			c.Error(result.Error)
			c.JSON(response.Error(res))
			return
		}

		clientPublic.BindAttributesFrom(client)
		c.JSON(response.Success(clientPublic, response.SuccessCreate))
		
	}
}