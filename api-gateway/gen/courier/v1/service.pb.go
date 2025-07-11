// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: courier/v1/service.proto

package courier_v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RegisterRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Password      string                 `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	Phone         string                 `protobuf:"bytes,3,opt,name=phone,proto3" json:"phone,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegisterRequest) Reset() {
	*x = RegisterRequest{}
	mi := &file_courier_v1_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegisterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterRequest) ProtoMessage() {}

func (x *RegisterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_courier_v1_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterRequest.ProtoReflect.Descriptor instead.
func (*RegisterRequest) Descriptor() ([]byte, []int) {
	return file_courier_v1_service_proto_rawDescGZIP(), []int{0}
}

func (x *RegisterRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *RegisterRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *RegisterRequest) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

type RegisterResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CourierId     string                 `protobuf:"bytes,1,opt,name=courier_id,json=courierId,proto3" json:"courier_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegisterResponse) Reset() {
	*x = RegisterResponse{}
	mi := &file_courier_v1_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegisterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterResponse) ProtoMessage() {}

func (x *RegisterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_courier_v1_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterResponse.ProtoReflect.Descriptor instead.
func (*RegisterResponse) Descriptor() ([]byte, []int) {
	return file_courier_v1_service_proto_rawDescGZIP(), []int{1}
}

func (x *RegisterResponse) GetCourierId() string {
	if x != nil {
		return x.CourierId
	}
	return ""
}

type LoginRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Phone         string                 `protobuf:"bytes,1,opt,name=phone,proto3" json:"phone,omitempty"`
	Password      string                 `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LoginRequest) Reset() {
	*x = LoginRequest{}
	mi := &file_courier_v1_service_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoginRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginRequest) ProtoMessage() {}

