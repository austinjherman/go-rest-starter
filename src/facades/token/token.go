package token

import (
	"aherman/src/http/response"
	tokenModels "aherman/src/models/token"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Token is a facade that allows us to interact with the token database
type Token struct {
	DB *gorm.DB
	SessionID string
}

// InWhitelist todo
func (env *Token) InWhitelist(id uuid.UUID) (bool, error) {
	var t tokenModels.Token
	// safe because this inline query is escaped by GORM / database/sql
	// https://gorm.io/it_IT/docs/security.html
	result := env.DB.Limit(1).Find(&t, "id = ?", id.String())
	if result.Error != nil {
		return false, result.Error
	}
	if t.ID == uuid.Nil {
		return false, response.ErrNotFound
	}
	return true, nil
}