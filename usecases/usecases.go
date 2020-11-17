package usecases

import (
	"errors"
	"fmt"
	"github.com/lanemets/claim-funnel/model"
	"log"
)

type Interactor struct {
	RpcClaim   RpcClaim
	RpcProfile RpcProfile
	BpmClient  BpmClaimClient
}

func (s *Interactor) CreateClaim(claim *model.Claim, profile *model.Profile) (*model.ClaimId, *model.ProcessDefinitionId, error) {
	claimId, err := s.RpcClaim.Create(claim, profile)
	if err != nil {
		errMsg := fmt.Sprintf("an error occurred on claim creation; error: %v", err)
		log.Println(errMsg)
		return nil, nil, errors.New(errMsg)
	}

	//TODO: move to goroutine?
	processId, err := s.BpmClient.StartProcessInstance(claimId)
	if err != nil {
		errMsg := fmt.Sprintf("an error occurred on starting process; claimId: %v, error: %v", claimId.Value, err)
		log.Println(errMsg)
		return claimId, nil, errors.New(errMsg)
		//TODO: put claim to retry queue for further processing?
	}

	log.Printf("process has been started successfully; claimId: %v, processId: %v\n", claimId.Value, processId)

	return claimId, processId, nil
}

func (s *Interactor) ConfirmClaim(claimId *model.ClaimId) error {
	return s.RpcClaim.ConfirmClaim(claimId)
}

func NewInteractor(
	claimClient RpcClaim,
	profileClient RpcProfile,
	bpmClient BpmClaimClient,
) *Interactor {
	return &Interactor{
		RpcClaim:   claimClient,
		RpcProfile: profileClient,
		BpmClient:  bpmClient,
	}
}
