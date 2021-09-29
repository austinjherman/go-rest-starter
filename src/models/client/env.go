package client

import "gorm.io/gorm"

// Env todo
type Env struct {
	Client Client
	DB *gorm.DB
}

// FindByID finds a client record by client ID.
func (env *Env) FindByID(id string) (*Client, error) {
	client := Client{}
	result := env.DB.Find(&client, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &client, nil
}