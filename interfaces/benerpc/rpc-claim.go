package benerpc

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	grpcClient "github.com/lanemets/claim-funnel/external/claim/gen"
	"github.com/lanemets/claim-funnel/model"
	"github.com/lanemets/claim-funnel/usecases"
	"log"
)

type ClaimClient struct {
	ctx *GrpcContext
}

func (client ClaimClient) Create(claim *model.Claim, profile *model.Profile) (*model.ClaimId, error) {
	request := &grpcClient.CreateClaimRequest{
		Claim: &grpcClient.Claim{
			Email:                 "lanemets.vv+xn3j@gmail.com",
			Amount:                "100.0",
			CurrencyCode:          "USD",
			Description:           "Test Claim To Create",
			ClientReferenceNumber: uuid.New().String(),
		},
		Profile: &grpcClient.Profile{
			ExternalId:  "0",
			ProfileType: 0,
			Entity: &grpcClient.Profile_Person{
				Person: &grpcClient.Person{
					FirstName: "TestFirsdt0",
					LastName:  "TestLastd0",
				},
			},
		},
	}

	claimServiceClient := grpcClient.NewClaimServiceClient(client.ctx.Connection())

	response, err := claimServiceClient.CreateClaim(client.ctx.Context(), request)
	if err != nil {
		errMsg := fmt.Sprintf("an error occured on creating claim: %v , error %v", claim, err)
		log.Printf(errMsg)
		return nil, errors.New(errMsg)
	}

	log.Printf("claim has been created successfully; claimId: %v", response.Id)

	return &model.ClaimId{Value: response.Id}, nil
}

func (client ClaimClient) ConfirmClaim(claimId *model.ClaimId) error {
	claimServiceClient := grpcClient.NewClaimServiceClient(client.ctx.Connection())

	req := &grpcClient.ConfirmClaimRequest{ClaimId: claimId.Value}
	_, err := claimServiceClient.ConfirmClaim(client.ctx.Context(), req)

	if err != nil {
		errMsg := fmt.Sprintf("an error has occurred on confirming claim: %v", err)
		log.Println(errMsg)
		return errors.New(errMsg)
	}

	log.Printf("claim has been confirmed successfully; claimId: %v", claimId)
	return nil
}

func NewRpcClaim(grpcContext *GrpcContext) usecases.RpcClaim {
	return ClaimClient{ctx: grpcContext}
}
