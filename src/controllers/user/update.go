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
		err := c.ShouldBindBodyWith(userUpdateable, binding.JSON)
		ok, httpResponse := app.Facades.Error.ShouldContinue(err, &response.ErrValidation)
		if !ok {
			c.JSON(response.Error(httpResponse))
			return
		}

		// update email if not empty
		if userUpdateable.Email != "" {
			err := app.Facades.User.EmailIsAvailable(userUpdateable.Email)
			ok, httpResponse := app.Facades.Error.ShouldContinue(err, &response.ErrUserEmailAlreadyRegistered)
			if !ok {
				c.JSON(response.Error(httpResponse))
				return
			}

			user.Email = userUpdateable.Email
		}

		// update password if not empty
		if userUpdateable.PasswordNew != "" {

			if !util.PasswordIsValid(userUpdateable.PasswordOld, user.Password) {
				httpResponse := response.ErrUserPasswordInvalid
				c.JSON(response.Error(httpResponse))
				return
			}

			user.Password, err = util.HashPassword(userUpdateable.PasswordNew)
			ok, httpResponse := app.Facades.Error.ShouldContinue(err, &response.ErrUnknown)
			if !ok {
				c.JSON(response.Error(httpResponse))
				return
			}
		}

		result := app.Facades.User.DB.Save(user)
		ok, httpResponse = app.Facades.Error.ShouldContinue(result.Error, &response.ErrResourceNotUpdated)
		if !ok {
			c.JSON(response.Error(httpResponse))
			return
		}

		userPublic.BindAttributes(user)
		c.JSON(response.Success(response.SuccessUpdate, userPublic))
	}
}