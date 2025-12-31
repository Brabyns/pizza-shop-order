package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Brabyns/pizza-shop-order/internal/models"
	"github.com/gin-gonic/gin"
)

type LoginData struct {
	Error string
}

type AdminDashaboardData struct {
	Username string
}

func (h *Handler) HandleLoginGet(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", LoginData{})
}

func (h *Handler) HandleLoginPost(c *gin.Context) {
	var form struct {
		Username string `form:"username" binding:"required,min=3,max=50"`
		Password string `form:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBind(&form); err != nil {
		c.HTML(http.StatusBadRequest, "login.tmpl", LoginData{
			Error: "Invalid input",
		})
		return
	}

	user, err := h.users.AuthenticateUser(form.Username, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			c.HTML(http.StatusUnauthorized, "login.tmpl", LoginData{
				Error: "Invalid username or password",
			})
			return
		}

		c.String(http.StatusInternalServerError, "Server error")
		return
	}

	SetSessionValue(c, "userID", fmt.Sprintf("%d", user.ID))
	SetSessionValue(c, "username", user.Username)

	c.Redirect(http.StatusSeeOther, "/admin")
}

func (h *Handler) HandleLogout(c *gin.Context) {
	if err := ClearSession(c); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/login")
}


func (h *Handler) ServeAdminDashboard(c *gin.Context) {
	username := GetSessionString(c, "username")

	if username == "" {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	c.HTML(http.StatusOK, "admin.tmpl", AdminDashaboardData{
		Username: username,
	})
}