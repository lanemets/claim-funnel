// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package profile

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// ClaimServiceClient is the client API for ClaimService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ClaimServiceClient interface {
	NotifyBeneficiary(ctx context.Context, in *NotifyBeneficiaryRequest, opts ...grpc.CallOption) (*NotifyBeneficiaryResponse, error)
	ConfirmClaim(ctx context.Context, in *ConfirmClaimRequest, opts ...grpc.CallOption) (*ConfirmClaimResponse, error)
	SetPaymentPending(ctx context.Context, in *SetPaymentPendingRequest, opts ...grpc.CallOption) (*SetPaymentPendingResponse, error)
	AcknowledgeClaim(ctx context.Context, in *AcknowledgeClaimRequest, opts ...grpc.CallOption) (*AcknowledgeClaimResponse, error)
}

type claimServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewClaimServiceClient(cc grpc.ClientConnInterface) ClaimServiceClient {
	return &claimServiceClient{cc}
}

func (c *claimServiceClient) NotifyBeneficiary(ctx context.Context, in *NotifyBeneficiaryRequest, opts ...grpc.CallOption) (*NotifyBeneficiaryResponse, error) {
	out := new(NotifyBeneficiaryResponse)
	err := c.cc.Invoke(ctx, "/profile.ClaimService/NotifyBeneficiary", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *claimServiceClient) ConfirmClaim(ctx context.Context, in *ConfirmClaimRequest, opts ...grpc.CallOption) (*ConfirmClaimResponse, error) {
	out := new(ConfirmClaimResponse)
	err := c.cc.Invoke(ctx, "/profile.ClaimService/ConfirmClaim", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *claimServiceClient) SetPaymentPending(ctx context.Context, in *SetPaymentPendingRequest, opts ...grpc.CallOption) (*SetPaymentPendingResponse, error) {
	out := new(SetPaymentPendingResponse)
	err := c.cc.Invoke(ctx, "/profile.ClaimService/SetPaymentPending", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *claimServiceClient) AcknowledgeClaim(ctx context.Context, in *AcknowledgeClaimRequest, opts ...grpc.CallOption) (*AcknowledgeClaimResponse, error) {
	out := new(AcknowledgeClaimResponse)
	err := c.cc.Invoke(ctx, "/profile.ClaimService/AcknowledgeClaim", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClaimServiceServer is the server API for ClaimService service.
// All implementations must embed UnimplementedClaimServiceServer
// for forward compatibility
type ClaimServiceServer interface {
	NotifyBeneficiary(context.Context, *NotifyBeneficiaryRequest) (*NotifyBeneficiaryResponse, error)
	ConfirmClaim(context.Context, *ConfirmClaimRequest) (*ConfirmClaimResponse, error)
	SetPaymentPending(context.Context, *SetPaymentPendingRequest) (*SetPaymentPendingResponse, error)
	AcknowledgeClaim(context.Context, *AcknowledgeClaimRequest) (*AcknowledgeClaimResponse, error)
	mustEmbedUnimplementedClaimServiceServer()
}

// UnimplementedClaimServiceServer must be embedded to have forward compatible implementations.
type UnimplementedClaimServiceServer struct {
}

func (UnimplementedClaimServiceServer) NotifyBeneficiary(context.Context, *NotifyBeneficiaryRequest) (*NotifyBeneficiaryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NotifyBeneficiary not implemented")
}
func (UnimplementedClaimServiceServer) ConfirmClaim(context.Context, *ConfirmClaimRequest) (*ConfirmClaimResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConfirmClaim not implemented")
}
func (UnimplementedClaimServiceServer) SetPaymentPending(context.Context, *SetPaymentPendingRequest) (*SetPaymentPendingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetPaymentPending not implemented")
}
func (UnimplementedClaimServiceServer) AcknowledgeClaim(context.Context, *AcknowledgeClaimRequest) (*AcknowledgeClaimResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AcknowledgeClaim not implemented")
}
func (UnimplementedClaimServiceServer) mustEmbedUnimplementedClaimServiceServer() {}

// UnsafeClaimServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ClaimServiceServer will
// result in compilation errors.
type UnsafeClaimServiceServer interface {
	mustEmbedUnimplementedClaimServiceServer()
}

func RegisterClaimServiceServer(s grpc.ServiceRegistrar, srv ClaimServiceServer) {
	s.RegisterService(&_ClaimService_serviceDesc, srv)
}

func _ClaimService_NotifyBeneficiary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NotifyBeneficiaryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClaimServiceServer).NotifyBeneficiary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/profile.ClaimService/NotifyBeneficiary",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClaimServiceServer).NotifyBeneficiary(ctx, req.(*NotifyBeneficiaryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClaimService_ConfirmClaim_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfirmClaimRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClaimServiceServer).ConfirmClaim(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/profile.ClaimService/ConfirmClaim",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClaimServiceServer).ConfirmClaim(ctx, req.(*ConfirmClaimRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClaimService_SetPaymentPending_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetPaymentPendingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClaimServiceServer).SetPaymentPending(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/profile.ClaimService/SetPaymentPending",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClaimServiceServer).SetPaymentPending(ctx, req.(*SetPaymentPendingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClaimService_AcknowledgeClaim_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AcknowledgeClaimRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClaimServiceServer).AcknowledgeClaim(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/profile.ClaimService/AcknowledgeClaim",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClaimServiceServer).AcknowledgeClaim(ctx, req.(*AcknowledgeClaimRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ClaimService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "profile.ClaimService",
	HandlerType: (*ClaimServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NotifyBeneficiary",
			Handler:    _ClaimService_NotifyBeneficiary_Handler,
		},
		{
			MethodName: "ConfirmClaim",
			Handler:    _ClaimService_ConfirmClaim_Handler,
		},
		{
			MethodName: "SetPaymentPending",
			Handler:    _ClaimService_SetPaymentPending_Handler,
		},
		{
			MethodName: "AcknowledgeClaim",
			Handler:    _ClaimService_AcknowledgeClaim_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "profile.proto",
}

// ProfilesServiceClient is the client API for ProfilesService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProfilesServiceClient interface {
	GetProfileByEmail(ctx context.Context, in *GetProfileByEmailRequest, opts ...grpc.CallOption) (*GetProfileByEmailResponse, error)
}

type profilesServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewProfilesServiceClient(cc grpc.ClientConnInterface) ProfilesServiceClient {
	return &profilesServiceClient{cc}
}

func (c *profilesServiceClient) GetProfileByEmail(ctx context.Context, in *GetProfileByEmailRequest, opts ...grpc.CallOption) (*GetProfileByEmailResponse, error) {
	out := new(GetProfileByEmailResponse)
	err := c.cc.Invoke(ctx, "/profile.ProfilesService/GetProfileByEmail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProfilesServiceServer is the server API for ProfilesService service.
// All implementations must embed UnimplementedProfilesServiceServer
// for forward compatibility
type ProfilesServiceServer interface {
	GetProfileByEmail(context.Context, *GetProfileByEmailRequest) (*GetProfileByEmailResponse, error)
	mustEmbedUnimplementedProfilesServiceServer()
}

// UnimplementedProfilesServiceServer must be embedded to have forward compatible implementations.
type UnimplementedProfilesServiceServer struct {
}

func (UnimplementedProfilesServiceServer) GetProfileByEmail(context.Context, *GetProfileByEmailRequest) (*GetProfileByEmailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProfileByEmail not implemented")
}
func (UnimplementedProfilesServiceServer) mustEmbedUnimplementedProfilesServiceServer() {}

// UnsafeProfilesServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProfilesServiceServer will
// result in compilation errors.
type UnsafeProfilesServiceServer interface {
	mustEmbedUnimplementedProfilesServiceServer()
}

func RegisterProfilesServiceServer(s grpc.ServiceRegistrar, srv ProfilesServiceServer) {
	s.RegisterService(&_ProfilesService_serviceDesc, srv)
}

func _ProfilesService_GetProfileByEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProfileByEmailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfilesServiceServer).GetProfileByEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/profile.ProfilesService/GetProfileByEmail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfilesServiceServer).GetProfileByEmail(ctx, req.(*GetProfileByEmailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ProfilesService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "profile.ProfilesService",
	HandlerType: (*ProfilesServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetProfileByEmail",
			Handler:    _ProfilesService_GetProfileByEmail_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "profile.proto",
}
