package user

import (
	"aherman/src/http/response"
	u "aherman/src/models/user"
	"aherman/src/types/container"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// Read gets a user from the database by ID
func Read(app *container.Container) gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			user         *u.User
			userReadable *u.Readable
			userPublic   *u.Public
		)

		// bind input to user variable.
		// if we can't, there was a validation error.
		if err := c.ShouldBindBodyWith(&userReadable, binding.JSON); err != nil {
			res := response.ErrValidation
			c.Error(res)
			c.JSON(response.Error(res))
			return
		}

		// find user
		user, err := app.User.FindByID(userReadable.ID)
		if err != nil {
			c.Error(err)
			c.JSON(response.Error(err))
			return
		}

		userPublic.BindAttributes(user)
		c.JSON(response.Success(userPublic, response.SuccessRead))
	}
}