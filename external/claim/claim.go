package claim

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"time"

	c "github.com/lanemets/claim-funnel/external/claim/gen"
)

const (
	address = "localhost:9002"
)

type ClaimId struct {
	Value string
}

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (claim Client) Create(request *c.CreateClaimRequest) *ClaimId {
	//TODO pool connections ?
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	client := c.NewClaimServiceClient(conn)
	claimId, err := client.CreateClaim(ctx, request)
	if err != nil {
		log.Fatalf("an error occured on creating claim: %v", err)
	}

	log.Printf("claim has been created successfully; claimId: %v", claimId)

	return &ClaimId{Value: claimId.Id}
}

func (claim Client) ConfirmClaim(claimId *ClaimId) {
	//TODO: get rid of this copy-paste shitshow
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	client := c.NewClaimServiceClient(conn)

	req := &c.ConfirmClaimRequest{ClaimId: claimId.Value}
	_, err = client.ConfirmClaim(ctx, req)

	if err != nil {
		log.Fatalf("an error occured on confirming claim: %v", err)
	}

	log.Printf("claim has been confirmed successfully; claimId: %v", claimId)
}
