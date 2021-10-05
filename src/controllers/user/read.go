package user

import (
	"aherman/src/container"
	"aherman/src/http/response"
	u "aherman/src/models/user"

	"github.com/gin-gonic/gin"
)

// Read gets a user from the database by ID
func Read(app *container.Container) gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			user       *u.User = app.Current.User
			userPublic *u.Public = &u.Public{}
		)

		userPublic.BindAttributes(user)
		c.JSON(response.Success(response.SuccessRead, userPublic))
	}
}