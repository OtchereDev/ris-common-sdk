// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v4.25.3
// source: pkg/proto/system.proto

package system

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

type SmsLogin struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Phone         string                 `protobuf:"bytes,2,opt,name=phone,proto3" json:"phone,omitempty"`
	Otp           string                 `protobuf:"bytes,3,opt,name=otp,proto3" json:"otp,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SmsLogin) Reset() {
	*x = SmsLogin{}
	mi := &file_pkg_proto_system_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SmsLogin) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SmsLogin) ProtoMessage() {}

func (x *SmsLogin) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_system_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SmsLogin.ProtoReflect.Descriptor instead.
func (*SmsLogin) Descriptor() ([]byte, []int) {
	return file_pkg_proto_system_proto_rawDescGZIP(), []int{0}
}

func (x *SmsLogin) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *SmsLogin) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

func (x *SmsLogin) GetOtp() string {
	if x != nil {
		return x.Otp
	}
	return ""
}

type EmailDoctorForgotPassword struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Email         string                 `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Code          string                 `protobuf:"bytes,2,opt,name=code,proto3" json:"code,omitempty"`
	DoctorName    string                 `protobuf:"bytes,3,opt,name=doctorName,proto3" json:"doctorName,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EmailDoctorForgotPassword) Reset() {
	*x = EmailDoctorForgotPassword{}
	mi := &file_pkg_proto_system_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EmailDoctorForgotPassword) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmailDoctorForgotPassword) ProtoMessage() {}

func (x *EmailDoctorForgotPassword) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_system_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmailDoctorForgotPassword.ProtoReflect.Descriptor instead.
func (*EmailDoctorForgotPassword) Descriptor() ([]byte, []int) {
	return file_pkg_proto_system_proto_rawDescGZIP(), []int{1}
}

func (x *EmailDoctorForgotPassword) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *EmailDoctorForgotPassword) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *EmailDoctorForgotPassword) GetDoctorName() string {
	if x != nil {
		return x.DoctorName
	}
	return ""
}

type EmailReport struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Email         string                 `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Link          string                 `protobuf:"bytes,2,opt,name=link,proto3" json:"link,omitempty"`
	Name          string                 `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EmailReport) Reset() {
	*x = EmailReport{}
	mi := &file_pkg_proto_system_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EmailReport) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmailReport) ProtoMessage() {}

func (x *EmailReport) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_system_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmailReport.ProtoReflect.Descriptor instead.
func (*EmailReport) Descriptor() ([]byte, []int) {
	return file_pkg_proto_system_proto_rawDescGZIP(), []int{2}
}

func (x *EmailReport) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *EmailReport) GetLink() string {
	if x != nil {
		return x.Link
	}
	return ""
}

func (x *EmailReport) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type SmsForgotPassword struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Phone         string                 `protobuf:"bytes,2,opt,name=phone,proto3" json:"phone,omitempty"`
	Otp           string                 `protobuf:"bytes,3,opt,name=otp,proto3" json:"otp,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SmsForgotPassword) Reset() {
	*x = SmsForgotPassword{}
	mi := &file_pkg_proto_system_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SmsForgotPassword) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SmsForgotPassword) ProtoMessage() {}

func (x *SmsForgotPassword) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_system_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SmsForgotPassword.ProtoReflect.Descriptor instead.
func (*SmsForgotPassword) Descriptor() ([]byte, []int) {
	return file_pkg_proto_system_proto_rawDescGZIP(), []int{3}
}

func (x *SmsForgotPassword) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *SmsForgotPassword) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

func (x *SmsForgotPassword) GetOtp() string {
	if x != nil {
		return x.Otp
	}
	return ""
}

