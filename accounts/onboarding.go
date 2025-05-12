package accounts

import (
	"backend/db"
	"backend/env"
	"backend/models"
	"backend/utils"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func postOnboarding(c fiber.Ctx) error {
	var body struct {
		Phone string
	}
	json.Unmarshal(c.Body(), &body)

	exists, account := models.CheckAccount(body.Phone)

	// generating code
	code := utils.GenCode(4)
	fmt.Println(code)

	// // // sending verification code on sms
	// err := utils.SendSMS(phone, code)
	// if err != nil {
	// 	return utils.MessageError(c, fmt.Sprintf("Codul nu a putut fi trimis.: %v", err))
	// }

	// setting verification code in redis
	hashedCode, err := bcrypt.GenerateFromPassword([]byte(code), 10)
	if err != nil {
		return utils.MessageError(c, "A aparut o problema tehnica, incercati mai tarziu.")
	}
	db.Set("code:"+body.Phone, string(hashedCode))

	if !exists {
		account.Create()
		account.Phone = body.Phone
	}

	token := account.GenStrippedToken()

	if env.DEV {
		return c.JSON(bson.M{
			"phone":     body.Phone,
			"token":     token,
			"newClient": !exists,
			"code":      code,
		})
	} else {
		return c.JSON(bson.M{
			"phone":     body.Phone,
			"token":     token,
			"newClient": !exists,
		})
	}
}

func postVerifyCode(c fiber.Ctx) error {
	var body struct {
		Code      string `json:"code"`
		NewClient bool   `json:"newClient"`
	}
	json.Unmarshal(c.Body(), &body)

	fmt.Println(body)

	// get account from token
	account := models.Account{}
	utils.GetLocals(c, "account", &account)

	fmt.Println(account.Phone)

	// getting code from redis
	hashedCode, _ := db.Get("code:" + account.Phone)

	if bcrypt.CompareHashAndPassword(
		[]byte(hashedCode), []byte(body.Code)) != nil {
		return utils.MessageError(c, "Codul introdus nu este corect")
	}

	// update account with phone
	// if new client
	var err error
	if body.NewClient {
		err = models.UpdateAccount(
			account.ID,
			bson.M{
				"phone": account.Phone,
			},
		)

		account.Socials = []models.Social{}

		if err != nil {
			return utils.Error(c, err)
		}
	} else {
		account, err = models.GetAccount(bson.M{
			"id": account.ID,
		})
		if err != nil {
			return utils.Error(c, err)
		}
	}

	// gen token
	token := account.GenAccountToken()

	return c.JSON(bson.M{
		"account": account,
		"token":   token,
	})
}

func postName(c fiber.Ctx) error {
	var body struct {
		FirstName string
		LastName  string
	}
	json.Unmarshal(c.Body(), &body)

	account := models.Account{}
	utils.GetLocals(c, "account", &account)

	// update local account with name
	account.FirstName = body.FirstName
	account.LastName = body.LastName

	// update account with name
	err := models.UpdateAccount(
		account.ID,
		bson.M{
			"firstName": body.FirstName,
			"lastName":  body.LastName,
		},
	)
	if err != nil {
		return utils.MessageError(c, "Eroare")
	}

	// generating the account token
	token := account.GenAccountToken()

	return c.JSON(bson.M{
		"account": account,
		"token":   token,
	})
}
