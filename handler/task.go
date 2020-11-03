package handler

import (
	"context"
	"fmt"
	c "github.com/lanemets/claim-funnel/external/claim/gen"
	p "github.com/lanemets/claim-funnel/external/profile/gen"
	"google.golang.org/grpc"
	"log"
	"time"
)

type Handler = func(variables map[string]string, businessKey string) (attributes map[string]interface{}, error error)

type ServiceHandler struct {
	Handler  Handler
	Topic    string
	WorkerId string
}

//TODO: get rid of address
func NotifyBeneficiary(address string) Handler {
	return func(variables map[string]string, businessKey string) (map[string]interface{}, error) {
		conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		claimId := businessKey
		profileId, exists := variables["profileId"]

		log.Printf("Notifying Beneficiary; claimId: %v, profileId: %v", claimId, profileId)

		client := p.NewClaimServiceClient(conn)
		req := &p.NotifyBeneficiaryRequest{ClaimId: claimId, ExistingUser: exists}

		_, err = client.NotifyBeneficiary(ctx, req)

		if err != nil {
			log.Fatalf("an error occured on notifying beneficiary: %v", err)
		}

		log.Printf("beneficiary has been successfully notified; claimId: %v", claimId)

		return make(map[string]interface{}), nil
	}
}

func SetPaymentPending(address string) Handler {
	return func(variables map[string]string, businessKey string) (map[string]interface{}, error) {
		conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		claimId := businessKey
		request := &p.SetPaymentPendingRequest{ClaimId: claimId}

		log.Printf("Creating SetPaymentPending Request; claimId: %v", claimId)

		client := p.NewClaimServiceClient(conn)
		_, err = client.SetPaymentPending(ctx, request)

		if err != nil {
			log.Fatalf("an error occured on SetPaymentPending request %v", err)
		}

		log.Printf(
			"SetPaymentPending operation has been successfully done; claimId: %v",
			claimId,
		)

		return make(map[string]interface{}), nil
	}
}

func AcknowledgeClaim(address string) Handler {
	return func(variables map[string]string, businessKey string) (map[string]interface{}, error) {
		conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		claimId := businessKey
		profileId := variables["profileId"]

		log.Printf("Creating AcknowledgeClaim Request; claimId: %v, profileId: %v", claimId, profileId)

		client := p.NewClaimServiceClient(conn)

		request := &p.AcknowledgeClaimRequest{
			ClaimId:   claimId,
			ProfileId: profileId,
		}

		_, err = client.AcknowledgeClaim(ctx, request)

		if err != nil {
			log.Fatalf("an error occured on AcknowledgeClaim request: %v", err)
		}

		log.Printf(
			"AcknowledgeClaim operation has been successfully done; claimId: %v, profileId: %v",
			claimId,
			profileId,
		)

		return make(map[string]interface{}), nil
	}
}

func GetClaimInfo(address string) Handler {
	return func(variables map[string]string, businessKey string) (map[string]interface{}, error) {
		conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		claimId := businessKey
		client := c.NewClaimServiceClient(conn)
		req := &c.GetClaimRequest{Id: claimId}

		claim, err := client.GetClaim(ctx, req)
		if err != nil {
			log.Fatalf("an error occured on get claim info: %v", err)
		}

		log.Printf("claim has been successfully retrieved; claimId: %v, claim: %v", claimId, claim)

		outputVariables := make(map[string]interface{})
		outputVariables["email"] = claim.Data.Email

		return outputVariables, nil
	}
}

func GetProfileByEmail(address string) Handler {
	return func(variables map[string]string, businessKey string) (map[string]interface{}, error) {
		conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		email := variables["email"]
		log.Printf("Creating GetProfileByEmail Request; email: %v", email)

		client := p.NewProfilesServiceClient(conn)
		req := &p.GetProfileByEmailRequest{Email: email}

		profile, err := client.GetProfileByEmail(ctx, req)
		if err != nil {
			log.Fatalf("an error occured on GetProfileByEmail request: %v", err)
		}

		log.Printf("Get Profile By Email operation has been successfully done; profile: %v", profile)

		outputVariables := make(map[string]interface{})

		switch response := profile.Response.(type) {
		case *p.GetProfileByEmailResponse_Exists:
			{
				log.Printf("Got ProfileId: %v", response.Exists.ProfileId)
				outputVariables["profileId"] = response.Exists.ProfileId
				outputVariables["profileExists"] = true
			}
		case *p.GetProfileByEmailResponse_NotFound:
			{
				log.Printf("No Profile Id Found")
				outputVariables["profileExists"] = false
			}
		default:
			_ = fmt.Errorf("GetProfileByEmailResponse has an unexpected type %v", response)
		}

		return outputVariables, nil
	}
}
