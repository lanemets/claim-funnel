package usecases

import (
	"errors"
	"fmt"
	"github.com/lanemets/claim-funnel/model"
	"log"
)

type Interactor struct {
	ClaimClient   RpcClaim
	ProfileClient RpcProfile
	BpmClient     BpmClaimClient
}

func (s *Interactor) CreateClaim(claim *model.Claim, profile *model.Profile) (*model.ClaimId, *model.ProcessDefinitionId, error) {
	claimId, err := s.ClaimClient.Create(claim, profile)
	if err != nil {
		errMsg := fmt.Sprintf("an error occurred on claim creation; claimId: %v, error: %v", claimId, err)
		log.Println(errMsg)
		return nil, nil, errors.New(errMsg)
	}

	processId, err := s.BpmClient.StartProcessInstance(claimId)
	if err != nil {
		errMsg := fmt.Sprintf("an error occurred on claim creation; claimId: %v, error: %v", claimId, err)
		log.Println(errMsg)
		return nil, nil, errors.New(errMsg)
	}

	return claimId, processId, nil
}

func (s *Interactor) ConfirmClaim(claimId *model.ClaimId) error {
	return s.ClaimClient.ConfirmClaim(claimId)
}

func NewInteractor(
	claimClient RpcClaim,
	profileClient RpcProfile,
	bpmClient BpmClaimClient,
) *Interactor {
	return &Interactor{
		ClaimClient:   claimClient,
		ProfileClient: profileClient,
		BpmClient:     bpmClient,
	}
}
