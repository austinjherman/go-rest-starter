package user

import (
	"aherman/src/http/response"
	u "aherman/src/models/user"
	"aherman/src/types/container"

	"github.com/gin-gonic/gin"
)

// Read gets a user from the database by ID
func Read(app *container.Container) gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			user       *u.User = app.User.Current
			userPublic *u.Public = &u.Public{}
		)

		userPublic.BindAttributes(user)
		c.JSON(response.Success(userPublic, response.SuccessRead))
	}
}