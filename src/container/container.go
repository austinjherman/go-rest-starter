package container

import (
	tokenFacades "aherman/src/facades/token"
	userFacades "aherman/src/facades/user"
	tokenModels "aherman/src/models/token"
	userModels "aherman/src/models/user"
)

// CurrentContainer todo
type CurrentContainer struct {
	Token *tokenModels.Token
	User *userModels.User 
}

// Container for dependency injection
type Container struct {
	Current *CurrentContainer
	Token *tokenFacades.Token
	User *userFacades.User
}
