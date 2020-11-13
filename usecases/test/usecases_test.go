package usecases

import (
	"github.com/lanemets/claim-funnel/model"
	"github.com/lanemets/claim-funnel/usecases"
	"github.com/lanemets/claim-funnel/usecases/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UseCaseSuite struct {
	suite.Suite
}

func TestUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(UseCaseSuite))
}

func (s *UseCaseSuite) TestCreateClaim() {
	// Given
	bpmClient := new(mocks.BpmClaimClient)

	claimIdExpected := &model.ClaimId{Value: "claimId"}
	processDefinitionIdExpected := &model.ProcessDefinitionId{Value: "processDefinition"}
	bpmClient.On("StartProcessInstance", claimIdExpected).Return(processDefinitionIdExpected, nil)

	claim := &model.Claim{
		Email:                 "",
		Amount:                "",
		CurrencyCode:          "",
		ClientReferenceNumber: "",
		Description:           "",
		CreatedOn:             "",
		ExpiresOn:             "",
	}
	profile := &model.Profile{
		ExternalId:    "",
		ProfileType:   0,
		ProfileEntity: nil,
		Address:       nil,
		Phone:         nil,
	}
	rpcClaim := new(mocks.RpcClaim)
	rpcClaim.On("Create", claim, profile).Return(claimIdExpected, nil)

	rpcProfile := new(mocks.RpcProfile)
	interactor := usecases.Interactor{
		ClaimClient:   rpcClaim,
		ProfileClient: rpcProfile,
		BpmClient:     bpmClient,
	}

	// When
	claimId, processDefinitionId, err := interactor.CreateClaim(claim, profile)

	// Then
	s.Equal(claimId, claimIdExpected)
	s.Equal(processDefinitionId, processDefinitionIdExpected)
	s.Nil(err)
}

func (s *UseCaseSuite) TestConfirmClaim() {
	// Given
	bpmClient := new(mocks.BpmClaimClient)

	rpcClaim := new(mocks.RpcClaim)
	claimId := &model.ClaimId{Value: "claimId"}
	rpcClaim.On("ConfirmClaim", claimId).Return(nil)

	rpcProfile := new(mocks.RpcProfile)

	interactor := usecases.Interactor{
		ClaimClient:   rpcClaim,
		ProfileClient: rpcProfile,
		BpmClient:     bpmClient,
	}

	// When
	err := interactor.ConfirmClaim(claimId)

	// Then
	s.Nil(err)
}
