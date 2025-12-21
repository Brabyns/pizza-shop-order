package main

import (
	"pizza-tracker-go/internal/models"
	"slices"
	"github.com/gin-gonic/gin/binding"

	"github.com/go-playground/validator/v10"
)

func RegisterCustomValidators() {
	if v, ok := binding.validator.Engine().(*validator.Validate); ok {
		v.RegisterValidattion("valid_pizza_type", createSliceValidator(models.PizzaTypes))
		v.RegisterValidattion("valid_pizza_size", createSliceValidator(models.PizzaSizes))
	}
} 

func createSliceValidator(allowedValues []string) validator.Func{
	return func(fl validator.FieldLevel) bool{
		return slices.Contains(allowedValues, fl.Field().String())
	}
}