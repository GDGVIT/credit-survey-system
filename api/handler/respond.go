package handler

import (
	"CreditBasedSurvey/model"
	"CreditBasedSurvey/pkg/activity"
	form2 "CreditBasedSurvey/pkg/form"
	"CreditBasedSurvey/pkg/response"
	usr "CreditBasedSurvey/pkg/user"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber"
	"math"
	"time"
)

func GetFormForResponse(ctx *fiber.Ctx) {
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
		ctx.Status(fiber.StatusNotFound)
		ctx.SendString("Form not found in our database.")
		return
	}
	if !formMap.IsPublished {
		ctx.Status(fiber.StatusNotFound)
		ctx.SendString("Form is not published.")
		return
	}

	if formMap.UserId == id.(string) {
		ctx.Status(fiber.StatusForbidden)
		ctx.SendString("Author of a survey can't respond to their own surveys.")
		return
	}

	authorizedRes := int64(math.Floor(formMap.CreditsAllotted * formMap.ResponseRate))
	if authorizedRes < formMap.NumResp {
		ctx.Status(fiber.StatusForbidden)
		ctx.SendString("Credits exhausted")
		return
	}

	quests, err := form2.GetQuestions(formid)
	if err != nil {
		ctx.Status(fiber.StatusNotFound)
		ctx.SendString("Error finding questions for the form.")
		return
	}
	surveyDetails := fiber.Map{}
	if formMap.ShowEmail {
		user, _ := usr.GetUser(formMap.UserId)
		surveyDetails["AuthorEmail"] = user.Email
		surveyDetails["AuthorName"] = user.Name
	}
	surveyDetails["Title"] = formMap.Title
	surveyDetails["Description"] = formMap.Description
	surveyDetails["Expiry"] = formMap.Expiry
	surveyDetails["IsAnonymous"] = formMap.IsAnonymous

	model, err := response.GetResponse(id.(string))
	if err == nil {
		ctx.Status(fiber.StatusNotAcceptable)
		ctx.JSON(fiber.Map{
			"Message":   "Already filled the form",
			"Responses": model,
			"Details":   surveyDetails,
			"body":      quests,
		})
		return
	}
	payload := fiber.Map{
		"Details": surveyDetails,
		"body":    quests,
	}

	ctx.JSON(payload)
}


func PostFormResponse(ctx *fiber.Ctx) {
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
		ctx.SendString("Error fetching form from out database")
		return
	}
	if !formMap.IsPublished {
		ctx.Status(fiber.StatusNotFound)
		ctx.SendString("Form is not published.")
		return
	}

	if formMap.UserId == id.(string) {
		ctx.Status(fiber.StatusForbidden)
		ctx.SendString("Author of a survey can't respond to their own surveys.")
		return
	}
	if formMap.Expiry.Before(time.Now()) {
		ctx.Status(fiber.StatusForbidden)
		ctx.SendString("Form Expired")
		return
	}
	authorizedRes := int64(math.Floor(formMap.CreditsAllotted * formMap.ResponseRate))
	if authorizedRes < formMap.NumResp {
		ctx.Status(fiber.StatusForbidden)
		ctx.SendString("Credits exhausted")
		return
	}
	ans := model.NewResponse()
	err = json.Unmarshal([]byte(ctx.Body()), &ans)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest)
		ctx.SendString(err.Error())
		return
	}
	ans.FormId = formid
	ans.UserId = id.(string)
	err = response.PostResponse(*ans)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest)
		ctx.SendString("Already Filled the Form.")
		return
	}
	err = form2.IncrementNumberOfResponsesByFormID(formid, 1) // Increasing the form response count
	activity.AddActivity(
		*model.NewActivity("TOOK A SURVEY",formMap.GainRate, formMap.FormId),activity.Responded, id.(string))
	ctx.Status(fiber.StatusCreated)
	ctx.SendString("OK")
}

