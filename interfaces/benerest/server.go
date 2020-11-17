package benerest

import (
	"fmt"
	camunda "github.com/citilinkru/camunda-client-go"
	"github.com/gin-gonic/gin"
	"github.com/lanemets/claim-funnel/interfaces/benerpc"
	"github.com/lanemets/claim-funnel/interfaces/bpm"
	"github.com/lanemets/claim-funnel/model"
	"github.com/lanemets/claim-funnel/usecases"
	"log"
	"time"
)

type Server interface {
	Start()
}

type Config struct {
}

type ClaimServer struct {
	//TODO: use config
	config     *Config
	interactor *Interactor
}

type Interactor interface {
	CreateClaim(claim *model.Claim, profile *model.Profile) (*model.ClaimId, *model.ProcessDefinitionId, error)
	ConfirmClaim(claimId *model.ClaimId) error
}

func NewServer() Server {
	return &ClaimServer{}
}

func (*ClaimServer) Start() {

	claimGrpcContext, claimErr := benerpc.NewGrpcContext("localhost:9002")
	if claimErr != nil {
		errMsg := fmt.Sprintf("an error has occurred on grpc context creating: %v", claimErr)
		log.Fatalf(errMsg)
	}
	defer claimGrpcContext.Close()

	profileGrpcContext, profileErr := benerpc.NewGrpcContext("localhost:9001")
	if profileErr != nil {
		errMsg := fmt.Sprintf("an error has occurred on grpc context creating: %v", profileErr)
		log.Fatalf(errMsg)
	}
	defer profileGrpcContext.Close()

	//TODO: move to config
	credentials := &bpm.Credentials{
		EndpointUrl: "http://localhost:8080/engine-rest",
		User:        "demo",
		Password:    "demo",
	}

	bpmClient := bpm.NewBpmClaimClient(
		bpm.NewCamundaClient(
			camunda.NewClient(
				camunda.ClientOptions{
					EndpointUrl: credentials.EndpointUrl,
					ApiUser:     credentials.User,
					ApiPassword: credentials.Password,
					Timeout:     time.Second * 10,
				},
			),
		),
	)

	//TODO: do not deploy process if the same exists
	bpmClient.DeployProcess(
		&bpm.Process{
			FilePath: "interfaces/bpm/resources/claim-process.bpmn",
		},
	)

	registerError := bpmClient.RegisterServiceHandlers(
		&bpm.WorkerConfig{
			MaxTasks:                  5,
			Retries:                   1,
			RetryTimeoutMillis:        5_000,
			LockDuration:              time.Second * 5,
			MaxParallelTaskPerHandler: 1,
			LongPollingTimeout:        time.Second * 1,
		},

		bpm.ServiceTaskHandler{
			Handler:  bpm.NotifyBeneficiary(claimGrpcContext),
			Topic:    bpm.NotifyBeneficiaryTopicName,
			WorkerId: "notify-beneficiary",
		},
		bpm.ServiceTaskHandler{
			Handler:  bpm.GetProfileByEmail(profileGrpcContext),
			Topic:    bpm.GetProfileByEmailTopicName,
			WorkerId: "get-profile-by-email",
		},
		bpm.ServiceTaskHandler{
			Handler:  bpm.GetClaimInfo(claimGrpcContext),
			Topic:    bpm.GetClaimTopicName,
			WorkerId: "get-claim",
		},
		bpm.ServiceTaskHandler{
			Handler:  bpm.AcknowledgeClaim(claimGrpcContext),
			Topic:    bpm.AcknowledgeClaimTopicName,
			WorkerId: "acknowledge-claim",
		},
		bpm.ServiceTaskHandler{
			Handler:  bpm.SetPaymentPending(claimGrpcContext),
			Topic:    bpm.SetPaymentPendingTopicName,
			WorkerId: "set-payment-pending",
		},
	)

	if registerError != nil {
		errMsg := fmt.Sprintf("an error has occurred on in handlers registration: %v", registerError)
		log.Fatalf(errMsg)
	}

	rpcClaim := benerpc.NewRpcClaim(claimGrpcContext)
	profileClient := benerpc.NewRpcProfile(claimGrpcContext)

	interactor := usecases.NewInteractor(rpcClaim, profileClient, bpmClient)

	engine := gin.Default()
	engine.
		GET(
			"/health",
			func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "healthy"})
			}).
		POST(
			"/v1/claims",
			CreateClaim(interactor),
		).
		POST(
			"/v1/claims/:claimId/confirm",
			func(c *gin.Context) {
				c.Param("claimId")
				//confirm claim here
			})

	_ = engine.Run(":8090")
}
