package container

import (
	errorFacades "aherman/src/facades/error"
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

// Facades todo
type Facades struct {
	Error *errorFacades.Error
	Token *tokenFacades.Token
	User *userFacades.User
}

// Container for dependency injection
type Container struct {
	Current *CurrentContainer
	Facades *Facades
}
