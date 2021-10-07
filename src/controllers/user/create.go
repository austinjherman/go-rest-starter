package user

import (
	"aherman/src/container"
	"aherman/src/http/response"
	userModels "aherman/src/models/user"
	"aherman/src/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// Create a user and persist to database
func Create(app *container.Container) gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			user          userModels.User
			userCreatable userModels.Creatable
			userPublic    userModels.Public
		)

		// validate user input
		err := c.ShouldBindBodyWith(&userCreatable, binding.JSON)
		ok, httpResponse := app.Facades.Error.ShouldContinue(err, &response.ErrValidation)
		if !ok {
			c.JSON(response.Error(httpResponse))
			return
		}

		// check if email is available
		err = app.Facades.User.EmailIsAvailable(userCreatable.Email)
		ok, httpResponse = app.Facades.Error.ShouldContinue(err, &response.ErrUserEmailAlreadyRegistered)
		if !ok {
			c.JSON(response.Error(httpResponse))
			return
		}

		// checks passed; create user
		passwordHash, err := util.HashPassword(userCreatable.Password)
		ok, httpResponse = app.Facades.Error.ShouldContinue(err, &response.ErrUnknown)
		if !ok {
			c.JSON(response.Error(httpResponse))
			return
		}
		
		user.Email = userCreatable.Email
		user.Password = passwordHash

		result := app.Facades.User.DB.Create(&user)
		ok, httpResponse = app.Facades.Error.ShouldContinue(result.Error, &response.ErrUnknown)
		if !ok {
			c.JSON(response.Error(httpResponse))
			return
		}

		userPublic.BindAttributes(&user)
		c.JSON(response.Success(response.SuccessCreate, userPublic))
	}
}