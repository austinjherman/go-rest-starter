package user

import (
	"aherman/src/container"
	"aherman/src/http/response"
	u "aherman/src/models/user"
	"aherman/src/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// Update gets a user from the database by ID
func Update(app *container.Container) gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			user           *u.User = app.Current.User
			userUpdateable *u.Updateable = &u.Updateable{}
			userPublic     *u.Public = &u.Public{}
		)

		// bind input to user variable.
		// if we can't, there was a validation error.
		if err := c.ShouldBindBodyWith(userUpdateable, binding.JSON); err != nil {
			c.Error(err)
			res := response.ErrValidation
			res.Description = err.Error()
			c.JSON(response.Error(res))
			return
		}

		if userUpdateable.Email != "" {
			err := app.User.EmailIsAvailable(userUpdateable.Email)
			if err != nil {
				c.Error(err)
				c.JSON(response.Error(err))
				return
			}

			user.Email = userUpdateable.Email
		}

		if userUpdateable.Password != "" {
			var err error
			user.Password, err = util.HashPassword(userUpdateable.Password)
			if err != nil {
				c.Error(err)
				c.JSON(response.Error(err))
				return
			}
		}

		result := app.User.DB.Save(user)
		if result.Error != nil {
			c.Error(result.Error)
			c.JSON(response.Error(result.Error))
			return
		}

		userPublic.BindAttributes(user)
		c.JSON(response.Success(response.SuccessUpdate, userPublic))
	}
}