package accounts

import (
	"backend/models"
	"backend/utils"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
)

func postBio(c fiber.Ctx) error {
	var body struct {
		Bio string `json:"bio"`
	}
	json.Unmarshal(c.Body(), &body)

	account := models.Account{}
	utils.GetLocals(c, "account", &account)

	err := models.UpdateAccount(
		account.ID,
		bson.M{
			"bio": body.Bio,
		},
	)
	if err != nil {
		return utils.Error(c, errors.New("Could not update bio"))
	}

	account.Bio = body.Bio
	token := account.GenAccountToken()

	return c.JSON(bson.M{
		"account": account,
		"token":   token,
	})
}

func patchSocial(c fiber.Ctx) error {
	var body models.Social
	json.Unmarshal(c.Body(), &body)

	account := models.Account{}
	utils.GetLocals(c, "account", &account)

	exists := false
	for i, s := range account.Socials {
		if s.Platform == body.Platform {
			exists = true
			account.Socials[i].Link = body.Link
		}
	}

	if !exists {
		account.Socials = append(account.Socials, body)
	}

	err := models.UpdateAccount(
		account.ID,
		bson.M{
			"socials": account.Socials,
		},
	)
	if err != nil {
		return utils.Error(c, errors.New("Could not update bio"))
	}

	token := account.GenAccountToken()

	return c.JSON(bson.M{
		"account": account,
		"token":   token,
	})
}

func postFollow(c fiber.Ctx) error {
	var body struct {
		FollowID string `json:"followID"`
	}
	json.Unmarshal(c.Body(), &body)

	// account from token
	account := models.Account{}
	utils.GetLocals(c, "account", &account)

	// get account again for good measure
	account, err := models.GetAccount(
		bson.M{
			"id": account.ID,
		},
	)
	if err != nil {
		return utils.Error(c, err)
	}

	// append the following
	exists := false
	for _, accountID := range account.Following {
		if accountID == body.FollowID {
			exists = true
		}
	}
	if !exists {
		account.Following = append(account.Following, body.FollowID)
	}

	// update the account
	err = models.UpdateAccount(
		account.ID,
		bson.M{
			"following": account.Following,
		},
	)
	if err != nil {
		return utils.Error(c, errors.New("Could not follow user"))
	}

	// get the following account
	following, err := models.GetAccount(
		bson.M{
			"id": body.FollowID,
		},
	)
	if err != nil {
		return utils.Error(c, errors.New("Could not get the other user"))
	}
	fmt.Println(following)

	// append to the followers of the following
	exists = false
	for _, accountID := range following.Followers {
		if accountID == account.ID {
			exists = true
		}
	}
	if !exists {
		following.Followers = append(following.Followers, account.ID)
	}

	// update the following account
	err = models.UpdateAccount(
		following.ID,
		bson.M{
			"followers": following.Followers,
		},
	)
	if err != nil {
		return utils.Error(c, errors.New("Could not follow the other user"))
	}

	fmt.Println(following)

	token := account.GenAccountToken()

	return c.JSON(bson.M{
		"token":   token,
		"account": account,
	})
}

func postUnfollow(c fiber.Ctx) error {
	var body struct {
		FollowID string `json:"followID"`
	}
	json.Unmarshal(c.Body(), &body)

	// account from token
	account := models.Account{}
	utils.GetLocals(c, "account", &account)

	// get account again for good measure
	account, err := models.GetAccount(
		bson.M{
			"id": account.ID,
		},
	)
	if err != nil {
		return utils.Error(c, err)
	}

	// remove from the following
	newFollowing := []string{}
	for _, accountID := range account.Following {
		if accountID != body.FollowID {
			newFollowing = append(newFollowing, accountID)
		}
	}

	// update the account
	err = models.UpdateAccount(
		account.ID,
		bson.M{
			"following": newFollowing,
		},
	)
	if err != nil {
		return utils.Error(c, errors.New("Could not follow user"))
	}
	account.Following = newFollowing

	// get the following account
	following, err := models.GetAccount(
		bson.M{
			"id": body.FollowID,
		},
	)
	if err != nil {
		return utils.Error(c, errors.New("Could not get the other user"))
	}
	fmt.Println(following)

	// remove from the followers of the following
	newFollowers := []string{}
	for _, accountID := range following.Followers {
		if accountID != account.ID {
			newFollowers = append(newFollowers, accountID)
		}
	}

	// update the following account
	err = models.UpdateAccount(
		following.ID,
		bson.M{
			"followers": newFollowers,
		},
	)
	if err != nil {
		return utils.Error(c, errors.New("Could not follow the other user"))
	}
	following.Followers = newFollowers

	fmt.Println(following)

	token := account.GenAccountToken()

	return c.JSON(bson.M{
		"token":   token,
		"account": account,
	})
}
