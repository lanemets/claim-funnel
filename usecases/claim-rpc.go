package usecases

import "github.com/lanemets/claim-funnel/model"

type RpcClaim interface {
	Create(claim *model.Claim, profile *model.Profile) (*model.ClaimId, error)
	ConfirmClaim(claimId *model.ClaimId) error
}
