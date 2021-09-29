package container

import (
	"aherman/src/models/client"
	"aherman/src/models/oauth"
	"aherman/src/models/user"
)

// Creator interface
type Creator interface {
	// Create(c *gin.Context)
	// Read(c *gin.Context)
	// Update(c *gin.Context)
	// Delete(c *gin.Context)
}

// Container for dependency injection
type Container struct {
	Client *client.Env
	OAuth *oauth.Env
	User *user.Env
}
