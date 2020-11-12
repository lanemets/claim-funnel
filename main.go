//func main() {
//
//	//register ServiceHandlers
//	//deploy process if needed
//	//start process
//
//	//--claim-create
//	////--start process here
//
//
//	createClaimRequest := &cb.CreateClaimRequest{
//		Claim: &cb.Claim{
//			Email:                 "lanemets.vv+x83j@gmail.com",
//			Amount:                "100.0",
//			CurrencyCode:          "USD",
//			Description:           "Test Claim To Create",
//			ClientReferenceNumber: uuid.New().String(),
//		},
//		Profile: &cb.Profile{
//			ExternalId:  "0",
//			ProfileType: 0,
//			Entity: &cb.Profile_Person{
//				Person: &cb.Person{
//					FirstName: "TestFirsdt0",
//					LastName:  "TestLastd0",
//				},
//			},
//		},
//	}
//
//	bpmClient := bpm.NewClient(
//		&config.BpmClient{
//			Credentials: &config.BpmCredentials{
//				EndpointUrl: "http://localhost:8080/engine-rest",
//				User:        "demo",
//				Password:    "demo",
//			},
//		},
//	)
//
//	//bpmClient.DeployProcess(
//	//	&config.BpmProcess{
//	//		FilePath: "bpm/resources/claim-process.bpmn",
//	//	},
//	//)
//
//	// setup service task handlers
//	const profileService = "localhost:9001"
//	const claimService = "localhost:9002"
//
//	bpmClient.RegisterServiceHandlers(
//		&config.Worker{
//			MaxTasks:                  5,
//			Retries:                   1,
//			RetryTimeoutMillis:        5_000,
//			LockDuration:              time.Second * 5,
//			MaxParallelTaskPerHandler: 1,
//			LongPollingTimeout:        time.Second * 1,
//		},
//
//		handler.ServiceHandler{
//			Handler:  handler.NotifyBeneficiary(claimService),
//			Topic:    config.NotifyBeneficiaryTopicName,
//			WorkerId: "notify-beneficiary",
//		},
//		handler.ServiceHandler{
//			Handler:  handler.GetProfileByEmail(profileService),
//			Topic:    config.GetProfileByEmailTopicName,
//			WorkerId: "get-profile-by-email",
//		},
//		handler.ServiceHandler{
//			Handler:  handler.GetClaimInfo(claimService),
//			Topic:    config.GetClaimTopicName,
//			WorkerId: "get-claim",
//		},
//		handler.ServiceHandler{
//			Handler:  handler.AcknowledgeClaim(claimService),
//			Topic:    config.AcknowledgeClaimTopicName,
//			WorkerId: "acknowledge-claim",
//		},
//		handler.ServiceHandler{
//			Handler:  handler.SetPaymentPending(claimService),
//			Topic:    config.SetPaymentPendingTopicName,
//			WorkerId: "set-payment-pending",
//		},
//	)
//
//	// clients for services
//	claimClient := claim.NewClient()
//	//profileClient := profile.NewClient()
//
//	// create claim before process start
//	claimId := claimClient.Create(createClaimRequest)
//
//	// start bpm process instance
//	bpmClient.StartProcessInstance(claimId)
//
//	// confirm claim
//	//profileClient.ConfirmClaim(claimId)
//	//log.Printf("Claim: %s has been confirmed sucessfully", claimId)
//
//	//bpmClient.CompleteClaimConfirmTask(claimId, processDefinitionId)
//	//log.Printf("BPM Task has been completed; processDefinitionId: %s", processDefinitionId)
//
//	//engine := gin.Default()
//	//_ = engine.Run(":8090")
//}

package main

import (
	"github.com/lanemets/claim-funnel/cmd"
	"github.com/lanemets/claim-funnel/interfaces/benerest"
)

func main() {
	cmd.Execute(
		benerest.NewServer(&benerest.Config{}),
	)
}
