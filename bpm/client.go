package bpm

import (
	"github.com/lanemets/claim-funnel/external/claim"
	"github.com/lanemets/claim-funnel/task"
	"log"
	"path/filepath"
)

const (
	ClaimProcessKey            = "claim-process"
	NotifyBeneficiaryTopicName = "notify-beneficiary"
	ClaimConfirmTaskId         = "claim-confirm"
)

type WorkerConfig struct {
	topicName    string
	retries      int
	retryTimeout int
}

type BpmClient interface {
	deployProcess(path string) error
	startProcessInstance(processKey string, businessKey string) (error, *ProcessDefinitionId)
	startExternalTaskWorker(config *WorkerConfig, consumer task.Consumer) error
	completeUserTask(businessKey string, processDefinitionId *ProcessDefinitionId) error
}

type ClaimBpmClient struct {
	bpmClient BpmClient
}

type Result struct {
	processId string
}

func NewClaimBpmClient() *ClaimBpmClient {
	return &ClaimBpmClient{
		bpmClient: NewCamundaClient(
			Credentials{
				endpointUrl: "http://localhost:8080/engine-rest",
				user:        "demo",
				password:    "demo",
			},
		),
	}
}

func (client ClaimBpmClient) DeployClaimProcess() {
	absPath, pathErr := filepath.Abs("bpm/resources/claim-process.bpmn")
	if pathErr != nil {
		log.Fatalf("unable to find process file, error: %v", pathErr)
	}

	err := client.bpmClient.deployProcess(absPath)
	if err != nil {
		log.Fatalf("unable to deploy claim process, error: %v", err)
	}
}

func (client ClaimBpmClient) StartClaimProcessInstance(claimId claim.Id) *ProcessDefinitionId {
	err, processDefinitionId := client.bpmClient.startProcessInstance(ClaimProcessKey, claimId.Value)
	if err != nil {
		log.Fatalf("unable to start claim process instance; claimId: %s, error: %s", claimId.Value, err)
	}

	return processDefinitionId
}

func (client ClaimBpmClient) StartBeneficiariesNotification(consumer task.Consumer) {
	config := &WorkerConfig{
		topicName:    NotifyBeneficiaryTopicName,
		retries:      5,
		retryTimeout: 5_000,
	}
	_ = client.bpmClient.startExternalTaskWorker(config, consumer)
}

func (client ClaimBpmClient) CompleteClaimConfirmTask(claimId claim.Id, processDefinitionId *ProcessDefinitionId) {
	err := client.bpmClient.completeUserTask(
		claimId.Value,
		processDefinitionId,
	)
	if err != nil {
		log.Fatalf("unable to complete claim confirm task; claimId: %s, error: %s", claimId.Value, err)
	}
}
