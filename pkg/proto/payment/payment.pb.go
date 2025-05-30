// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.3
// source: pkg/proto/payment.proto

package payment

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
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

type AppointmentPayment struct {
	state              protoimpl.MessageState `protogen:"open.v1"`
	Id                 uint32                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	AppointmentIds     []uint32               `protobuf:"varint,2,rep,packed,name=appointment_ids,json=appointmentIds,proto3" json:"appointment_ids,omitempty"`
	Amount             float64                `protobuf:"fixed64,3,opt,name=amount,proto3" json:"amount,omitempty"`
	HasDiscount        bool                   `protobuf:"varint,4,opt,name=has_discount,json=hasDiscount,proto3" json:"has_discount,omitempty"`
	IsReceipted        bool                   `protobuf:"varint,5,opt,name=is_receipted,json=isReceipted,proto3" json:"is_receipted,omitempty"`
	PaymentMethod      string                 `protobuf:"bytes,6,opt,name=payment_method,json=paymentMethod,proto3" json:"payment_method,omitempty"`
	DiscountAmount     float64                `protobuf:"fixed64,7,opt,name=discount_amount,json=discountAmount,proto3" json:"discount_amount,omitempty"`
	DiscountAuthorizer string                 `protobuf:"bytes,8,opt,name=discount_authorizer,json=discountAuthorizer,proto3" json:"discount_authorizer,omitempty"`
	PaidAt             *timestamppb.Timestamp `protobuf:"bytes,9,opt,name=paid_at,json=paidAt,proto3" json:"paid_at,omitempty"`
	IsDeleted          bool                   `protobuf:"varint,10,opt,name=is_deleted,json=isDeleted,proto3" json:"is_deleted,omitempty"`
	ReceiptUrl         string                 `protobuf:"bytes,11,opt,name=receipt_url,json=receiptUrl,proto3" json:"receipt_url,omitempty"`
	ReceiptNumber      string                 `protobuf:"bytes,12,opt,name=receipt_number,json=receiptNumber,proto3" json:"receipt_number,omitempty"`
	InsuranceCompany   string                 `protobuf:"bytes,13,opt,name=insurance_company,json=insuranceCompany,proto3" json:"insurance_company,omitempty"`
	InsuranceId        string                 `protobuf:"bytes,14,opt,name=insurance_id,json=insuranceId,proto3" json:"insurance_id,omitempty"`
	PaymentDetails     string                 `protobuf:"bytes,15,opt,name=payment_details,json=paymentDetails,proto3" json:"payment_details,omitempty"`
	unknownFields      protoimpl.UnknownFields
	sizeCache          protoimpl.SizeCache
}

func (x *AppointmentPayment) Reset() {
	*x = AppointmentPayment{}
	mi := &file_pkg_proto_payment_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AppointmentPayment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppointmentPayment) ProtoMessage() {}

func (x *AppointmentPayment) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_payment_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppointmentPayment.ProtoReflect.Descriptor instead.
func (*AppointmentPayment) Descriptor() ([]byte, []int) {
	return file_pkg_proto_payment_proto_rawDescGZIP(), []int{0}
}

func (x *AppointmentPayment) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *AppointmentPayment) GetAppointmentIds() []uint32 {
	if x != nil {
		return x.AppointmentIds
	}
	return nil
}

func (x *AppointmentPayment) GetAmount() float64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *AppointmentPayment) GetHasDiscount() bool {
	if x != nil {
		return x.HasDiscount
	}
	return false
}

func (x *AppointmentPayment) GetIsReceipted() bool {
	if x != nil {
		return x.IsReceipted
	}
	return false
}

func (x *AppointmentPayment) GetPaymentMethod() string {
	if x != nil {
		return x.PaymentMethod
	}
	return ""
}

func (x *AppointmentPayment) GetDiscountAmount() float64 {
	if x != nil {
		return x.DiscountAmount
	}
	return 0
}

func (x *AppointmentPayment) GetDiscountAuthorizer() string {
	if x != nil {
		return x.DiscountAuthorizer
	}
	return ""
}

func (x *AppointmentPayment) GetPaidAt() *timestamppb.Timestamp {
	if x != nil {
		return x.PaidAt
	}
	return nil
}

func (x *AppointmentPayment) GetIsDeleted() bool {
	if x != nil {
		return x.IsDeleted
	}
	return false
}

func (x *AppointmentPayment) GetReceiptUrl() string {
	if x != nil {
		return x.ReceiptUrl
	}
	return ""
}

func (x *AppointmentPayment) GetReceiptNumber() string {
	if x != nil {
		return x.ReceiptNumber
	}
	return ""
}

func (x *AppointmentPayment) GetInsuranceCompany() string {
	if x != nil {
		return x.InsuranceCompany
	}
	return ""
}

func (x *AppointmentPayment) GetInsuranceId() string {
	if x != nil {
		return x.InsuranceId
	}
	return ""
}