type LoginWithPhone struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Phone         string                 `protobuf:"bytes,2,opt,name=phone,proto3" json:"phone,omitempty"`
	Otp           string                 `protobuf:"bytes,3,opt,name=otp,proto3" json:"otp,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LoginWithPhone) Reset() {
	*x = LoginWithPhone{}
	mi := &file_pkg_proto_system_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoginWithPhone) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginWithPhone) ProtoMessage() {}

func (x *LoginWithPhone) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_system_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginWithPhone.ProtoReflect.Descriptor instead.
func (*LoginWithPhone) Descriptor() ([]byte, []int) {
	return file_pkg_proto_system_proto_rawDescGZIP(), []int{4}
}

func (x *LoginWithPhone) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *LoginWithPhone) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

func (x *LoginWithPhone) GetOtp() string {
	if x != nil {
		return x.Otp
	}
	return ""
}

type PushNotificationMessage struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Devices       []string               `protobuf:"bytes,1,rep,name=devices,proto3" json:"devices,omitempty"`
	Title         string                 `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Message       string                 `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PushNotificationMessage) Reset() {
	*x = PushNotificationMessage{}
	mi := &file_pkg_proto_system_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PushNotificationMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushNotificationMessage) ProtoMessage() {}

func (x *PushNotificationMessage) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_system_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushNotificationMessage.ProtoReflect.Descriptor instead.
func (*PushNotificationMessage) Descriptor() ([]byte, []int) {
	return file_pkg_proto_system_proto_rawDescGZIP(), []int{5}
}

func (x *PushNotificationMessage) GetDevices() []string {
	if x != nil {
		return x.Devices
	}
	return nil
}

func (x *PushNotificationMessage) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *PushNotificationMessage) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type EmailDoctorCreated struct {
	state             protoimpl.MessageState `protogen:"open.v1"`
	Name              string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Email             string                 `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	GeneratedPassword string                 `protobuf:"bytes,3,opt,name=generatedPassword,proto3" json:"generatedPassword,omitempty"`
	unknownFields     protoimpl.UnknownFields
	sizeCache         protoimpl.SizeCache
}

func (x *EmailDoctorCreated) Reset() {
	*x = EmailDoctorCreated{}
	mi := &file_pkg_proto_system_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EmailDoctorCreated) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmailDoctorCreated) ProtoMessage() {}

func (x *EmailDoctorCreated) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_system_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmailDoctorCreated.ProtoReflect.Descriptor instead.
func (*EmailDoctorCreated) Descriptor() ([]byte, []int) {
	return file_pkg_proto_system_proto_rawDescGZIP(), []int{6}
}

func (x *EmailDoctorCreated) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *EmailDoctorCreated) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *EmailDoctorCreated) GetGeneratedPassword() string {
	if x != nil {
		return x.GeneratedPassword
	}
	return ""
}

type WhatsappReport struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Phone         string                 `protobuf:"bytes,1,opt,name=phone,proto3" json:"phone,omitempty"`
	FirstName     string                 `protobuf:"bytes,2,opt,name=firstName,proto3" json:"firstName,omitempty"`
	Link          string                 `protobuf:"bytes,3,opt,name=link,proto3" json:"link,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *WhatsappReport) Reset() {
	*x = WhatsappReport{}
	mi := &file_pkg_proto_system_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *WhatsappReport) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WhatsappReport) ProtoMessage() {}

func (x *WhatsappReport) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_system_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WhatsappReport.ProtoReflect.Descriptor instead.
func (*WhatsappReport) Descriptor() ([]byte, []int) {
	return file_pkg_proto_system_proto_rawDescGZIP(), []int{7}
}

func (x *WhatsappReport) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

func (x *WhatsappReport) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}

func (x *WhatsappReport) GetLink() string {
	if x != nil {
		return x.Link
	}
	return ""
}

type SendSms struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	PhoneNumber   string                 `protobuf:"bytes,1,opt,name=phoneNumber,proto3" json:"phoneNumber,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SendSms) Reset() {
	*x = SendSms{}
	mi := &file_pkg_proto_system_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SendSms) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendSms) ProtoMessage() {}

func (x *SendSms) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_system_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendSms.ProtoReflect.Descriptor instead.
func (*SendSms) Descriptor() ([]byte, []int) {
	return file_pkg_proto_system_proto_rawDescGZIP(), []int{8}
}

func (x *SendSms) GetPhoneNumber() string {
	if x != nil {
		return x.PhoneNumber
	}
	return ""
}

func (x *SendSms) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_pkg_proto_system_proto protoreflect.FileDescriptor

