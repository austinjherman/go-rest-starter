/**************************************************************************************
  User Facade
	-------------
	The user facade has access to the user database connection. It is passed as a
	dependency to all controllers. With it, we can abstract some of the logic out of
	the controllers and handle it here.
**************************************************************************************/

package user

import (
	"aherman/src/http/response"
	tokenModels "aherman/src/models/token"
	userModels "aherman/src/models/user"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User is a facade that allows to interact with the user database
type User struct {
	DB *gorm.DB
}


/**************************************************************************************
  Private Helper Functions
**************************************************************************************/

// get binds a user by the query and args or returns an error
func (env *User) get(user *userModels.User, query interface{}, args ...interface{}) error {

	// run the query
	result := env.DB.Model(&userModels.User{}).
		Where(query, args...).
		Find(&user)
	
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// revoke tokens will revoke a user's tokens by query and args
func (env *User) revokeTokens(query interface{}, args ...interface{}) error {
	// tokens array for binding
	tokens := []tokenModels.Token{}

	// query for tokens with matching user_id and session_id
	result := env.DB.Model(&tokenModels.Token{}).
		Where(query, args...).
		Find(&tokens)

	if result.Error != nil {
		return result.Error
	}

	// tokenIDS for binding
	tokenIDs := []uuid.UUID{}
	for _, token := range tokens {
		tokenIDs = append(tokenIDs, token.ID)
	}

	// delete all tokens with matching user_id and session_id
	err := env.DB.Delete(&tokenModels.Token{}, tokenIDs).Error
	if err != nil {
		return err
	}

	return nil
}


/**************************************************************************************
  Public Functions
**************************************************************************************/

// BindByEmail binds user the email or returns an error
func (env *User) BindByEmail(user *userModels.User, email string) error {
	err := env.get(user, "email = ?", email)
	if err != nil {
		return err
	}
	if user.ID == uuid.Nil {
		return response.ErrUserNotFound
	}
	return nil
}

// BindByID binds user by id or returns an error
func (env *User) BindByID(user *userModels.User, id string) error {
	err := env.get(user, "id = ?", id)
	if err != nil {
		return err
	}
	if user.ID == uuid.Nil {
		return response.ErrUserNotFound
	}
	return nil
}

// EmailIsAvailable checks whether or not an email is available
func (env *User) EmailIsAvailable(email string) error {
	user := &userModels.User{}
	err := env.get(user, "email = ?", email)
	if err != nil {
		return err
	}
	if user.ID != uuid.Nil {
		return response.ErrUserEmailAlreadyRegistered
	}
	return nil
}

// RevokeTokenByID todo
func (env *User) RevokeTokenByID(tokenID string) error {
	return env.revokeTokens("id = ?", tokenID)
}

// RevokeTokensBySession todo
func (env *User) RevokeTokensBySession(userID uuid.UUID, sessionID string) error {
	return env.revokeTokens("user_id = ? AND session_id = ?", userID, sessionID)
}

// RevokeTokensByUser todo
func (env *User) RevokeTokensByUser(userID uuid.UUID) error {
	return env.revokeTokens("user_id = ?", userID)
}

// WhitelistTokens todo
func (env *User) WhitelistTokens(userID uuid.UUID, tokens []tokenModels.Token) error {
	var u *userModels.User = &userModels.User{}
	u.ID = userID
	err := env.DB.Model(u).Association(userModels.TokensFieldName).Append(tokens)
	if err != nil {
		return err
	}
	return nil
}