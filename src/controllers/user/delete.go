package user

import (
	"aherman/src/http/response"
	userModels "aherman/src/models/user"
	"aherman/src/container"

	"github.com/gin-gonic/gin"
)

// Delete removes a user from the database
func Delete(app *container.Container) gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			user       *userModels.User   = app.Current.User
			userPublic *userModels.Public = &userModels.Public{}
		)

		userPublic.BindAttributes(user)

		result := app.User.DB.Delete(&userModels.User{}, user.ID)

		if result.Error != nil {
			c.Error(result.Error)
			c.JSON(response.Error(result.Error))
			return
		}

		c.JSON(response.Success(response.SuccessDelete, userPublic))
	}
}