package benerpc

import (
	"github.com/lanemets/claim-funnel/usecases"
)

type ProfileClient struct {
	ctx *GrpcContext
}

func NewRpcProfile(grpcContext *GrpcContext) usecases.RpcProfile {
	return ProfileClient{
		ctx: grpcContext,
	}
}
