package bpm

import (
	"errors"
	"fmt"
	"github.com/lanemets/claim-funnel/model"
	"github.com/lanemets/claim-funnel/usecases"
	"log"
	"path/filepath"
)

type BpmClient interface {
	RegisterExternalTaskWorker(workerConfig *WorkerConfig, handler ServiceTaskHandler) error
	DeployProcess(path string) error
	StartProcessInstance(processKey string, businessKey string) (*usecases.ProcessDefinitionId, error)
	CompleteUserTask(businessKey string, taskId string, processDefinitionId *usecases.ProcessDefinitionId) error
}

type BpmClaimClient struct {
	client BpmClient
}

func NewBpmClaimClient(client BpmClient) BpmClaimClient {
	return BpmClaimClient{
		client: client,
	}
}

func (client BpmClaimClient) DeployProcess(processConfig *Process) {
	absPath, pathErr := filepath.Abs(processConfig.FilePath)
	if pathErr != nil {
		log.Fatalf("unable to find process file, error: %v", pathErr)
	}

	err := client.client.DeployProcess(absPath)
	if err != nil {
		log.Fatalf("unable to deploy claim process, error: %v", err)
	}
}

func (client BpmClaimClient) StartProcessInstance(claimId *model.ClaimId) (*usecases.ProcessDefinitionId, error) {
	processDefinitionId, err := client.client.StartProcessInstance(ClaimProcessKey, claimId.Value)
	if err != nil {
		errMsg := fmt.Sprintf("unable to start claim process instance; claimId: %s, error: %s", claimId.Value, err)
		log.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	return processDefinitionId, nil
}

func (client BpmClaimClient) RegisterServiceHandlers(workerConfig *WorkerConfig, handlers ...ServiceTaskHandler) error {
	for _, h := range handlers {
		err := client.client.RegisterExternalTaskWorker(workerConfig, h)
		if err != nil {
			errMsg := fmt.Sprintf("an error has occurred during external task worker registration; err: %v", err)
			log.Println(errMsg)
			return errors.New(errMsg)
		}
	}
	return nil
}

func (client BpmClaimClient) CompleteClaimConfirmTask(
	claimId *model.ClaimId,
	processDefinitionId *usecases.ProcessDefinitionId,
) {
	err := client.client.CompleteUserTask(
		claimId.Value,
		ClaimConfirmTaskId,
		processDefinitionId,
	)
	if err != nil {
		log.Fatalf("unable to complete claim confirm task; claimId: %s, error: %s", claimId.Value, err)
	}
}
