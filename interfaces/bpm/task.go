package bpm

import (
	"fmt"
	c "github.com/lanemets/claim-funnel/external/claim/gen"
	p "github.com/lanemets/claim-funnel/external/profile/gen"
	"github.com/lanemets/claim-funnel/interfaces/benerpc"
	"log"
)

type Handler = func(variables map[string]string, businessKey string) (attributes map[string]interface{}, error error)

type ServiceTaskHandler struct {
	Handler  Handler
	Topic    string
	WorkerId string
}

func NotifyBeneficiary(grpcCtx *benerpc.GrpcContext) Handler {
	return func(variables map[string]string, businessKey string) (map[string]interface{}, error) {
		claimId := businessKey
		profileId, exists := variables["profileId"]

		log.Printf("Notifying Beneficiary; claimId: %v, profileId: %v", claimId, profileId)

		client := c.NewClaimServiceClient(grpcCtx.Connection())
		req := &c.NotifyBeneficiaryRequest{ClaimId: claimId, ExistingUser: exists}

		_, err := client.NotifyBeneficiary(grpcCtx.Context(), req)

		if err != nil {
			log.Fatalf("an error occured on notifying beneficiary: %v", err)
		}

		log.Printf("beneficiary has been successfully notified; claimId: %v", claimId)

		return make(map[string]interface{}), nil
	}
}

func SetPaymentPending(grpcCtx *benerpc.GrpcContext) Handler {
	return func(variables map[string]string, businessKey string) (map[string]interface{}, error) {

		claimId := businessKey
		request := &c.SetPaymentPendingRequest{ClaimId: claimId}

		log.Printf("Creating SetPaymentPending Request; claimId: %v", claimId)

		client := c.NewClaimServiceClient(grpcCtx.Connection())
		_, err := client.SetPaymentPending(grpcCtx.Context(), request)

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

func AcknowledgeClaim(grpcCtx *benerpc.GrpcContext) Handler {
	return func(variables map[string]string, businessKey string) (map[string]interface{}, error) {

		claimId := businessKey
		profileId := variables["profileId"]

		log.Printf("Creating AcknowledgeClaim Request; claimId: %v, profileId: %v", claimId, profileId)

		client := c.NewClaimServiceClient(grpcCtx.Connection())

		request := &c.AcknowledgeClaimRequest{
			ClaimId:   claimId,
			ProfileId: profileId,
		}

		_, err := client.AcknowledgeClaim(grpcCtx.Context(), request)

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

func GetClaimInfo(grpcCtx *benerpc.GrpcContext) Handler {
	return func(variables map[string]string, businessKey string) (map[string]interface{}, error) {
		claimId := businessKey
		client := c.NewClaimServiceClient(grpcCtx.Connection())
		req := &c.GetClaimRequest{Id: claimId}

		claim, err := client.GetClaim(grpcCtx.Context(), req)
		if err != nil {
			log.Fatalf("an error occured on get claim info: %v", err)
		}

		log.Printf("claim has been successfully retrieved; claimId: %v, claim: %v", claimId, claim)

		outputVariables := make(map[string]interface{})
		outputVariables["email"] = claim.Data.Email

		return outputVariables, nil
	}
}

func GetProfileByEmail(grpcCtx *benerpc.GrpcContext) Handler {
	return func(variables map[string]string, businessKey string) (map[string]interface{}, error) {

		email := variables["email"]
		log.Printf("Creating GetProfileByEmail Request; email: %v", email)

		client := p.NewProfilesServiceClient(grpcCtx.Connection())
		req := &p.GetProfileByEmailRequest{Email: email}

		profile, err := client.GetProfileByEmail(grpcCtx.Context(), req)
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
