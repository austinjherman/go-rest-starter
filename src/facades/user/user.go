package user

import (
	"aherman/src/http/response"
	tokenModels "aherman/src/models/token"
	userModels "aherman/src/models/user"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User is a facade that allows to interact with the user database
type User struct {
	DB *gorm.DB
}

// EmailIsAvailable checks to see if the user's email is already registered.
func (env *User) EmailIsAvailable(email string) error {
	var u userModels.User
	// safe because this inline query is escaped by GORM / database/sql
	// https://gorm.io/it_IT/docs/security.html
	if result := env.DB.First(&u, "email = ?", email); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil
		}
		return result.Error
	}
	if u.ID != uuid.Nil {
		return response.ErrEmailNotAvailable
	}
	return nil
}

// FindByEmail todo
func (env *User) FindByEmail(email string) (*userModels.User, error) {
	var u userModels.User
	// safe because this inline query is escaped by GORM / database/sql
	// https://gorm.io/it_IT/docs/security.html
	result := env.DB.Limit(1).Find(&u, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}
	if u.ID == uuid.Nil {
		return nil, response.ErrEmailNotFound
	}
	return &u, nil
}

// FindByID todo
func (env *User) FindByID(id uuid.UUID) (*userModels.User, error) {
	var u userModels.User
	// safe because this inline query is escaped by GORM / database/sql
	// https://gorm.io/it_IT/docs/security.html
	result := env.DB.Limit(1).Find(&u, "id = ?", id.String())
	if result.Error != nil {
		return nil, result.Error
	}
	if u.ID == uuid.Nil {
		return nil, response.ErrNotFound
	}
	return &u, nil
}

// RevokeSession todo
func (env *User) RevokeSession(userID uuid.UUID, sessionID string) error {

	// tokens array for binding
	tokens := []tokenModels.Token{}

	// query for tokens with matching user_id and session_id
	result := env.DB.Model(&tokenModels.Token{}).
		Where("user_id = ? AND session_id = ?", userID, sessionID).
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