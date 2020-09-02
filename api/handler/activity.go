package handler

import (
	"CreditBasedSurvey/pkg/activity"
	usr "CreditBasedSurvey/pkg/user"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber"
)

func GetUserActivity(ctx *fiber.Ctx) {
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["id"]
	if id == nil {
		ctx.Status(fiber.StatusBadRequest)
		ctx.SendString("Malformed JWT")
		return
	}
	me, err := activity.GetActivityLogs(id.(string))
	if err != nil {
		ctx.Status(fiber.StatusNotFound)
		ctx.SendString("No Logs found for the user")
		return
	}

	ctx.Status(fiber.StatusOK)
	ctx.JSON(fiber.Map{
		"status": "User Found.",
		"body":   me,
	})

}

func DeleteUser(ctx *fiber.Ctx) {
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["id"]
	if id == nil {
		ctx.Status(fiber.StatusBadRequest)
		ctx.SendString("Malformed JWT")
		return
	}
	res, err := usr.DeleteUser(id.(string))
	if err != nil {
		ctx.Status(fiber.StatusNotAcceptable)
		ctx.SendString("UNABLE TO DELETE THE USER")
		return
	}
	if res == 0 {
		ctx.Status(fiber.StatusNotFound)
		ctx.SendString(fmt.Sprint("No user with id:", id))
		return
	}
	ctx.Status(fiber.StatusOK)
	ctx.SendString("DELETED")
}