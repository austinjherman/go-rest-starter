package user

import (
	"aherman/src/http/response"
	userModels "aherman/src/models/user"
	"aherman/src/types/container"

	"github.com/gin-gonic/gin"
)

// Delete removes a user from the database
func Delete(app *container.Container) gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			user       *userModels.User = app.User.Current
			userPublic *userModels.Public
		)

		userPublic.BindAttributes(user)

		result := app.User.DB.Delete(&userModels.User{}, user.ID)

		if result.Error != nil {
			c.Error(result.Error)
			res := response.ErrDatabase
			c.JSON(response.Error(res))
			return
		}

		c.JSON(response.Success(userPublic, response.SuccessDelete))

	}
}