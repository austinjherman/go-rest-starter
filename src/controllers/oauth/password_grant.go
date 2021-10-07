package oauth

import (
	"aherman/src/container"
	"aherman/src/enums"
	"aherman/src/http/response"
	"aherman/src/logging"
	oauthModels "aherman/src/models/oauth"
	tokenModels "aherman/src/models/token"
	userModels "aherman/src/models/user"
	"aherman/src/util"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// PasswordGrant is now considered a legacy grant type by OAuth2, however,
// we still need a way to allow users to log into our application.
// This grant type is only allowed for trusted clients.
func PasswordGrant(app *container.Container) gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			credentials *oauthModels.PasswordGrantRequest = &oauthModels.PasswordGrantRequest{}
			user *userModels.User = &userModels.User{}
		)

		// bind input. if we can't, there was a validation error.
		err := c.ShouldBindBodyWith(credentials, binding.JSON)
		ok, httpResponse := app.Facades.Error.ShouldContinue(err, &response.ErrValidation)
		if !ok {
			c.JSON(response.Error(httpResponse))
			return
		}

		// check if the request was made by a trusted client
		if credentials.ClientID != enums.TrustedClientID || credentials.ClientSecret != enums.TrustedClientSecret {
			httpResponse := response.ErrClientUnknown
			serverErr := logging.NewServerError(errors.New("request made by an untrusted client")).
				BindClientErr(&httpResponse)
			c.Error(serverErr)
			c.JSON(response.Error(httpResponse))
			return
		}

		// Get the user and check credentials
		err = app.Facades.User.BindByEmail(user, credentials.Email)
		ok, httpResponse = app.Facades.Error.ShouldContinue(err, &response.ErrUserEmailNotFound)
		if !ok {
			c.JSON(response.Error(httpResponse))
			return
		}

		if !util.PasswordIsValid(credentials.Password, user.Password) {
			httpResponse := response.ErrUserPasswordInvalid
			c.JSON(response.Error(httpResponse))
			return
		}

		// passed all checks, so let's generate tokens
		accessToken, accessTokenClaims, err := app.Facades.Token.NewToken(
			enums.JWTTokenTypeAccess, user.ID,
		)
		ok, httpResponse = app.Facades.Error.ShouldContinue(err, &response.ErrUnknown)
		if err != nil {
			c.JSON(response.Error(httpResponse))
			return
		}

		refreshToken, _, err := app.Facades.Token.NewToken(
			enums.JWTTokenTypeRefresh, user.ID,
		)
		ok, httpResponse = app.Facades.Error.ShouldContinue(err, &response.ErrUnknown)
		if err != nil {
			c.JSON(response.Error(httpResponse))
			return
		}

		// whitelist the provided tokens
		tokens := []tokenModels.Token{
			*accessToken,
			*refreshToken,
		}

		err = app.Facades.User.WhitelistTokens(user.ID, tokens)
		ok, httpResponse = app.Facades.Error.ShouldContinue(err, &response.ErrUnknown)
		if !ok {
			c.JSON(response.Error(httpResponse))
			return
		}

		data := &oauthModels.Success{
			AccessToken: accessToken.Token,
			ExpiresIn: time.Unix(accessTokenClaims.ExpiresAt, 0),
			RefreshToken: refreshToken.Token,
			Scope: "",
			TokenType: enums.JWTTokenTypeBearerOutgoing,
		}

		c.JSON(response.Success(response.SuccessLogin, data))
	}
}