var file_pkg_proto_system_proto_rawDesc = string([]byte{
	0x0a, 0x16, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x79, 0x73, 0x74,
	0x65, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x46, 0x0a, 0x08, 0x53, 0x6d, 0x73, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x14, 0x0a, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x70, 0x68, 0x6f, 0x6e, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6f, 0x74, 0x70, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6f, 0x74, 0x70, 0x22, 0x65, 0x0a, 0x19, 0x45, 0x6d, 0x61, 0x69, 0x6c,
	0x44, 0x6f, 0x63, 0x74, 0x6f, 0x72, 0x46, 0x6f, 0x72, 0x67, 0x6f, 0x74, 0x50, 0x61, 0x73, 0x73,
	0x77, 0x6f, 0x72, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f,
	0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x1e,
	0x0a, 0x0a, 0x64, 0x6f, 0x63, 0x74, 0x6f, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x64, 0x6f, 0x63, 0x74, 0x6f, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x4b,
	0x0a, 0x0b, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x14, 0x0a,
	0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d,
	0x61, 0x69, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x6c, 0x69, 0x6e, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6c, 0x69, 0x6e, 0x6b, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x4f, 0x0a, 0x11, 0x53,
	0x6d, 0x73, 0x46, 0x6f, 0x72, 0x67, 0x6f, 0x74, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6f, 0x74,
	0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6f, 0x74, 0x70, 0x22, 0x4c, 0x0a, 0x0e,
	0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x57, 0x69, 0x74, 0x68, 0x50, 0x68, 0x6f, 0x6e, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6f, 0x74, 0x70, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6f, 0x74, 0x70, 0x22, 0x63, 0x0a, 0x17, 0x50, 0x75,
	0x73, 0x68, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x73, 0x12,
	0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22,
	0x6c, 0x0a, 0x12, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x44, 0x6f, 0x63, 0x74, 0x6f, 0x72, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61,
	0x69, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12,
	0x2c, 0x0a, 0x11, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x50, 0x61, 0x73, 0x73,
	0x77, 0x6f, 0x72, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x67, 0x65, 0x6e, 0x65,
	0x72, 0x61, 0x74, 0x65, 0x64, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0x58, 0x0a,
	0x0e, 0x57, 0x68, 0x61, 0x74, 0x73, 0x61, 0x70, 0x70, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x12,
	0x14, 0x0a, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x70, 0x68, 0x6f, 0x6e, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x66, 0x69, 0x72, 0x73, 0x74, 0x4e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x66, 0x69, 0x72, 0x73, 0x74, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6c, 0x69, 0x6e, 0x6b, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6c, 0x69, 0x6e, 0x6b, 0x22, 0x45, 0x0a, 0x07, 0x53, 0x65, 0x6e, 0x64, 0x53,
	0x6d, 0x73, 0x12, 0x20, 0x0a, 0x0b, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x4e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x42, 0x12,
	0x5a, 0x10, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x79, 0x73, 0x74,
	0x65, 0x6d, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_pkg_proto_system_proto_rawDescOnce sync.Once
	file_pkg_proto_system_proto_rawDescData []byte
)

func file_pkg_proto_system_proto_rawDescGZIP() []byte {
	file_pkg_proto_system_proto_rawDescOnce.Do(func() {
		file_pkg_proto_system_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_pkg_proto_system_proto_rawDesc), len(file_pkg_proto_system_proto_rawDesc)))
	})
	return file_pkg_proto_system_proto_rawDescData
}

var file_pkg_proto_system_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_pkg_proto_system_proto_goTypes = []any{
	(*SmsLogin)(nil),                  // 0: proto.SmsLogin
	(*EmailDoctorForgotPassword)(nil), // 1: proto.EmailDoctorForgotPassword
	(*EmailReport)(nil),               // 2: proto.EmailReport
	(*SmsForgotPassword)(nil),         // 3: proto.SmsForgotPassword
	(*LoginWithPhone)(nil),            // 4: proto.LoginWithPhone
	(*PushNotificationMessage)(nil),   // 5: proto.PushNotificationMessage
	(*EmailDoctorCreated)(nil),        // 6: proto.EmailDoctorCreated
	(*WhatsappReport)(nil),            // 7: proto.WhatsappReport
	(*SendSms)(nil),                   // 8: proto.SendSms
}
var file_pkg_proto_system_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pkg_proto_system_proto_init() }
func file_pkg_proto_system_proto_init() {
	if File_pkg_proto_system_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_pkg_proto_system_proto_rawDesc), len(file_pkg_proto_system_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_proto_system_proto_goTypes,
		DependencyIndexes: file_pkg_proto_system_proto_depIdxs,
		MessageInfos:      file_pkg_proto_system_proto_msgTypes,
	}.Build()
	File_pkg_proto_system_proto = out.File
	file_pkg_proto_system_proto_goTypes = nil
	file_pkg_proto_system_proto_depIdxs = nil
}
