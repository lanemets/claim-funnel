package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lanemets/claim-funnel/bpm"
	"github.com/lanemets/claim-funnel/external/claim"
	"github.com/lanemets/claim-funnel/external/profile"
	"time"
)

func main() {
	//TODO wrap with goroutines
	bpmClient := bpm.NewClaimBpmClient()
	bpmClient.DeployClaimProcess()

	profileClient := profile.Client()
	bpmClient.StartBeneficiariesNotification(profileClient.NotifyBeneficiaryTask())

	claimClient := claim.Client()
	claimId := claimClient.Create()
	processDefinitionId := bpmClient.StartClaimProcessInstance(claimId)

	time.Sleep(time.Second * 5)

	// on claim confirm from mail link
	profileClient.ConfirmClaim(claimId)
	bpmClient.CompleteClaimConfirmTask(claimId, processDefinitionId)

	engine := gin.Default()
	//engine.POST("/v1/claims")

	_ = engine.Run(":8090")
}
