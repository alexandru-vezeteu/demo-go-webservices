
package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const _ = grpc.SupportPackageIsVersion9

const (
	IdentityService_Login_FullMethodName              = "/idm.IdentityService/Login"
	IdentityService_VerifyToken_FullMethodName        = "/idm.IdentityService/VerifyToken"
	IdentityService_RevokeToken_FullMethodName        = "/idm.IdentityService/RevokeToken"
	IdentityService_CheckPermission_FullMethodName    = "/idm.IdentityService/CheckPermission"
	IdentityService_WriteRelationships_FullMethodName = "/idm.IdentityService/WriteRelationships"
	IdentityService_ReadRelationships_FullMethodName  = "/idm.IdentityService/ReadRelationships"
)

type IdentityServiceClient interface {
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	VerifyToken(ctx context.Context, in *VerifyTokenRequest, opts ...grpc.CallOption) (*VerifyTokenResponse, error)
	RevokeToken(ctx context.Context, in *RevokeTokenRequest, opts ...grpc.CallOption) (*RevokeTokenResponse, error)
	CheckPermission(ctx context.Context, in *CheckPermissionRequest, opts ...grpc.CallOption) (*CheckPermissionResponse, error)
	WriteRelationships(ctx context.Context, in *WriteRelationshipsRequest, opts ...grpc.CallOption) (*WriteRelationshipsResponse, error)
	ReadRelationships(ctx context.Context, in *ReadRelationshipsRequest, opts ...grpc.CallOption) (*ReadRelationshipsResponse, error)
}

type identityServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewIdentityServiceClient(cc grpc.ClientConnInterface) IdentityServiceClient {
	return &identityServiceClient{cc}
}

func (c *identityServiceClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, IdentityService_Login_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *identityServiceClient) VerifyToken(ctx context.Context, in *VerifyTokenRequest, opts ...grpc.CallOption) (*VerifyTokenResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(VerifyTokenResponse)
	err := c.cc.Invoke(ctx, IdentityService_VerifyToken_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *identityServiceClient) RevokeToken(ctx context.Context, in *RevokeTokenRequest, opts ...grpc.CallOption) (*RevokeTokenResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RevokeTokenResponse)
	err := c.cc.Invoke(ctx, IdentityService_RevokeToken_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *identityServiceClient) CheckPermission(ctx context.Context, in *CheckPermissionRequest, opts ...grpc.CallOption) (*CheckPermissionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CheckPermissionResponse)
	err := c.cc.Invoke(ctx, IdentityService_CheckPermission_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *identityServiceClient) WriteRelationships(ctx context.Context, in *WriteRelationshipsRequest, opts ...grpc.CallOption) (*WriteRelationshipsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(WriteRelationshipsResponse)
	err := c.cc.Invoke(ctx, IdentityService_WriteRelationships_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *identityServiceClient) ReadRelationships(ctx context.Context, in *ReadRelationshipsRequest, opts ...grpc.CallOption) (*ReadRelationshipsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReadRelationshipsResponse)
	err := c.cc.Invoke(ctx, IdentityService_ReadRelationships_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type IdentityServiceServer interface {
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	VerifyToken(context.Context, *VerifyTokenRequest) (*VerifyTokenResponse, error)
	RevokeToken(context.Context, *RevokeTokenRequest) (*RevokeTokenResponse, error)
	CheckPermission(context.Context, *CheckPermissionRequest) (*CheckPermissionResponse, error)
	WriteRelationships(context.Context, *WriteRelationshipsRequest) (*WriteRelationshipsResponse, error)
	ReadRelationships(context.Context, *ReadRelationshipsRequest) (*ReadRelationshipsResponse, error)
	mustEmbedUnimplementedIdentityServiceServer()
}

type UnimplementedIdentityServiceServer struct{}

func (UnimplementedIdentityServiceServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedIdentityServiceServer) VerifyToken(context.Context, *VerifyTokenRequest) (*VerifyTokenResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method VerifyToken not implemented")
}
func (UnimplementedIdentityServiceServer) RevokeToken(context.Context, *RevokeTokenRequest) (*RevokeTokenResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method RevokeToken not implemented")
}
func (UnimplementedIdentityServiceServer) CheckPermission(context.Context, *CheckPermissionRequest) (*CheckPermissionResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method CheckPermission not implemented")
}
func (UnimplementedIdentityServiceServer) WriteRelationships(context.Context, *WriteRelationshipsRequest) (*WriteRelationshipsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method WriteRelationships not implemented")
}
func (UnimplementedIdentityServiceServer) ReadRelationships(context.Context, *ReadRelationshipsRequest) (*ReadRelationshipsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method ReadRelationships not implemented")
}
func (UnimplementedIdentityServiceServer) mustEmbedUnimplementedIdentityServiceServer() {}
func (UnimplementedIdentityServiceServer) testEmbeddedByValue()                         {}

type UnsafeIdentityServiceServer interface {
	mustEmbedUnimplementedIdentityServiceServer()
}

func RegisterIdentityServiceServer(s grpc.ServiceRegistrar, srv IdentityServiceServer) {
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&IdentityService_ServiceDesc, srv)
}

func _IdentityService_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IdentityServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IdentityService_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IdentityServiceServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IdentityService_VerifyToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IdentityServiceServer).VerifyToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IdentityService_VerifyToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IdentityServiceServer).VerifyToken(ctx, req.(*VerifyTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IdentityService_RevokeToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RevokeTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IdentityServiceServer).RevokeToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IdentityService_RevokeToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IdentityServiceServer).RevokeToken(ctx, req.(*RevokeTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IdentityService_CheckPermission_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckPermissionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IdentityServiceServer).CheckPermission(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IdentityService_CheckPermission_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IdentityServiceServer).CheckPermission(ctx, req.(*CheckPermissionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IdentityService_WriteRelationships_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WriteRelationshipsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IdentityServiceServer).WriteRelationships(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IdentityService_WriteRelationships_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IdentityServiceServer).WriteRelationships(ctx, req.(*WriteRelationshipsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IdentityService_ReadRelationships_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadRelationshipsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IdentityServiceServer).ReadRelationships(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IdentityService_ReadRelationships_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IdentityServiceServer).ReadRelationships(ctx, req.(*ReadRelationshipsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var IdentityService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "idm.IdentityService",
	HandlerType: (*IdentityServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _IdentityService_Login_Handler,
		},
		{
			MethodName: "VerifyToken",
			Handler:    _IdentityService_VerifyToken_Handler,
		},
		{
			MethodName: "RevokeToken",
			Handler:    _IdentityService_RevokeToken_Handler,
		},
		{
			MethodName: "CheckPermission",
			Handler:    _IdentityService_CheckPermission_Handler,
		},
		{
			MethodName: "WriteRelationships",
			Handler:    _IdentityService_WriteRelationships_Handler,
		},
		{
			MethodName: "ReadRelationships",
			Handler:    _IdentityService_ReadRelationships_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/idm.proto",
}
