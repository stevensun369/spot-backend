package models

import (
	"backend/db"
	"backend/env"
	"backend/utils"
	"strings"
	"time"

	sj "github.com/brianvoe/sjwt"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
)

type Social struct {
	Platform string `bson:"platform" json:"platform"`
	Link     string `bson:"link" json:"link"`
}

type Account struct {
	ID       string `bson:"id" json:"id"`
	Username string `bson:"username" json:"username"`

	Phone     string `bson:"phone" json:"phone"`
	FirstName string `bson:"firstName" json:"firstName"`
	LastName  string `bson:"lastName" json:"lastName"`

	Bio     string   `bson:"bio" json:"bio"`
	Socials []Social `bson:"socials" json:"socials"`

	Following []string `bson:"following" json:"following"`
	Followers []string `bson:"followers" json:"followers"`
}

func (account Account) GenStrippedToken() string {
	acc := struct {
		Phone string `bson:"phone" json:"phone"`
		ID    string `bson:"id" json:"id"`
	}{
		Phone: account.Phone,
		ID:    account.ID,
	}
	claims, _ := sj.ToClaims(acc)

	claims.SetExpiresAt(time.Now().Add(365 * 24 * time.Hour))
	token := claims.Generate(env.JWTKey)
	return token
}

func (account Account) GenAccountToken() string {
	claims, _ := sj.ToClaims(account)
	claims.SetExpiresAt(time.Now().Add(365 * 24 * time.Hour))
	// 1 year = 365 days * 24 hours in a day

	token := claims.Generate(env.JWTKey)
	return token
}

func ParseAccountToken(token string) (Account, error) {
	hasVerified := sj.Verify(token, env.JWTKey)

	if !hasVerified {
		return Account{}, nil
	}

	claims, _ := sj.Parse(token)
	err := claims.Validate()
	account := Account{}
	claims.ToStruct(&account)

	return account, err
}

func AccountMiddleware(c fiber.Ctx) error {
	var token string

	authHeader := c.Get("Authorization")

	if string(authHeader) != "" && strings.HasPrefix(string(authHeader), "Bearer") {

		tokens := strings.Fields(string(authHeader))
		if len(tokens) == 2 {
			token = tokens[1]
		}
		if token == "" {
			return utils.MessageError(c, "no token")
		}

		account, err := ParseAccountToken(token)
		if err != nil {
			return utils.MessageError(c, "a aparut o eroare")
		}

		c.Locals("id", account.ID)
		utils.SetLocals(c, "account", account)
	}

	if token == "" {
		return utils.MessageError(c, "no token")
	}

	return c.Next()
}

func (account *Account) Create() error {
	// generating ID
	account.ID = utils.GenID(6)

	account.Socials = []Social{}
	account.Following = []string{}
	account.Followers = []string{}

	// creating account
	_, err := db.Accounts.InsertOne(db.Ctx, account)

	return err
}

func UpdateAccount(id string, updates any) error {
	_, err := db.Accounts.UpdateOne(
		db.Ctx,
		bson.M{"id": id},
		bson.M{
			"$set": updates,
		},
	)

	return err
}

func GetAccount(query any) (Account, error) {
	var account Account

	err := db.Accounts.FindOne(
		db.Ctx,
		query,
	).Decode(&account)

	return account, err
}

func CheckAccount(phone string) (bool, Account) {
	var account Account

	err := db.Accounts.FindOne(
		db.Ctx, bson.M{
			"phone": phone,
		},
	).Decode(&account)

	if err != nil {
		return false, Account{}
	} else {
		return true, account
	}
}

func CheckAccountUsername(username string) (bool, Account) {
	var account Account

	err := db.Accounts.FindOne(
		db.Ctx, bson.M{
			"username": username,
		},
	).Decode(&account)

	if err != nil {
		return false, Account{}
	} else {
		return true, account
	}
}
