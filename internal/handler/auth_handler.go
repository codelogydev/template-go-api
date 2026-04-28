package handler

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/codelogydev/core-go/oauth"
	"github.com/codelogydev/core-go/response"
	"github.com/codelogydev/template-go-api/internal/service"
)

type AuthHandler struct {
	service service.UserService
}

func NewAuthHandler(service service.UserService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	b := make([]byte, 16)
	rand.Read(b)
	state := hex.EncodeToString(b)

	c.SetCookie("oauth_state", state, 300, "/", "", false, true)
	c.Redirect(http.StatusTemporaryRedirect, oauth.GoogleAuthURL(state))
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	cookieState, err := c.Cookie("oauth_state")
	if err != nil || cookieState != c.Query("state") {
		response.BadRequest(c, "invalid oauth state")
		return
	}

	userInfo, err := oauth.GoogleExchange(c.Request.Context(), c.Query("code"))
	if err != nil {
		response.Error(c, 500, "google authentication failed")
		return
	}

	result, err := h.service.LoginWithGoogle(userInfo.ID, userInfo.Email, userInfo.Name)
	if err != nil {
		response.Error(c, 500, "login failed")
		return
	}

	response.Success(c, result)
}
