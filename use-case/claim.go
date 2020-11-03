package use_case

import (
	"github.com/lanemets/claim-funnel/bpm"
	"github.com/lanemets/claim-funnel/external/claim"
	"github.com/lanemets/claim-funnel/external/profile"
	"log"
)

type Claim struct {
	bpmClient     *bpm.Client
	claimClient   *claim.Client
	profileClient *profile.Client

	claimId      *claim.ClaimId
	processDefinitionId *bpm.ProcessDefinitionId
}

func (uc Claim) Confirm(claimId claim.ClaimId, processDefinitionId bpm.ProcessDefinitionId) {
	uc.profileClient.ConfirmClaim(&claimId)
	log.Printf("Claim: %s has been confirmes sucessfully", claimId)

	uc.bpmClient.CompleteClaimConfirmTask(&claimId, &processDefinitionId)
	log.Printf("BPM Task has been completed; processDefinitionId: %s", uc.processDefinitionId)
}
