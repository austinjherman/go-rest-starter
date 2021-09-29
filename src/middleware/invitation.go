package middleware

import (
	"aherman/src/enums"
	"aherman/src/http/response"
	"errors"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type invitationRequest struct {
	InvitationCode string `json:"invitation_code"`
}

// Invitation is a middleware that requires the invitation code
// be provided in the JSON request
func Invitation() gin.HandlerFunc {
	return func(c *gin.Context) {
		
		// todo: should move this to it's own middleware
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		request := invitationRequest{}
		err := c.ShouldBindBodyWith(&request, binding.JSON)

		if err != nil {
			res := response.ErrUnknown
			c.Error(res)
			c.AbortWithStatusJSON(res.Status, res)
			return
		}

		if request.InvitationCode == "" {
			res := response.ErrNoInvitation
			c.Error(res)
			c.AbortWithStatusJSON(res.Status, res)
			return
		}

		if request.InvitationCode != os.Getenv(enums.InvitationCodeEnvironmentKey) {
			res := response.ErrNoInvitation
			c.Error(errors.New("Incorrect invitation code"))
			c.AbortWithStatusJSON(res.Status, res)
			return
		}

		c.Next()
	}
}
