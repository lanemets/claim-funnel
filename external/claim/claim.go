package claim

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"log"
	"time"

	cb "github.com/lanemets/claim-funnel/external/claim/gen"
)

const (
	address = "localhost:9002"
)

type Id struct {
	Value string
}

type Claim struct {
}

func Client() *Claim {
	return &Claim{}
}

func (claim Claim) Create() Id {
	//TODO pool connections ?
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	//TODO move to params
	request := &cb.CreateClaimRequest{
		Claim: &cb.Claim{
			Email:                 "lanemets.vv+0m7@gmail.com",
			Amount:                "100.0",
			CurrencyCode:          "USD",
			ClientReferenceNumber: uuid.New().String(),
			Description: "Test",
		},
		Profile: &cb.Profile{
			ExternalId:  "0",
			ProfileType: 0,
			Entity: &cb.Profile_Person{
				Person: &cb.Person{
					FirstName: "TestFirst0",
					LastName:  "TestLast0",
				},
			},
		},
	}

	client := cb.NewClaimServiceClient(conn)
	claimId, err := client.CreateClaim(ctx, request)
	if err != nil {
		log.Fatalf("an error occured on creating claim: %v", err)
	}

	log.Printf("claim has been created successfully; claimId: %v", claimId)

	return Id{Value: claimId.Id}
}
