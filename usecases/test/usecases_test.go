package usecases_test

import (
	"errors"
	"fmt"
	"github.com/lanemets/claim-funnel/model"
	"github.com/lanemets/claim-funnel/usecases"
	"github.com/lanemets/claim-funnel/usecases/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
)

const (
	ConfirmClaim         = "ConfirmClaim"
	StartProcessInstance = "StartProcessInstance"
	Create               = "Create"
)

type UseCaseSuite struct {
	suite.Suite

	bpmClient  mocks.BpmClaimClient
	rpcClaim   mocks.RpcClaim
	rpcProfile mocks.RpcProfile
}

func (suite *UseCaseSuite) SetupTest() {
	suite.rpcProfile = mocks.RpcProfile{}
	suite.rpcClaim = mocks.RpcClaim{}
	suite.bpmClient = mocks.BpmClaimClient{}
}

func TestUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(UseCaseSuite))
}

func (suite *UseCaseSuite) TestCreateClaim() {
	// Given
	dataProvider := []struct {
		claimId             *model.ClaimId
		processDefinitionId *model.ProcessDefinitionId

		claimError error
		bpmError   error

		processDefinitionIdExpected *model.ProcessDefinitionId
		claimIdExpected             *model.ClaimId

		errorExpected error
	}{
		{
			&model.ClaimId{Value: "claimId0"},
			&model.ProcessDefinitionId{Value: "processDefinitionId0"},
			nil,
			nil,

			&model.ProcessDefinitionId{Value: "processDefinitionId0"},
			&model.ClaimId{Value: "claimId0"},
			nil,
		},
		{
			nil,
			&model.ProcessDefinitionId{Value: "processDefinitionId1"},
			errors.New("create claim transport error"),
			nil,

			nil,
			nil,
			errors.New("an error occurred on claim creation; error: create claim transport error"),
		},
		{
			&model.ClaimId{Value: "claimId2"},
			nil,
			nil,
			errors.New("start process error"),

			nil,
			&model.ClaimId{Value: "claimId2"},
			errors.New("an error occurred on starting process; claimId: claimId2, error: start process error"),
		},
		{
			nil,
			nil,
			errors.New("create claim transport error"),
			errors.New("start process error"),

			nil,
			nil,
			errors.New("an error occurred on claim creation; error: create claim transport error"),
		},
	}

	for i, data := range dataProvider {

		claim := &model.Claim{
			Email:        fmt.Sprint(i),
			Amount:       "",
			CurrencyCode: "",
			Description:  "",
		}
		profile := &model.Profile{
			ExternalID:  fmt.Sprint(i),
			ProfileType: "PERSON",
			Address:     model.Address{},
			Phone:       model.Phone{},
		}

		suite.rpcClaim.On(Create, claim, profile).
			Return(data.claimId, data.claimError)

		suite.bpmClient.On(StartProcessInstance, data.claimId).
			Return(data.processDefinitionId, data.bpmError)

		interactor := &usecases.Interactor{
			RpcClaim:   &suite.rpcClaim,
			BpmClient:  &suite.bpmClient,
		}

		claimId, processDefinitionId, err := interactor.CreateClaim(claim, profile)

		// Then
		suite.Equal(data.claimIdExpected, claimId)
		suite.Equal(data.processDefinitionIdExpected, processDefinitionId)
		suite.Equal(data.errorExpected, err)
	}
}

func (suite *UseCaseSuite) TestConfirmClaim() {
	dataProvider := []struct {
		confirmResult    interface{}
		interactorResult error
	}{
		{nil, nil},
		{
			errors.New("an error has occurred on confirming claim"),
			errors.New("an error has occurred on confirming claim"),
		},
	}

	for i, data := range dataProvider {

		// Given
		claimId := &model.ClaimId{Value: fmt.Sprint("claimId", i)}
		suite.rpcClaim.On(ConfirmClaim, claimId).Return(data.confirmResult)

		interactor := &usecases.Interactor{
			RpcClaim:   &suite.rpcClaim,
			BpmClient:  &suite.bpmClient,
		}

		// When
		err := interactor.ConfirmClaim(claimId)

		// Then
		suite.Equal(err, data.interactorResult)
	}
}
