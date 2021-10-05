package oauth

import (
	"aherman/src/container"
	"aherman/src/enums"
	"aherman/src/http/response"
	oauthModels "aherman/src/models/oauth"
	tokenModels "aherman/src/models/token"
	userModels "aherman/src/models/user"
	"aherman/src/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// PasswordGrant is now considered a legacy grant type by OAuth2, however,
// we still need a way to allow users to log into our application.
// This grant type is only allowed for trusted clients.
func PasswordGrant(app *container.Container) gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			credentials *oauthModels.PasswordGrantRequest = &oauthModels.PasswordGrantRequest{}
		)

		// bind input. if we can't, there was a validation error.
		if err := c.ShouldBindBodyWith(credentials, binding.JSON); err != nil {
			c.Error(err)
			res := response.ErrValidation
			res.Description = err.Error()
			c.JSON(response.Error(res))
			return
		}

		// check if the request was made by a trusted client
		if credentials.ClientID != enums.TrustedClientID || credentials.ClientSecret != enums.TrustedClientSecret {
			err := response.ErrNotATrustedClient
			c.Error(err)
			c.JSON(response.Error(err))
			return
		}

		// Get the user and check credentials
		user, err := app.User.FindByEmail(credentials.Email)
		if err != nil {
			c.Error(err)
			c.JSON(response.Error(err))
		}
		if user.ID == uuid.Nil {
			res := response.ErrEmailNotFound
			c.JSON(response.Error(res))
			return
		}
		if !util.PasswordIsValid(credentials.Password, user.Password) {
			res := response.ErrInvalidPassword
			c.JSON(response.Error(res))
			return
		}

		// passed all checks, so let's generate tokens

		// first we'll generate a session token so we can revoke tokens
		// by session when a user logs out (i.e. all tokens generated during this request)
		sessionID := util.StringRandom(16)

		// access token
		accessTokenExp := enums.JWTAccessTokenExpiry
		accessTokenID := uuid.New()
		accessTokenClaims := &userModels.JWTClaims{
			ID: accessTokenID,
			SessionID: sessionID,
			StandardClaims: jwt.StandardClaims{
				// In JWT, the expiry time is expressed as unix milliseconds
				ExpiresAt: accessTokenExp.Unix(),
			},
			TokenType: enums.JWTTokenTypeAccess,
			UserID: user.ID,
		}
		accessToken := jwt.NewWithClaims(enums.JWTSigningMethod, *accessTokenClaims)

		// convert to string
		accessTokenString, err := accessToken.SignedString([]byte(enums.AppSecret))
		if err != nil {
			c.Error(err)
			res := response.ErrUnknown
			c.JSON(response.Error(res))
			return
		}

		// refresh token
		refreshTokenExp := enums.JWTRefreshTokenExpiry
		refreshTokenID := uuid.New()
		refreshTokenClaims := &userModels.JWTClaims{
			ID: refreshTokenID,
			SessionID: sessionID,
			StandardClaims: jwt.StandardClaims{
				// In JWT, the expiry time is expressed as unix milliseconds
				ExpiresAt: refreshTokenExp.Unix(),
			},
			TokenType: enums.JWTTokenTypeRefresh,
			UserID: user.ID,
		}
		refreshToken := jwt.NewWithClaims(enums.JWTSigningMethod, *refreshTokenClaims)

		// convert to string
		refreshTokenString, err := refreshToken.SignedString([]byte(enums.AppSecret))
		if err != nil {
			c.Error(err)
			res := response.ErrUnknown
			c.JSON(response.Error(res))
			return
		}

		data := &oauthModels.Success{
			AccessToken: accessTokenString,
			ExpiresIn: accessTokenExp,
			RefreshToken: refreshTokenString,
			Scope: "",
			TokenType: enums.JWTTokenTypeBearerOutgoing,
		}

		tokens := []tokenModels.Token{
			{
				ID: accessTokenID,
				UserID: user.ID,
				SessionID: sessionID,
				Token: accessTokenString,
				Type: enums.JWTTokenTypeAccess,
			},
			{
				ID: refreshTokenID,
				UserID: user.ID,
				SessionID: sessionID,
				Token: refreshTokenString,
				Type: enums.JWTTokenTypeRefresh,
			},
		}

		// whitelist the provided tokens
		err = app.User.WhitelistTokens(user.ID, tokens)
		if err != nil {
			c.Error(err)
			res := response.ErrUnknown
			c.JSON(response.Error(res))
			return
		}

		c.JSON(response.Success(response.SuccessLogin, data))
	}
}