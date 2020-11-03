package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lanemets/claim-funnel/bpm"
	"github.com/lanemets/claim-funnel/config"
	"github.com/lanemets/claim-funnel/external/claim"
	cb "github.com/lanemets/claim-funnel/external/claim/gen"
	"github.com/lanemets/claim-funnel/handler"
	"time"
)

func main() {
	createClaimRequest := &cb.CreateClaimRequest{
		Claim: &cb.Claim{
			Email:                 "lanemets.vv+0lo@gmail.com",
			Amount:                "100.0",
			CurrencyCode:          "USD",
			Description:           "Test Claim To Create",
			ClientReferenceNumber: uuid.New().String(),
		},
		Profile: &cb.Profile{
			ExternalId:  "0",
			ProfileType: 0,
			Entity: &cb.Profile_Person{
				Person: &cb.Person{
					FirstName: "TestFirsdt0",
					LastName:  "TestLastd0",
				},
			},
		},
	}

	bpmClient := bpm.NewClient(
		&config.BpmClient{
			Credentials: &config.BpmCredentials{
				EndpointUrl: "http://localhost:8080/engine-rest",
				User:        "demo",
				Password:    "demo",
			},
		},
	)

	//bpmClient.DeployProcess(
	//	&config.BpmProcess{
	//		FilePath: "bpm/resources/claim-process.bpmn",
	//	},
	//)

	// setup service task handlers
	const profileService = "localhost:9001"
	const claimService = "localhost:9002"

	bpmClient.RegisterServiceHandlers(
		&config.Worker{
			MaxTasks:                  5,
			Retries:                   1,
			RetryTimeoutMillis:        5_000,
			LockDuration:              time.Second * 5,
			MaxParallelTaskPerHandler: 1,
			LongPollingTimeout:        time.Second * 1,
		},

		handler.ServiceHandler{
			Handler:  handler.NotifyBeneficiary(profileService),
			Topic:    config.NotifyBeneficiaryTopicName,
			WorkerId: "notify-beneficiary",
		},
		handler.ServiceHandler{
			Handler:  handler.GetProfileByEmail(profileService),
			Topic:    config.GetProfileByEmailTopicName,
			WorkerId: "get-profile-by-email",
		},
		handler.ServiceHandler{
			Handler:  handler.GetClaimInfo(claimService),
			Topic:    config.GetClaimTopicName,
			WorkerId: "get-claim",
		},
		handler.ServiceHandler{
			Handler:  handler.AcknowledgeClaim(profileService),
			Topic:    config.AcknowledgeClaimTopicName,
			WorkerId: "acknowledge-claim",
		},
		handler.ServiceHandler{
			Handler:  handler.SetPaymentPending(profileService),
			Topic:    config.SetPaymentPendingTopicName,
			WorkerId: "set-payment-pending",
		},
	)

	// clients for services
	claimClient := claim.NewClient()
	//profileClient := profile.NewClient()

	// create claim before process start
	claimId := claimClient.Create(createClaimRequest)

	// start bpm process instance
	bpmClient.StartProcessInstance(claimId)

	// confirm claim
	//profileClient.ConfirmClaim(claimId)
	//log.Printf("Claim: %s has been confirmed sucessfully", claimId)

	//bpmClient.CompleteClaimConfirmTask(claimId, processDefinitionId)
	//log.Printf("BPM Task has been completed; processDefinitionId: %s", processDefinitionId)

	engine := gin.Default()
	_ = engine.Run(":8090")
}