func (x *AppointmentPayment) GetPaymentDetails() string {
	if x != nil {
		return x.PaymentDetails
	}
	return ""
}

var File_pkg_proto_payment_proto protoreflect.FileDescriptor

var file_pkg_proto_payment_proto_rawDesc = string([]byte{
	0x0a, 0x17, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x61, 0x79, 0x6d,
	0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x61, 0x70, 0x70, 0x6f, 0x69,
	0x6e, 0x74, 0x6d, 0x65, 0x6e, 0x74, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc1, 0x04, 0x0a, 0x12, 0x41, 0x70, 0x70, 0x6f,
	0x69, 0x6e, 0x74, 0x6d, 0x65, 0x6e, 0x74, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x69, 0x64, 0x12, 0x27,
	0x0a, 0x0f, 0x61, 0x70, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x0e, 0x61, 0x70, 0x70, 0x6f, 0x69, 0x6e, 0x74,
	0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12,
	0x21, 0x0a, 0x0c, 0x68, 0x61, 0x73, 0x5f, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x68, 0x61, 0x73, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x69, 0x73, 0x5f, 0x72, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74,
	0x65, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x69, 0x73, 0x52, 0x65, 0x63, 0x65,
	0x69, 0x70, 0x74, 0x65, 0x64, 0x12, 0x25, 0x0a, 0x0e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74,
	0x5f, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x70,
	0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x27, 0x0a, 0x0f,
	0x64, 0x69, 0x73, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x41,
	0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x2f, 0x0a, 0x13, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x5f, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65, 0x72, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x12, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x41, 0x75, 0x74, 0x68,
	0x6f, 0x72, 0x69, 0x7a, 0x65, 0x72, 0x12, 0x33, 0x0a, 0x07, 0x70, 0x61, 0x69, 0x64, 0x5f, 0x61,
	0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x06, 0x70, 0x61, 0x69, 0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x69,
	0x73, 0x5f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x09, 0x69, 0x73, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65,
	0x63, 0x65, 0x69, 0x70, 0x74, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x72, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x55, 0x72, 0x6c, 0x12, 0x25, 0x0a, 0x0e, 0x72,
	0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x0c, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0d, 0x72, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x4e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x12, 0x2b, 0x0a, 0x11, 0x69, 0x6e, 0x73, 0x75, 0x72, 0x61, 0x6e, 0x63, 0x65, 0x5f,
	0x63, 0x6f, 0x6d, 0x70, 0x61, 0x6e, 0x79, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x69,
	0x6e, 0x73, 0x75, 0x72, 0x61, 0x6e, 0x63, 0x65, 0x43, 0x6f, 0x6d, 0x70, 0x61, 0x6e, 0x79, 0x12,
	0x21, 0x0a, 0x0c, 0x69, 0x6e, 0x73, 0x75, 0x72, 0x61, 0x6e, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18,
	0x0e, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x69, 0x6e, 0x73, 0x75, 0x72, 0x61, 0x6e, 0x63, 0x65,
	0x49, 0x64, 0x12, 0x27, 0x0a, 0x0f, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x64, 0x65,
	0x74, 0x61, 0x69, 0x6c, 0x73, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x70, 0x61, 0x79,
	0x6d, 0x65, 0x6e, 0x74, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x42, 0x13, 0x5a, 0x11, 0x70,
	0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_pkg_proto_payment_proto_rawDescOnce sync.Once
	file_pkg_proto_payment_proto_rawDescData []byte
)

func file_pkg_proto_payment_proto_rawDescGZIP() []byte {
	file_pkg_proto_payment_proto_rawDescOnce.Do(func() {
		file_pkg_proto_payment_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_pkg_proto_payment_proto_rawDesc), len(file_pkg_proto_payment_proto_rawDesc)))
	})
	return file_pkg_proto_payment_proto_rawDescData
}

var file_pkg_proto_payment_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_pkg_proto_payment_proto_goTypes = []any{
	(*AppointmentPayment)(nil),    // 0: appointment.AppointmentPayment
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
}
var file_pkg_proto_payment_proto_depIdxs = []int32{
	1, // 0: appointment.AppointmentPayment.paid_at:type_name -> google.protobuf.Timestamp
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_pkg_proto_payment_proto_init() }
func file_pkg_proto_payment_proto_init() {
	if File_pkg_proto_payment_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_pkg_proto_payment_proto_rawDesc), len(file_pkg_proto_payment_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_proto_payment_proto_goTypes,
		DependencyIndexes: file_pkg_proto_payment_proto_depIdxs,
		MessageInfos:      file_pkg_proto_payment_proto_msgTypes,
	}.Build()
	File_pkg_proto_payment_proto = out.File
	file_pkg_proto_payment_proto_goTypes = nil
	file_pkg_proto_payment_proto_depIdxs = nil
}
