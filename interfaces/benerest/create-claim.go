package benerest

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lanemets/claim-funnel/model"
	"github.com/stroiman/go-automapper"
	"log"
	"net/http"
)

func CreateClaim(interactor Interactor) func(c *gin.Context) {
	return func(c *gin.Context) {
		var request CreateClaimRequest
		err := c.ShouldBindJSON(&request)

		if err != nil {
			errMsg := fmt.Sprintf("an error has occurred on parsing request; err: %v", err)
			log.Println(errMsg)

			c.JSON(http.StatusBadRequest, gin.H{"error": errors.New(errMsg)})
		}

		log.Printf("request: %v", request)

		var claim = &model.Claim{}
		automapper.MapLoose(request.Claim, claim)

		log.Printf("claim: %v", claim)

		var profile = &model.Profile{}
		automapper.MapLoose(request.Profile, profile)

		log.Printf("profile: %v", profile)

		claimId, processId, err := interactor.CreateClaim(claim, profile)
		if err != nil {
			errMsg := fmt.Sprintf("error on claim creation; err: %v", err)
			log.Println(errMsg)

			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.New(errMsg)})
		}
		log.Printf("process started: %v", processId)

		c.JSON(200, claimId)
	}
}
