/**************************************************************************************
  Token Whitelist Facade
	-----------------------
	The token whitelist facade has access to the token whitelist database connection.
	It is passed as a dependency to all controllers. With it, we can abstract some of
	the logic out the the controllers and handle it here.
**************************************************************************************/

package token

import (
	"aherman/src/enums"
	"aherman/src/http/response"
	tokenModels "aherman/src/models/token"
	"aherman/src/util"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Token is a facade that allows us to interact with the token database
type Token struct {
	DB *gorm.DB
}

/**************************************************************************************
  Private Helper Functions
**************************************************************************************/

// get binds a token by the query and args or returns an error
func (env *Token) get(token *tokenModels.Token, query interface{}, args ...interface{}) error {

	// run the query
	result := env.DB.Model(&tokenModels.Token{}).
		Where(query, args...).
		Find(&token)
	
	if result.Error != nil {
		return result.Error
	}

	if token.ID == uuid.Nil {
		return response.ErrTokenNotFound
	}

	return nil
}


/**************************************************************************************
  Public Functions
**************************************************************************************/

// BindByID todo
func (env *Token) BindByID(token *tokenModels.Token, id uuid.UUID) error {
	return env.get(token, "id = ?", id)
}

// ParseWithClaims todo
func (env *Token) ParseWithClaims(tokenStr string, claims *tokenModels.JWTClaims) (*jwt.Token, error) {

	// parse the token with the user claims
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			err := fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			return nil, err
		}
		return []byte(enums.AppSecret), nil
	})

	// check for token parsing error
	if err != nil {
		return nil, err
	}

	return token, nil
}

// NewToken todo
func (env *Token) NewToken(tokenType string, subjectID uuid.UUID) (
	*tokenModels.Token, *tokenModels.JWTClaims, error,
) {

	// only allow known token types
	if tokenType != enums.JWTTokenTypeAccess && tokenType != enums.JWTTokenTypeRefresh {
		return nil, nil, errors.New("unknown token type requested")
	}

	// set the token expiry
	tokenExp := time.Time{}
	if tokenType == enums.JWTTokenTypeAccess {
		tokenExp = enums.JWTAccessTokenExpiry
	}
	if tokenType == enums.JWTTokenTypeRefresh {
		tokenExp = enums.JWTRefreshTokenExpiry
	}

	// generate ids and expiry
	sessionID := util.StringRandom(16)
	tokenID := uuid.New()

	// create the claims
	claims := &tokenModels.JWTClaims{
		ID: tokenID,
		SessionID: sessionID,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: tokenExp.Unix(),
			Id: tokenID.String(),
			IssuedAt: time.Now().Unix(),
			Subject: subjectID.String(),
		},
		TokenType: tokenType,
	}

	jwtToken := jwt.NewWithClaims(enums.JWTSigningMethod, *claims)

	token := &tokenModels.Token{
		ID: tokenID,
		UserID: subjectID,
		SessionID: sessionID,
		Type: tokenType,
	}
	jwtTokenStr, err := util.JWTToString(jwtToken)
	if err != nil {
		return nil, nil, err
	}
	token.Token = jwtTokenStr
	
	return token, claims, nil
}
