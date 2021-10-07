package oauth

import (
	"aherman/src/container"
	"aherman/src/enums"
	"aherman/src/http/response"
	oauthModels "aherman/src/models/oauth"
	tokenModels "aherman/src/models/token"
	userModels "aherman/src/models/user"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang-jwt/jwt"
)

type authorizationHeader struct {
	Authorization string `header:"Authorization"`
}

// RefreshForAccess contains the logic that allows us to swap refresh tokens
// for access tokens
func RefreshForAccess(app *container.Container) gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			bearer *authorizationHeader = &authorizationHeader{}
			claims *tokenModels.JWTClaims = &tokenModels.JWTClaims{}
			jwtToken *jwt.Token = &jwt.Token{}
			token *tokenModels.Token = &tokenModels.Token{}
			user *userModels.User = &userModels.User{}
		)

		// validate header
		err := c.ShouldBindWith(bearer, binding.Header)
		ok, httpResponse := app.Facades.Error.ShouldContinue(err, &response.ErrValidation)
		if !ok {
			c.JSON(httpResponse.Status, httpResponse)
			return
		}

		// split bearer header
		bearerHeader := strings.Split(bearer.Authorization, " ")
		if len(bearerHeader) != 2 || bearerHeader[0] != "Bearer" {
			httpResponse := response.ErrTokenInvalid
			c.JSON(httpResponse.Status, httpResponse)
			return
		}

		// validate the token
		jwtToken, err = app.Facades.Token.ParseWithClaims(bearerHeader[1], claims)
		ok, httpResponse = app.Facades.Error.ShouldContinue(err, &response.ErrTokenInvalid)
		if !ok {
			c.JSON(response.Error(httpResponse))
			return
		}

		// find the user
		err = app.Facades.User.BindByID(user, claims.Subject)
		ok, httpResponse = app.Facades.Error.ShouldContinue(err, &response.ErrResourceNotFound)
		if !ok {
			c.JSON(response.Error(httpResponse))
			return
		}

		// explicitly check for expiry
		// todo: revoke this token if it it's still in the whitelist
		now := time.Now()
		tokenExpiry := time.Unix(claims.ExpiresAt, 0)
		if tokenExpiry.Before(now) {
			app.Facades.User.RevokeTokenByID(claims.ID.String())
			httpResponse := response.ErrTokenExpired
			c.JSON(response.Error(httpResponse))
			return
		}

		// We'll only allow the refresh token flow
		if claims.TokenType != enums.JWTTokenTypeRefresh {
			httpResponse := response.ErrTokenTypeInvalid
			c.JSON(response.Error(httpResponse))
			return
		}

		// We'll disallow invalid tokens, too.
		if !jwtToken.Valid {
			httpResponse := response.ErrTokenInvalid
			c.JSON(response.Error(httpResponse))
			return
		}

		// check to see if the token is in the whitelist
		err = app.Facades.Token.BindByID(token, claims.ID)
		ok, httpResponse = app.Facades.Error.ShouldContinue(err, &response.ErrTokenNotFound)
		if !ok {
			c.JSON(response.Error(httpResponse))
			return
		}

		// passed all checks so let's generate a token
		accessToken, accessTokenClaims, err := app.Facades.Token.NewToken(
			enums.JWTTokenTypeAccess, user.ID,
		)
		ok, httpResponse = app.Facades.Error.ShouldContinue(err, &response.ErrUnknown)
		if !ok {
			c.JSON(response.Error(httpResponse))
			return
		}

		// whitelist the provided tokens
		tokens := []tokenModels.Token{
			*accessToken,
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
			Scope: "",
			TokenType: enums.JWTTokenTypeBearerOutgoing,
		}

		c.JSON(response.Success(response.SuccessLogin, data))
	}
}
