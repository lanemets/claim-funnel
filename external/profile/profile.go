package profile

import (
	"context"
	"github.com/lanemets/claim-funnel/external/claim"
	pb "github.com/lanemets/claim-funnel/external/profile/gen"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	address = "localhost:9001"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

type ProfileId struct {
	Value string
}

//TODO: move to claim-service
func (profile Client) ConfirmClaim(claimId *claim.ClaimId) {
	//TODO: get rid of this copy-paste shitshow
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	client := pb.NewClaimServiceClient(conn)

	req := &pb.ConfirmClaimRequest{ClaimId: claimId.Value}
	_, err = client.ConfirmClaim(ctx, req)

	if err != nil {
		log.Fatalf("an error occured on confirming claim: %v", err)
	}

	log.Printf("beneficiary has been successfully notified; claimId: %v", claimId)
}