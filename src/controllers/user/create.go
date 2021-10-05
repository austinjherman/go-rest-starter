package user

import (
	"aherman/src/container"
	"aherman/src/http/response"
	u "aherman/src/models/user"
	"aherman/src/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// Create a user and persist to database
func Create(app *container.Container) gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			user          u.User
			userCreatable u.Creatable
			userPublic    u.Public
		)

		// bind input to user variable.
		// if we can't, there was a validation error.
		if err := c.ShouldBindBodyWith(&userCreatable, binding.JSON); err != nil {
			c.Error(err)
			res := response.ErrValidation
			res.Description = err.Error()
			c.JSON(response.Error(res))
			return
		}

		// check if email is available
		if err := app.User.EmailIsAvailable(userCreatable.Email); err != nil {
			c.Error(err)
			c.JSON(response.Error(err))
			return
		}

		// checks passed; create user
		passwordHash, err := util.HashPassword(userCreatable.Password)
		if err != nil {
			c.Error(err)
			c.JSON(response.Error(err))
			return
		}
		
		user.Email = userCreatable.Email
		user.Password = passwordHash

		if result := app.User.DB.Create(&user); result.Error != nil {
			c.Error(result.Error)
			c.JSON(response.Error(result.Error))
			return
		}

		userPublic.BindAttributes(&user)
		c.JSON(response.Success(response.SuccessCreate, userPublic))
	}
}