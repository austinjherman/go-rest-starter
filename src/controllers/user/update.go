package user

import (
	"aherman/src/http/response"
	u "aherman/src/models/user"
	"aherman/src/types/container"
	"aherman/src/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// Update gets a user from the database by ID
func Update(app *container.Container) gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			user           *u.User = app.User.Current
			userUpdateable *u.Updateable
			userPublic     *u.Public
		)

		// bind input to user variable.
		// if we can't, there was a validation error.
		if err := c.ShouldBindBodyWith(userUpdateable, binding.JSON); err != nil {
			res := response.ErrValidation
			c.Error(res)
			c.JSON(response.Error(res))
			return
		}

		if userUpdateable.Email != "" {
			user.Email = userUpdateable.Email
		}

		if userUpdateable.Password != "" {
			var err error
			user.Password, err = util.HashPassword(userUpdateable.Password)
			if err != nil {
				c.Error(err)
				res := response.ErrUnknown
				c.JSON(response.Error(res))
				return
			}
		}

		app.User.DB.Save(user)

		userPublic.BindAttributes(user)
		c.JSON(response.Success(userPublic, response.SuccessRead))
	}
}