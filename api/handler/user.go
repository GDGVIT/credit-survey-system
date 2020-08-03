package handler

import (
	usr "CreditBasedSurvey/pkg/user"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber"
)

func GetUser(ctx *fiber.Ctx) {
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["id"]
	if id == nil {
		ctx.Status(fiber.StatusBadRequest)
		ctx.SendString("Malformed JWT")
		return
	}
	me, err := usr.GetUser(id.(string))
	if err != nil {
		ctx.Status(fiber.StatusNotFound)
		ctx.SendString("No user found for the user")
		return
	}

	ctx.Status(fiber.StatusOK)
	ctx.JSON(fiber.Map{
		"status": "User Found.",
		"body":   me,
	})

}
