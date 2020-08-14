package handler

import (
	"CreditBasedSurvey/model"
	form2 "CreditBasedSurvey/pkg/form"
	"CreditBasedSurvey/pkg/sanitizer"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
)

func PostFormQuestions(ctx *fiber.Ctx) {
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["id"]
	if id == nil {
		ctx.Status(fiber.StatusBadRequest)
		ctx.SendString("Malformed JWT")
		return
	}
	formid := ctx.Params("formid")
	if formid == "" {
		ctx.Status(fiber.StatusBadRequest)
		ctx.SendString("Malformed URI")
		return
	}
	formDetails, err := form2.GetFormDetails(formid)
	if err != nil {
		ctx.Status(fiber.StatusNotFound)
		ctx.SendString("FORM NOT FOUND")
		return
	}

	if formDetails.UserId != id.(string) {
		ctx.Status(fiber.StatusForbidden)
		ctx.SendString("You are not the owner of the form.")
		return
	}

	var m model.FormQuestion
	err = json.Unmarshal([]byte(ctx.Body()), &m)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest)
		ctx.SendString(err.Error())
		return
	}

	e := sanitizer.SanitizeQuestions(&m, formid)
	if e {
		ctx.Status(fiber.StatusBadRequest)
		ctx.SendString("Malformed Payload.")
		return
	}
	err = form2.PostQuestions(m)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest)
		ctx.SendString("Error saving in the database")
		return
	}

	ctx.Status(fiber.StatusCreated)
	ctx.SendString("OK")
}

func GetFormQuestions(ctx *fiber.Ctx) {
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["id"]
	if id == nil {
		ctx.Status(fiber.StatusBadRequest)
		ctx.SendString("Malformed JWT")
		return
	}
	formid := ctx.Params("formid")
	if formid == "" {
		ctx.Status(fiber.StatusBadRequest)
		ctx.SendString("Malformed URI")
		return
	}
	quests, err := form2.GetQuestions(formid)
	if err != nil {
		ctx.Status(fiber.StatusNotFound)
		ctx.SendString("Error finding questions for the form.")
	}
	ctx.JSON(bson.M{
		"id":   formid,
		"data": quests,
	})
}

func PostUpdatedFormQuestions(ctx *fiber.Ctx) {
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["id"]
	if id == nil {
		ctx.Status(fiber.StatusBadRequest)
		ctx.SendString("Malformed JWT")
		return
	}
	formid := ctx.Params("formid")
	if formid == "" {
		ctx.Status(fiber.StatusBadRequest)
		ctx.SendString("Malformed URI")
		return
	}

	formDetails, err := form2.GetFormDetails(formid)
	if err != nil {
		ctx.Status(fiber.StatusNotFound)
		ctx.SendString("FORM NOT FOUND")
		return
	}

	if formDetails.UserId != id.(string) {
		ctx.Status(fiber.StatusForbidden)
		ctx.SendString("You are not the owner of the form.")
		return
	}

	var m model.FormQuestion
	err = json.Unmarshal([]byte(ctx.Body()), &m)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest)
		ctx.SendString(err.Error())
		return
	}

	e := sanitizer.SanitizeQuestions(&m, formid)
	if e {
		ctx.Status(fiber.StatusBadRequest)
		ctx.SendString("Malformed Payload.")
		return
	}
	err = form2.ReplaceQuestions(m)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest)
		ctx.SendString("Error saving in the database")
		return
	}

	ctx.Status(fiber.StatusCreated)
	ctx.SendString("OK")
}