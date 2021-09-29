package oauth

import (
	"aherman/src/enums"
	"aherman/src/http/response"
	oauthModels "aherman/src/models/oauth"
	userModels "aherman/src/models/user"
	"aherman/src/types/container"
	"aherman/src/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// PasswordGrant todo
func PasswordGrant(app *container.Container) gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			// client      *clientModels.Client
			credentials oauthModels.PasswordGrantCredentials
			user        *userModels.User
		)

		// bind input to credentials variable.
		// if we can't, there was a validation error.
		if err := c.ShouldBindBodyWith(&credentials, binding.JSON); err != nil {
			c.Error(err)
			res := response.ErrValidation
			res.Description = err.Error()
			c.JSON(response.Error(res))
			return
		}

		// Get the client and check credentials
		// client, err := app.Client.FindByID(credentials.ClientID.String())
		// if err != nil {
		// 	c.Error(err)
		// 	res := response.ErrDatabase
		// 	c.JSON(response.Error(res))
		// 	return
		// }
		// if client.ID == uuid.Nil {
		// 	res := response.ErrNotFound
		// 	c.JSON(response.Error(res))
		// 	return
		// }
		// if client.Secret != credentials.ClientSecret {
		// 	res := response.ErrBadCredentials
		// 	c.JSON(response.Error(res))
		// 	return
		// }

		// Get the user and check credentials
		user, err := app.User.FindByEmail(credentials.Email)
		if err != nil {
			c.Error(err)
			res := response.ErrDatabase
			c.JSON(response.Error(res))
		}
		if user.ID == uuid.Nil {
			res := response.ErrNotFound
			c.JSON(response.Error(res))
			return
		}
		if !util.PasswordIsValid(credentials.Password, user.Password) {
			res := response.ErrBadCredentials
			c.JSON(response.Error(res))
			return
		}

		// passed all checks, so let's generate a token
		exp := enums.JWTAccessTokenExpiry
		claims := &userModels.JWTClaims{
			ID: user.ID,
			StandardClaims: jwt.StandardClaims{
				// In JWT, the expiry time is expressed as unix milliseconds
				ExpiresAt: exp.Unix(),
			},
		}
		token := jwt.NewWithClaims(enums.JWTSigningMethod, *claims)

		// convert to string
		tokenString, err := token.SignedString([]byte(enums.AppSecret))
		if err != nil {
			c.Error(err)
			res := response.ErrUnknown
			c.JSON(response.Error(res))
			return
		}

		c.JSON(response.Success(
			struct {
				AccessToken string `json:"access_token"`
			}{
				AccessToken: tokenString,
			},
			response.SuccessCreate,
		))

	}
}