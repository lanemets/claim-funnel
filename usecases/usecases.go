package usecases

import (
	"errors"
	"fmt"
	"github.com/lanemets/claim-funnel/model"
	"log"
)

type RpcClaim interface {
	Create(claim *model.Claim, profile *model.Profile) (*model.ClaimId, error)
	ConfirmClaim(claimId model.ClaimId)
}

type RpcProfile interface {
}

type ProcessDefinitionId struct {
	Value string
}

type BpmClaimClient interface {
	StartProcessInstance(claimId *model.ClaimId) (*ProcessDefinitionId, error)
}

type Interactor struct {
	claimClient   RpcClaim
	profileClient RpcProfile
	bpmClient     BpmClaimClient
}

func (s *Interactor) CreateClaim(claim *model.Claim, profile *model.Profile) (*ProcessDefinitionId, error) {
	claimId, err := s.claimClient.Create(claim, profile)
	if err != nil {
		errMsg := fmt.Sprintf("an error occurred on claim creation; claimId: %v, error: %v", claimId, err)
		log.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	processId, err := s.bpmClient.StartProcessInstance(claimId)
	if err != nil {
		errMsg := fmt.Sprintf("an error occurred on claim creation; claimId: %v, error: %v", claimId, err)
		log.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	return processId, nil
}

func (s *Interactor) ConfirmClaim(claimId model.ClaimId) {
	s.claimClient.ConfirmClaim(claimId)
}

func NewInteractor(
	claimClient RpcClaim,
	profileClient RpcProfile,
	bpmClient BpmClaimClient,
) *Interactor {
	return &Interactor{
		claimClient:   claimClient,
		profileClient: profileClient,
		bpmClient:     bpmClient,
	}
}
