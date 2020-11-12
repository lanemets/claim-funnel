package bpm

type Credentials struct {
	EndpointUrl string
	User        string
	Password    string
}

type Process struct {
	FilePath string
}

const (
	ClaimProcessKey       = "claim-process"
	ClaimConfirmTopicName = "claim-confirm"

	AcknowledgeClaimTopicName  = "acknowledge-claim"
	SetPaymentPendingTopicName = "set-payment-pending"
	GetClaimTopicName          = "get-current-claim"
	NotifyBeneficiaryTopicName = "notify-beneficiary"
	GetProfileByEmailTopicName = "get-profile-by-email"

	ClaimConfirmTaskId = "claim-confirm"
)
