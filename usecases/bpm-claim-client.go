package usecases

import "github.com/lanemets/claim-funnel/model"

type BpmClaimClient interface {
	StartProcessInstance(claimId *model.ClaimId) (*model.ProcessDefinitionId, error)
}
