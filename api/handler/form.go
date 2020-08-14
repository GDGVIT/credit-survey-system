package handler

import (
	"CreditBasedSurvey/model"
	"CreditBasedSurvey/pkg/activity"
	form2 "CreditBasedSurvey/pkg/form"
	"CreditBasedSurvey/pkg/sanitizer"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber"
)

func CreateNewForm(ctx *fiber.Ctx) {
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["id"]
	if id == nil {
		ctx.Status(fiber.StatusBadRequest)
		ctx.SendString("Malformed JWT")
		return
	}

	body := []byte(ctx.Body())

	var form model.Form
	err := json.Unmarshal(body, &form)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest)
		ctx.SendString(err.Error())
		return
	}
	payloadStatus := sanitizer.CheckFormDetailsPayload(&form)
	if payloadStatus {
		ctx.Status(fiber.StatusBadRequest)
		ctx.JSON(fiber.Map{
			"status": "invalid inputs",
			"message": []string{
				"title may be missing",
				"description may be missing",
				"credits may be negitive or zero",
			},
		})
		return
	}
	sanitizer.SanitizeFormDetails(id.(string),&form)

	res, err := form2.CreateForms(form)
	if err != nil {
		ctx.Status(fiber.StatusBadGateway)
		ctx.SendString(err.Error())
		return
	}
	if err != nil {
		ctx.Status(fiber.StatusBadGateway)
		ctx.SendString(err.Error())
		return
	}
	ctx.Status(fiber.StatusCreated)
	ctx.JSON(fiber.Map{
		"FormId": res,
	})
}

func GetAuthorizedForm(ctx *fiber.Ctx) {
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
	formMap, err := form2.GetFormDetails(formid)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest)
		ctx.SendString(err.Error())
		return
	}
	ctx.JSON(formMap)
}

func PublishForm(ctx *fiber.Ctx) {
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
	formMap, err := form2.GetFormDetails(formid)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest)
		ctx.SendString(err.Error())
		return
	}
	creditsOfUser, err := activity.GetCredits(id.(string))
	if err != nil {
		ctx.Status(fiber.StatusBadRequest)
		ctx.SendString("USER NOT FOUND!")
		return
	}
	if formMap.CreditsAllotted > creditsOfUser {
		ctx.Status(fiber.StatusBadRequest)
		ctx.SendString("You don't have enough credits.")
		return
	}
	if !formMap.IsPublished {
		activity.AddActivity(*model.NewActivity("CREATED A FORM", -1.0*formMap.CreditsAllotted, formid),
			activity.Created, id.(string))
		form2.SetPublishStatus(formid, true)
		ctx.Status(fiber.StatusOK)
		ctx.SendString("PUBLISHED")
		return
	} else {
		ctx.Status(fiber.StatusBadRequest)
		ctx.SendString("ALREADY PUBLISHED!")
		return
	}

	ctx.Status(fiber.StatusFailedDependency)
	ctx.SendString("Error Occurred.")
}

func UnPublishForm(ctx *fiber.Ctx) {
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
	formMap, err := form2.GetFormDetails(formid)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest)
		ctx.SendString(err.Error())
		return
	}
	creditsOfUser, err := activity.GetCredits(id.(string))
	if err != nil {
		ctx.Status(fiber.StatusBadRequest)
		ctx.SendString("USER NOT FOUND!")
		return
	}
	if formMap.CreditsAllotted > creditsOfUser {
		ctx.Status(fiber.StatusBadRequest)
		ctx.SendString("You don't have enough credits.")
		return
	}
	if formMap.IsPublished {
		amount := formMap.CreditsAllotted - (float64(formMap.NumResp) / formMap.ResponseRate)
		activity.AddActivity(*model.NewActivity("Reverse transfer because you unpublished a form",
			amount, formid),
			activity.ReverseTransfer, id.(string))
		form2.SetPublishStatus(formid, false)
		ctx.Status(fiber.StatusOK)
		ctx.SendString("PUBLISHED")
		return
	} else {
		ctx.Status(fiber.StatusBadRequest)
		ctx.SendString("ALREADY UNPUBLISHED!")
		return
	}
	ctx.Status(fiber.StatusFailedDependency)
	ctx.SendString("Error Occurred.")
}