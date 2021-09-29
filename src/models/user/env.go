package user

import (
	"aherman/src/http/response"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Env todo
type Env struct {
	Current *User
	DB *gorm.DB
}

// EmailIsAvailable checks to see if the user's email is already registered.
func (env *Env) EmailIsAvailable(email string) error {
	var u User
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
func (env *Env) FindByEmail(email string) (*User, error) {
	var u User
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
func (env *Env) FindByID(id uuid.UUID) (*User, error) {
	var u User
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