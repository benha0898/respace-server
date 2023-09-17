package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/account"
	"github.com/stripe/stripe-go/v75/accountlink"
)

type PrefilledInfo struct {
	RespaceId string                 `json:"respaceId"`
	Email     string                 `json:"email"`
	FirstName string                 `json:"firstName"`
	LastName  string                 `json:"lastName"`
	Phone     string                 `json:"phoneNumber"`
	Birthday  stripe.PersonDOBParams `json:"birthday"`
	Address   stripe.AddressParams   `json:"address"`
}

func main() {
	port := "1235"

	stripe.Key = "sk_test_51NCVqgDsu5gTqmrhhu88xpcL1c3QquxbGhBNP7TFePVDVcSkzb5zySvkGzrjiGGsQazKDZJ0dqJ7k2DuuAa3SEUb00FaMWZ1jW"

	// Set up API endpoints
	api := gin.Default()

	api.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	api.POST("/createAccount", createAccount)

	api.Run(":" + port)
}

func createAccount(c *gin.Context) {
	// If there is any prefilled info, use properties inside AccountParams
	var prefilledInfo PrefilledInfo
	if err := c.ShouldBindJSON(&prefilledInfo); err != nil {
		raiseError(c, err, 500, "Error parsing request body")
	}
	params := &stripe.AccountParams{
		Type:         stripe.String(string(stripe.AccountTypeExpress)),
		BusinessType: stripe.String(string(stripe.AccountBusinessTypeIndividual)),
		Email:        &prefilledInfo.Email,
		Individual: &stripe.PersonParams{
			FirstName: &prefilledInfo.FirstName,
			LastName:  &prefilledInfo.LastName,
			Email:     &prefilledInfo.Email,
			Phone:     &prefilledInfo.Phone,
			DOB:       &prefilledInfo.Birthday,
			Address:   &prefilledInfo.Address,
			Metadata: map[string]string{
				"respaceId": prefilledInfo.RespaceId,
			},
		},
	}
	resultAccount, err := account.New(params)
	fmt.Println(*resultAccount)

	if err != nil {
		raiseError(c, err, 500, "Error creating a Connect Stripe account")
	} else {
		accountLinkParams := &stripe.AccountLinkParams{
			Account:    stripe.String(resultAccount.ID),
			ReturnURL:  stripe.String("https://google.com"),
			RefreshURL: stripe.String("https://youtube.com"),
			Type:       stripe.String("account_onboarding"),
		}
		resultAccountLink, err := accountlink.New(accountLinkParams)

		if err != nil {
			raiseError(c, err, 500, "Error creating an Account Link")
		} else {
			fmt.Println(*resultAccountLink)
			c.JSON(200, resultAccountLink)
		}
	}
}

func raiseError(c *gin.Context, err error, statusCode int, message string) {
	c.Error(err)
	c.AbortWithStatusJSON(statusCode, gin.H{
		"error":  message + ": " + err.Error(),
		"status": statusCode,
	})
}
