package bpm

import (
	"github.com/lanemets/claim-funnel/config"
	"github.com/lanemets/claim-funnel/external/claim"
	"github.com/lanemets/claim-funnel/handler"
	"log"
	"path/filepath"
)

type Client struct {
	client BpmClient
}

type ProcessDefinitionId struct {
	value string
}

type BpmClient interface {
	deployProcess(path string) error
	startProcessInstance(processKey string, businessKey string) (error, *ProcessDefinitionId)

	registerExternalTaskWorker(
		workerConfig *config.Worker,
		serviceHandler handler.ServiceHandler,
	) error

	completeUserTask(businessKey string, taskId string, processDefinitionId *ProcessDefinitionId) error
}

type Credentials struct {
	endpointUrl string
	user        string
	password    string
}

func NewClient(config *config.BpmClient) *Client {
	return &Client{
		client: NewCamundaClient(config.Credentials),
	}
}

func (client Client) DeployProcess(processConfig *config.BpmProcess) {
	absPath, pathErr := filepath.Abs(processConfig.FilePath)
	if pathErr != nil {
		log.Fatalf("unable to find process file, error: %v", pathErr)
	}

	err := client.client.deployProcess(absPath)
	if err != nil {
		log.Fatalf("unable to deploy claim process, error: %v", err)
	}
}

func (client Client) StartProcessInstance(claimId *claim.ClaimId) *ProcessDefinitionId {
	err, processDefinitionId := client.client.startProcessInstance(config.ClaimProcessKey, claimId.Value)
	if err != nil {
		log.Fatalf("unable to start claim process instance; claimId: %s, error: %s", claimId.Value, err)
	}

	return processDefinitionId
}

func (client Client) RegisterServiceHandlers(workerConfig *config.Worker, handlers ...handler.ServiceHandler) {
	for _, h := range handlers {
		_ = client.client.registerExternalTaskWorker(workerConfig, h)
	}
}

func (client Client) CompleteClaimConfirmTask(claimId *claim.ClaimId, processDefinitionId *ProcessDefinitionId) {
	err := client.client.completeUserTask(
		claimId.Value,
		"claim-confirm",
		processDefinitionId,
	)
	if err != nil {
		log.Fatalf("unable to complete claim confirm task; claimId: %s, error: %s", claimId.Value, err)
	}
}
