package main

import (
	"CreditBasedSurvey/api/handler"
	"CreditBasedSurvey/model"
	usr "CreditBasedSurvey/pkg/user"
	"CreditBasedSurvey/pkg/utils"
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber"
	jwtware "github.com/gofiber/jwt"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var JwtSecret = "itsAnewJWTEverytime"

var conf oauth2.Config

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	googleClientId := os.Getenv("CLIENT_ID")
	googleClientSecret := os.Getenv("CLIENT_SECRET")
	conf = oauth2.Config{
		ClientID:     googleClientId,
		ClientSecret: googleClientSecret,
		RedirectURL:  "http://127.0.0.1:5000/user/auth",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid",
		},
		Endpoint: google.Endpoint}
}

func ErrorOccurred(err error, ctx *fiber.Ctx) {
	ctx.SendStatus(fiber.StatusInternalServerError)
	ctx.SendString(err.Error())
}


func main() {
	err := utils.SetClient()
	if err != nil {
		panic(err.Error())
	}
	defer utils.CloseClient()

	Init()
	app := fiber.New()

	// API Groups
	user := app.Group("/user")
	{
		user.Get("/new", func(ctx *fiber.Ctx) {
			ctx.Status(fiber.StatusTemporaryRedirect)
			ctx.Set("Location", conf.AuthCodeURL("signup"))
		})
		// REMINDER first Login query parameter to be added.
		user.Get("/auth", func(ctx *fiber.Ctx) {
			if ctx.Query("state") != "signup" {
				ctx.Status(fiber.StatusForbidden)
				ctx.SendString("Not a sign-in op")
				return
			} else {
				tok, err := conf.Exchange(context.TODO(), ctx.Query("code"))
				if err != nil {
					ErrorOccurred(err, ctx)
					return
				}
				client := conf.Client(context.TODO(), tok)
				var bodyJson map[string]interface{}
				resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
				if err != nil {
					ErrorOccurred(err, ctx)
					return
				}
				defer resp.Body.Close()
				data, _ := ioutil.ReadAll(resp.Body)
				err = json.Unmarshal(data, &bodyJson)
				if err != nil {
					ErrorOccurred(err, ctx)
					return
				}

				// CREATE a user
				newUser := model.NewUser(bodyJson["email"].(string), bodyJson["name"].(string))
				_id, err := usr.InitializeUser(*newUser, 100)

				// CREATING TOKEN
				token := jwt.New(jwt.SigningMethodHS256)
				claims := token.Claims.(jwt.MapClaims)
				claims["id"] = _id
				claims["owner"] = true
				claims["expiry"] = time.Now().Add(time.Hour * 72).Unix()

				// GENERATING TOKEN

				t, err := token.SignedString([]byte(JwtSecret))
				if err != nil {
					ErrorOccurred(err, ctx)
					return
				}

				ctx.JSON(fiber.Map{
					"jwt": t,
					"tok": tok,
				})
			}
		})
	}

	api := app.Group("/api", jwtware.New(jwtware.Config{
		SigningKey: []byte(JwtSecret),
	}))
	v1 := api.Group("/v1")
	{
		v1.Get("/activity/me", handler.GetUserActivity) // Route to fetch activity
		user := v1.Group("/user")
		{
			user.Delete("/me", handler.DeleteUser) // DELETE current user
			user.Get("/me", handler.GetUser)       // GET current user
			user.Get("/me/forms", handler.GetUserFormsByID) // GET forms of a user
		}

		form := v1.Group("/form")
		{
			// Group of routes for forms
			// Insecure content for the new form
			form.Post("/new", handler.CreateNewForm)
			form.Get("/:formid", handler.GetAuthorizedForm)
			// ROUTE to POST questions of a form
			form.Post("/:formid/questions", handler.PostFormQuestions)
			form.Get("/:formid/questions", handler.GetFormQuestions)
			form.Post("/:formid/questUpdate", handler.PostUpdatedFormQuestions)
			form.Put("/:formid/publish", handler.PublishForm)
			form.Put("/:formid/un-publish", handler.UnPublishForm)
		}

		respond := v1.Group("/respond")
		{
			respond.Get("/:formid", handler.GetFormForResponse)
			respond.Post("/:formid", handler.PostFormResponse)
		}
	}
	app.Listen(os.Getenv("PORT"))
}