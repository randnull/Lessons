// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: tutors.proto

package __

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	UserService_CreateUser_FullMethodName             = "/users.UserService/CreateUser"
	UserService_GetUserById_FullMethodName            = "/users.UserService/GetUserById"
	UserService_GetTutorById_FullMethodName           = "/users.UserService/GetTutorById"
	UserService_GetTutorInfoById_FullMethodName       = "/users.UserService/GetTutorInfoById"
	UserService_GetStudentById_FullMethodName         = "/users.UserService/GetStudentById"
	UserService_GetUserByTelegramId_FullMethodName    = "/users.UserService/GetUserByTelegramId"
	UserService_GetAllUsers_FullMethodName            = "/users.UserService/GetAllUsers"
	UserService_GetAllTutorsPagination_FullMethodName = "/users.UserService/GetAllTutorsPagination"
	UserService_UpdateBioTutor_FullMethodName         = "/users.UserService/UpdateBioTutor"
	UserService_UpdateTags_FullMethodName             = "/users.UserService/UpdateTags"
	UserService_ChangeTutorActive_FullMethodName      = "/users.UserService/ChangeTutorActive"
	UserService_ChangeTutorName_FullMethodName        = "/users.UserService/ChangeTutorName"
	UserService_CreateNewResponse_FullMethodName      = "/users.UserService/CreateNewResponse"
	UserService_AddResponsesToTutor_FullMethodName    = "/users.UserService/AddResponsesToTutor"
	UserService_CreateReview_FullMethodName           = "/users.UserService/CreateReview"
	UserService_GetReview_FullMethodName              = "/users.UserService/GetReview"
	UserService_GetReviews_FullMethodName             = "/users.UserService/GetReviews"
	UserService_SetReviewActive_FullMethodName        = "/users.UserService/SetReviewActive"
	UserService_BanUser_FullMethodName                = "/users.UserService/BanUser"
)

// UserServiceClient is the client API for UserService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserServiceClient interface {
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateResponse, error)
	GetUserById(ctx context.Context, in *GetById, opts ...grpc.CallOption) (*User, error)
	GetTutorById(ctx context.Context, in *GetById, opts ...grpc.CallOption) (*Tutor, error)
	GetTutorInfoById(ctx context.Context, in *GetById, opts ...grpc.CallOption) (*TutorDetails, error)
	GetStudentById(ctx context.Context, in *GetById, opts ...grpc.CallOption) (*User, error)
	GetUserByTelegramId(ctx context.Context, in *GetByTelegramId, opts ...grpc.CallOption) (*User, error)
	GetAllUsers(ctx context.Context, in *GetAllRequest, opts ...grpc.CallOption) (*GetAllResponse, error)
	GetAllTutorsPagination(ctx context.Context, in *GetAllTutorsPaginationRequest, opts ...grpc.CallOption) (*GetTutorsPaginationResponse, error)
	UpdateBioTutor(ctx context.Context, in *UpdateBioRequest, opts ...grpc.CallOption) (*Success, error)
	UpdateTags(ctx context.Context, in *UpdateTagsRequest, opts ...grpc.CallOption) (*Success, error)
	ChangeTutorActive(ctx context.Context, in *SetActiveTutorById, opts ...grpc.CallOption) (*Success, error)
	ChangeTutorName(ctx context.Context, in *ChangeNameRequest, opts ...grpc.CallOption) (*Success, error)
	CreateNewResponse(ctx context.Context, in *CreateResponseRequest, opts ...grpc.CallOption) (*Success, error)
	AddResponsesToTutor(ctx context.Context, in *AddResponseToTutorRequest, opts ...grpc.CallOption) (*AddResponseToTutorResponse, error)
	CreateReview(ctx context.Context, in *CreateReviewRequest, opts ...grpc.CallOption) (*CreateReviewResponse, error)
	GetReview(ctx context.Context, in *GetReviewRequest, opts ...grpc.CallOption) (*Review, error)
	GetReviews(ctx context.Context, in *GetReviewsRequest, opts ...grpc.CallOption) (*GetReviewsResponse, error)
	SetReviewActive(ctx context.Context, in *SetReviewsActiveRequest, opts ...grpc.CallOption) (*Success, error)
	BanUser(ctx context.Context, in *BanUserRequest, opts ...grpc.CallOption) (*Success, error)
}

type userServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateResponse)
	err := c.cc.Invoke(ctx, UserService_CreateUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetUserById(ctx context.Context, in *GetById, opts ...grpc.CallOption) (*User, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(User)
	err := c.cc.Invoke(ctx, UserService_GetUserById_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetTutorById(ctx context.Context, in *GetById, opts ...grpc.CallOption) (*Tutor, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Tutor)
	err := c.cc.Invoke(ctx, UserService_GetTutorById_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetTutorInfoById(ctx context.Context, in *GetById, opts ...grpc.CallOption) (*TutorDetails, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TutorDetails)
	err := c.cc.Invoke(ctx, UserService_GetTutorInfoById_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetStudentById(ctx context.Context, in *GetById, opts ...grpc.CallOption) (*User, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(User)
	err := c.cc.Invoke(ctx, UserService_GetStudentById_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetUserByTelegramId(ctx context.Context, in *GetByTelegramId, opts ...grpc.CallOption) (*User, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(User)
	err := c.cc.Invoke(ctx, UserService_GetUserByTelegramId_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetAllUsers(ctx context.Context, in *GetAllRequest, opts ...grpc.CallOption) (*GetAllResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetAllResponse)
	err := c.cc.Invoke(ctx, UserService_GetAllUsers_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetAllTutorsPagination(ctx context.Context, in *GetAllTutorsPaginationRequest, opts ...grpc.CallOption) (*GetTutorsPaginationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetTutorsPaginationResponse)
	err := c.cc.Invoke(ctx, UserService_GetAllTutorsPagination_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) UpdateBioTutor(ctx context.Context, in *UpdateBioRequest, opts ...grpc.CallOption) (*Success, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Success)
	err := c.cc.Invoke(ctx, UserService_UpdateBioTutor_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) UpdateTags(ctx context.Context, in *UpdateTagsRequest, opts ...grpc.CallOption) (*Success, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Success)
	err := c.cc.Invoke(ctx, UserService_UpdateTags_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) ChangeTutorActive(ctx context.Context, in *SetActiveTutorById, opts ...grpc.CallOption) (*Success, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Success)
	err := c.cc.Invoke(ctx, UserService_ChangeTutorActive_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) ChangeTutorName(ctx context.Context, in *ChangeNameRequest, opts ...grpc.CallOption) (*Success, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Success)
	err := c.cc.Invoke(ctx, UserService_ChangeTutorName_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) CreateNewResponse(ctx context.Context, in *CreateResponseRequest, opts ...grpc.CallOption) (*Success, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Success)
	err := c.cc.Invoke(ctx, UserService_CreateNewResponse_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) AddResponsesToTutor(ctx context.Context, in *AddResponseToTutorRequest, opts ...grpc.CallOption) (*AddResponseToTutorResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddResponseToTutorResponse)
	err := c.cc.Invoke(ctx, UserService_AddResponsesToTutor_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) CreateReview(ctx context.Context, in *CreateReviewRequest, opts ...grpc.CallOption) (*CreateReviewResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateReviewResponse)
	err := c.cc.Invoke(ctx, UserService_CreateReview_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetReview(ctx context.Context, in *GetReviewRequest, opts ...grpc.CallOption) (*Review, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Review)
	err := c.cc.Invoke(ctx, UserService_GetReview_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetReviews(ctx context.Context, in *GetReviewsRequest, opts ...grpc.CallOption) (*GetReviewsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetReviewsResponse)
	err := c.cc.Invoke(ctx, UserService_GetReviews_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) SetReviewActive(ctx context.Context, in *SetReviewsActiveRequest, opts ...grpc.CallOption) (*Success, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Success)
	err := c.cc.Invoke(ctx, UserService_SetReviewActive_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) BanUser(ctx context.Context, in *BanUserRequest, opts ...grpc.CallOption) (*Success, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Success)
	err := c.cc.Invoke(ctx, UserService_BanUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServiceServer is the server API for UserService service.
// All implementations must embed UnimplementedUserServiceServer
// for forward compatibility.
type UserServiceServer interface {
	CreateUser(context.Context, *CreateUserRequest) (*CreateResponse, error)
	GetUserById(context.Context, *GetById) (*User, error)
	GetTutorById(context.Context, *GetById) (*Tutor, error)
	GetTutorInfoById(context.Context, *GetById) (*TutorDetails, error)
	GetStudentById(context.Context, *GetById) (*User, error)
	GetUserByTelegramId(context.Context, *GetByTelegramId) (*User, error)
	GetAllUsers(context.Context, *GetAllRequest) (*GetAllResponse, error)
	GetAllTutorsPagination(context.Context, *GetAllTutorsPaginationRequest) (*GetTutorsPaginationResponse, error)
	UpdateBioTutor(context.Context, *UpdateBioRequest) (*Success, error)
	UpdateTags(context.Context, *UpdateTagsRequest) (*Success, error)
	ChangeTutorActive(context.Context, *SetActiveTutorById) (*Success, error)
	ChangeTutorName(context.Context, *ChangeNameRequest) (*Success, error)
	CreateNewResponse(context.Context, *CreateResponseRequest) (*Success, error)
	AddResponsesToTutor(context.Context, *AddResponseToTutorRequest) (*AddResponseToTutorResponse, error)
	CreateReview(context.Context, *CreateReviewRequest) (*CreateReviewResponse, error)
	GetReview(context.Context, *GetReviewRequest) (*Review, error)
	GetReviews(context.Context, *GetReviewsRequest) (*GetReviewsResponse, error)
	SetReviewActive(context.Context, *SetReviewsActiveRequest) (*Success, error)
	BanUser(context.Context, *BanUserRequest) (*Success, error)
	mustEmbedUnimplementedUserServiceServer()
}

// UnimplementedUserServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedUserServiceServer struct{}

func (UnimplementedUserServiceServer) CreateUser(context.Context, *CreateUserRequest) (*CreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (UnimplementedUserServiceServer) GetUserById(context.Context, *GetById) (*User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserById not implemented")
}
func (UnimplementedUserServiceServer) GetTutorById(context.Context, *GetById) (*Tutor, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTutorById not implemented")
}
func (UnimplementedUserServiceServer) GetTutorInfoById(context.Context, *GetById) (*TutorDetails, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTutorInfoById not implemented")
}
func (UnimplementedUserServiceServer) GetStudentById(context.Context, *GetById) (*User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStudentById not implemented")
}
func (UnimplementedUserServiceServer) GetUserByTelegramId(context.Context, *GetByTelegramId) (*User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserByTelegramId not implemented")
}
func (UnimplementedUserServiceServer) GetAllUsers(context.Context, *GetAllRequest) (*GetAllResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllUsers not implemented")
}
func (UnimplementedUserServiceServer) GetAllTutorsPagination(context.Context, *GetAllTutorsPaginationRequest) (*GetTutorsPaginationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllTutorsPagination not implemented")
}
func (UnimplementedUserServiceServer) UpdateBioTutor(context.Context, *UpdateBioRequest) (*Success, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateBioTutor not implemented")
}
func (UnimplementedUserServiceServer) UpdateTags(context.Context, *UpdateTagsRequest) (*Success, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTags not implemented")
}
func (UnimplementedUserServiceServer) ChangeTutorActive(context.Context, *SetActiveTutorById) (*Success, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeTutorActive not implemented")
}
func (UnimplementedUserServiceServer) ChangeTutorName(context.Context, *ChangeNameRequest) (*Success, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeTutorName not implemented")
}
func (UnimplementedUserServiceServer) CreateNewResponse(context.Context, *CreateResponseRequest) (*Success, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateNewResponse not implemented")
}
func (UnimplementedUserServiceServer) AddResponsesToTutor(context.Context, *AddResponseToTutorRequest) (*AddResponseToTutorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddResponsesToTutor not implemented")
}
func (UnimplementedUserServiceServer) CreateReview(context.Context, *CreateReviewRequest) (*CreateReviewResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateReview not implemented")
}
func (UnimplementedUserServiceServer) GetReview(context.Context, *GetReviewRequest) (*Review, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetReview not implemented")
}
func (UnimplementedUserServiceServer) GetReviews(context.Context, *GetReviewsRequest) (*GetReviewsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetReviews not implemented")
}
func (UnimplementedUserServiceServer) SetReviewActive(context.Context, *SetReviewsActiveRequest) (*Success, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetReviewActive not implemented")
}
func (UnimplementedUserServiceServer) BanUser(context.Context, *BanUserRequest) (*Success, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BanUser not implemented")
}
func (UnimplementedUserServiceServer) mustEmbedUnimplementedUserServiceServer() {}
func (UnimplementedUserServiceServer) testEmbeddedByValue()                     {}

// UnsafeUserServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServiceServer will
// result in compilation errors.
type UnsafeUserServiceServer interface {
	mustEmbedUnimplementedUserServiceServer()
}

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	// If the following call pancis, it indicates UnimplementedUserServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&UserService_ServiceDesc, srv)
}

func _UserService_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_CreateUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetUserById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetById)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetUserById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_GetUserById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetUserById(ctx, req.(*GetById))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetTutorById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetById)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetTutorById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_GetTutorById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetTutorById(ctx, req.(*GetById))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetTutorInfoById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetById)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetTutorInfoById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_GetTutorInfoById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetTutorInfoById(ctx, req.(*GetById))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetStudentById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetById)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetStudentById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_GetStudentById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetStudentById(ctx, req.(*GetById))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetUserByTelegramId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByTelegramId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetUserByTelegramId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_GetUserByTelegramId_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetUserByTelegramId(ctx, req.(*GetByTelegramId))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetAllUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetAllUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_GetAllUsers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetAllUsers(ctx, req.(*GetAllRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetAllTutorsPagination_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllTutorsPaginationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetAllTutorsPagination(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_GetAllTutorsPagination_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetAllTutorsPagination(ctx, req.(*GetAllTutorsPaginationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_UpdateBioTutor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateBioRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).UpdateBioTutor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_UpdateBioTutor_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).UpdateBioTutor(ctx, req.(*UpdateBioRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_UpdateTags_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateTagsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).UpdateTags(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_UpdateTags_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).UpdateTags(ctx, req.(*UpdateTagsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_ChangeTutorActive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetActiveTutorById)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).ChangeTutorActive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_ChangeTutorActive_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).ChangeTutorActive(ctx, req.(*SetActiveTutorById))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_ChangeTutorName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangeNameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).ChangeTutorName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_ChangeTutorName_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).ChangeTutorName(ctx, req.(*ChangeNameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_CreateNewResponse_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateResponseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).CreateNewResponse(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_CreateNewResponse_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).CreateNewResponse(ctx, req.(*CreateResponseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_AddResponsesToTutor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddResponseToTutorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).AddResponsesToTutor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_AddResponsesToTutor_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).AddResponsesToTutor(ctx, req.(*AddResponseToTutorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_CreateReview_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateReviewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).CreateReview(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_CreateReview_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).CreateReview(ctx, req.(*CreateReviewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetReview_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetReviewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetReview(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_GetReview_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetReview(ctx, req.(*GetReviewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetReviews_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetReviewsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetReviews(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_GetReviews_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetReviews(ctx, req.(*GetReviewsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_SetReviewActive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetReviewsActiveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).SetReviewActive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_SetReviewActive_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).SetReviewActive(ctx, req.(*SetReviewsActiveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_BanUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BanUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).BanUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_BanUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).BanUser(ctx, req.(*BanUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserService_ServiceDesc is the grpc.ServiceDesc for UserService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "users.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateUser",
			Handler:    _UserService_CreateUser_Handler,
		},
		{
			MethodName: "GetUserById",
			Handler:    _UserService_GetUserById_Handler,
		},
		{
			MethodName: "GetTutorById",
			Handler:    _UserService_GetTutorById_Handler,
		},
		{
			MethodName: "GetTutorInfoById",
			Handler:    _UserService_GetTutorInfoById_Handler,
		},
		{
			MethodName: "GetStudentById",
			Handler:    _UserService_GetStudentById_Handler,
		},
		{
			MethodName: "GetUserByTelegramId",
			Handler:    _UserService_GetUserByTelegramId_Handler,
		},
		{
			MethodName: "GetAllUsers",
			Handler:    _UserService_GetAllUsers_Handler,
		},
		{
			MethodName: "GetAllTutorsPagination",
			Handler:    _UserService_GetAllTutorsPagination_Handler,
		},
		{
			MethodName: "UpdateBioTutor",
			Handler:    _UserService_UpdateBioTutor_Handler,
		},
		{
			MethodName: "UpdateTags",
			Handler:    _UserService_UpdateTags_Handler,
		},
		{
			MethodName: "ChangeTutorActive",
			Handler:    _UserService_ChangeTutorActive_Handler,
		},
		{
			MethodName: "ChangeTutorName",
			Handler:    _UserService_ChangeTutorName_Handler,
		},
		{
			MethodName: "CreateNewResponse",
			Handler:    _UserService_CreateNewResponse_Handler,
		},
		{
			MethodName: "AddResponsesToTutor",
			Handler:    _UserService_AddResponsesToTutor_Handler,
		},
		{
			MethodName: "CreateReview",
			Handler:    _UserService_CreateReview_Handler,
		},
		{
			MethodName: "GetReview",
			Handler:    _UserService_GetReview_Handler,
		},
		{
			MethodName: "GetReviews",
			Handler:    _UserService_GetReviews_Handler,
		},
		{
			MethodName: "SetReviewActive",
			Handler:    _UserService_SetReviewActive_Handler,
		},
		{
			MethodName: "BanUser",
			Handler:    _UserService_BanUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "tutors.proto",
}