func (x *LoginRequest) ProtoReflect() protoreflect.Message {
	mi := &file_courier_v1_service_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginRequest.ProtoReflect.Descriptor instead.
func (*LoginRequest) Descriptor() ([]byte, []int) {
	return file_courier_v1_service_proto_rawDescGZIP(), []int{2}
}

func (x *LoginRequest) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

func (x *LoginRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type LoginResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Token         string                 `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LoginResponse) Reset() {
	*x = LoginResponse{}
	mi := &file_courier_v1_service_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoginResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginResponse) ProtoMessage() {}

func (x *LoginResponse) ProtoReflect() protoreflect.Message {
	mi := &file_courier_v1_service_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginResponse.ProtoReflect.Descriptor instead.
func (*LoginResponse) Descriptor() ([]byte, []int) {
	return file_courier_v1_service_proto_rawDescGZIP(), []int{3}
}

func (x *LoginResponse) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type AuthenticateRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Token         string                 `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AuthenticateRequest) Reset() {
	*x = AuthenticateRequest{}
	mi := &file_courier_v1_service_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AuthenticateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthenticateRequest) ProtoMessage() {}

func (x *AuthenticateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_courier_v1_service_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthenticateRequest.ProtoReflect.Descriptor instead.
func (*AuthenticateRequest) Descriptor() ([]byte, []int) {
	return file_courier_v1_service_proto_rawDescGZIP(), []int{4}
}

func (x *AuthenticateRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type AuthenticateResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CourierId     string                 `protobuf:"bytes,1,opt,name=courier_id,json=courierId,proto3" json:"courier_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AuthenticateResponse) Reset() {
	*x = AuthenticateResponse{}
	mi := &file_courier_v1_service_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AuthenticateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthenticateResponse) ProtoMessage() {}

func (x *AuthenticateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_courier_v1_service_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthenticateResponse.ProtoReflect.Descriptor instead.
func (*AuthenticateResponse) Descriptor() ([]byte, []int) {
	return file_courier_v1_service_proto_rawDescGZIP(), []int{5}
}

func (x *AuthenticateResponse) GetCourierId() string {
	if x != nil {
		return x.CourierId
	}
	return ""
}

var File_courier_v1_service_proto protoreflect.FileDescriptor

const file_courier_v1_service_proto_rawDesc = "" +
	"\n" +
	"\x18courier/v1/service.proto\x12\n" +
	"courier.v1\"W\n" +
	"\x0fRegisterRequest\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x12\x1a\n" +
	"\bpassword\x18\x02 \x01(\tR\bpassword\x12\x14\n" +
	"\x05phone\x18\x03 \x01(\tR\x05phone\"1\n" +
	"\x10RegisterResponse\x12\x1d\n" +
	"\n" +
	"courier_id\x18\x01 \x01(\tR\tcourierId\"@\n" +
	"\fLoginRequest\x12\x14\n" +
	"\x05phone\x18\x01 \x01(\tR\x05phone\x12\x1a\n" +
	"\bpassword\x18\x02 \x01(\tR\bpassword\"%\n" +
	"\rLoginResponse\x12\x14\n" +
	"\x05token\x18\x01 \x01(\tR\x05token\"+\n" +
	"\x13AuthenticateRequest\x12\x14\n" +
	"\x05token\x18\x01 \x01(\tR\x05token\"5\n" +
	"\x14AuthenticateResponse\x12\x1d\n" +
	"\n" +
	"courier_id\x18\x01 \x01(\tR\tcourierId2\xec\x01\n" +
	"\x12CourierAuthService\x12E\n" +
	"\bRegister\x12\x1b.courier.v1.RegisterRequest\x1a\x1c.courier.v1.RegisterResponse\x12<\n" +
	"\x05Login\x12\x18.courier.v1.LoginRequest\x1a\x19.courier.v1.LoginResponse\x12Q\n" +
	"\fAuthenticate\x12\x1f.courier.v1.AuthenticateRequest\x1a .courier.v1.AuthenticateResponseBPZNgithub.com/OlegDokuchaev/clean-ddd-app/api-gateway/proto/courier/v1;courier_v1b\x06proto3"

var (
	file_courier_v1_service_proto_rawDescOnce sync.Once
	file_courier_v1_service_proto_rawDescData []byte
)

func file_courier_v1_service_proto_rawDescGZIP() []byte {
	file_courier_v1_service_proto_rawDescOnce.Do(func() {
		file_courier_v1_service_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_courier_v1_service_proto_rawDesc), len(file_courier_v1_service_proto_rawDesc)))
	})
	return file_courier_v1_service_proto_rawDescData
}

var file_courier_v1_service_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_courier_v1_service_proto_goTypes = []any{
	(*RegisterRequest)(nil),      // 0: courier.v1.RegisterRequest
	(*RegisterResponse)(nil),     // 1: courier.v1.RegisterResponse
	(*LoginRequest)(nil),         // 2: courier.v1.LoginRequest
	(*LoginResponse)(nil),        // 3: courier.v1.LoginResponse
	(*AuthenticateRequest)(nil),  // 4: courier.v1.AuthenticateRequest
	(*AuthenticateResponse)(nil), // 5: courier.v1.AuthenticateResponse
}
var file_courier_v1_service_proto_depIdxs = []int32{
	0, // 0: courier.v1.CourierAuthService.Register:input_type -> courier.v1.RegisterRequest
	2, // 1: courier.v1.CourierAuthService.Login:input_type -> courier.v1.LoginRequest
	4, // 2: courier.v1.CourierAuthService.Authenticate:input_type -> courier.v1.AuthenticateRequest
	1, // 3: courier.v1.CourierAuthService.Register:output_type -> courier.v1.RegisterResponse
	3, // 4: courier.v1.CourierAuthService.Login:output_type -> courier.v1.LoginResponse
	5, // 5: courier.v1.CourierAuthService.Authenticate:output_type -> courier.v1.AuthenticateResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_courier_v1_service_proto_init() }
func file_courier_v1_service_proto_init() {
	if File_courier_v1_service_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_courier_v1_service_proto_rawDesc), len(file_courier_v1_service_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_courier_v1_service_proto_goTypes,
		DependencyIndexes: file_courier_v1_service_proto_depIdxs,
		MessageInfos:      file_courier_v1_service_proto_msgTypes,
	}.Build()
	File_courier_v1_service_proto = out.File
	file_courier_v1_service_proto_goTypes = nil
	file_courier_v1_service_proto_depIdxs = nil
}
