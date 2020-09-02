package handler

import (
	"CreditBasedSurvey/model"
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

func GetUserFormsByID(ctx *fiber.Ctx)  {
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


	var forms = new([]model.Form)
	err = usr.GetAllFormsOfUserID(id.(string), forms)
	if err != nil {
		ctx.Status(fiber.StatusNotFound)
		ctx.SendString("Error fetching the forms of the current user.")
		return
	}


	ctx.Status(fiber.StatusOK)
	ctx.JSON(fiber.Map{
		"status": "User Found.",
		"userDetails":   me,
		"forms": *forms,
	})
}