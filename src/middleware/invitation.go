package middleware

import (
	"aherman/src/container"
	"aherman/src/enums"
	"aherman/src/http/response"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type invitationRequest struct {
	InvitationCode string `json:"invitation_code"`
}

// Invitation is a middleware that requires the invitation code
// be provided in the JSON request
func Invitation(app *container.Container) gin.HandlerFunc {
	return func(c *gin.Context) {

		request := &invitationRequest{}

		err := c.ShouldBindBodyWith(request, binding.JSON)
		ok, httpResponse := app.Facades.Error.ShouldContinue(err, &response.ErrValidation)
		if !ok {
			c.AbortWithStatusJSON(httpResponse.Status, httpResponse)
			return
		}

		if request.InvitationCode == "" {
			httpResponse := response.ErrInvitationCodeInvalid
			c.AbortWithStatusJSON(httpResponse.Status, httpResponse)
			return
		}

		if request.InvitationCode != enums.InvitationCode {
			httpResponse := response.ErrInvitationCodeInvalid
			c.AbortWithStatusJSON(httpResponse.Status, httpResponse)
			return
		}

		c.Next()
	}
}
