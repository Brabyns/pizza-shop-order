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
	Orders []models.Order
	Statuses []string
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
	orders, err := h.orders.GetAllOrders()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error fetchinh orders")
		return
	}
	username := GetSessionString(c, "username")

	if username == "" {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	c.HTML(http.StatusOK, "admin.tmpl", AdminDashaboardData{
		Orders: orders,
		Statuses: models.OrderStatuses,
		Username: username,
	})
}

func (h *Handler) HandleOrderPut(c *gin.Context){
	orderId := c.Param("id")
	newStatus := c.PostForm("status")

	if err := h.orders.UpdateOrderStatus(orderId, newStatus); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	h.notificationManager.Notify("order:"+orderId, "order_updated")

	c.Redirect(http.StatusSeeOther, "/admin")
}

func (h *Handler) HandleOrderDelete(c *gin.Context){
	orderID := c.Param("id")

	if err := h.orders.DeleteOrder(orderID); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin")
}