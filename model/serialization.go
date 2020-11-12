package model

import (
	proto "github.com/lanemets/claim-funnel/external/claim/gen"
)

func fromDto(claim *proto.Claim) *Claim {
	return &Claim{
		Email:                 claim.Email,
		Amount:                claim.Amount,
		CurrencyCode:          claim.CurrencyCode,
		ClientReferenceNumber: claim.ClientReferenceNumber,
		Description:           claim.Description,
		CreatedOn:             claim.CreatedOn,
		ExpiresOn:             claim.ExpiresOn,
	}
}
