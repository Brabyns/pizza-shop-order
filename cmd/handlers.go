package main

import (
	 "github.com/Brabyns/pizza-shop-order/internal/models"
)

type Handler struct {
	orders *models.OrderModel
	users *models.UserModel
}

func NewHandler(dbModel *models.DBModel) *Handler{
	return &Handler{
		orders: &dbModel.Order,
		users: &dbModel.User,
	}
}