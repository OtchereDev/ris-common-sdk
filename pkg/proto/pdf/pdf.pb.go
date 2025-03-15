// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v3.12.4
// source: pkg/proto/pdf.proto

package pdf

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

type PaymentItem struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Code          string                 `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Description   string                 `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Quantity      int32                  `protobuf:"varint,3,opt,name=quantity,proto3" json:"quantity,omitempty"`
	UnitPrice     string                 `protobuf:"bytes,4,opt,name=unitPrice,proto3" json:"unitPrice,omitempty"`
	Amount        string                 `protobuf:"bytes,5,opt,name=amount,proto3" json:"amount,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PaymentItem) Reset() {
	*x = PaymentItem{}
	mi := &file_pkg_proto_pdf_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PaymentItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PaymentItem) ProtoMessage() {}

func (x *PaymentItem) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_pdf_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PaymentItem.ProtoReflect.Descriptor instead.
func (*PaymentItem) Descriptor() ([]byte, []int) {
	return file_pkg_proto_pdf_proto_rawDescGZIP(), []int{0}
}

func (x *PaymentItem) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *PaymentItem) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *PaymentItem) GetQuantity() int32 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

func (x *PaymentItem) GetUnitPrice() string {
	if x != nil {
		return x.UnitPrice
	}
	return ""
}

func (x *PaymentItem) GetAmount() string {
	if x != nil {
		return x.Amount
	}
	return ""
}

type ReceiptData struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	FileName      string                 `protobuf:"bytes,1,opt,name=fileName,proto3" json:"fileName,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	PatientId     string                 `protobuf:"bytes,3,opt,name=patientId,proto3" json:"patientId,omitempty"`
	PaymentDate   string                 `protobuf:"bytes,4,opt,name=paymentDate,proto3" json:"paymentDate,omitempty"`
	ReceiptNumber string                 `protobuf:"bytes,5,opt,name=receiptNumber,proto3" json:"receiptNumber,omitempty"`
	PaymentMethod string                 `protobuf:"bytes,6,opt,name=paymentMethod,proto3" json:"paymentMethod,omitempty"`
	Discount      string                 `protobuf:"bytes,7,opt,name=discount,proto3" json:"discount,omitempty"`
	SubTotal      string                 `protobuf:"bytes,8,opt,name=subTotal,proto3" json:"subTotal,omitempty"`
	Total         string                 `protobuf:"bytes,9,opt,name=total,proto3" json:"total,omitempty"`
	PaymentItems  []*PaymentItem         `protobuf:"bytes,10,rep,name=paymentItems,proto3" json:"paymentItems,omitempty"`
	Id            uint32                 `protobuf:"varint,11,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ReceiptData) Reset() {
	*x = ReceiptData{}
	mi := &file_pkg_proto_pdf_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReceiptData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReceiptData) ProtoMessage() {}

func (x *ReceiptData) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_pdf_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReceiptData.ProtoReflect.Descriptor instead.
func (*ReceiptData) Descriptor() ([]byte, []int) {
	return file_pkg_proto_pdf_proto_rawDescGZIP(), []int{1}
}

func (x *ReceiptData) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *ReceiptData) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ReceiptData) GetPatientId() string {
	if x != nil {
		return x.PatientId
	}
	return ""
}

func (x *ReceiptData) GetPaymentDate() string {
	if x != nil {
		return x.PaymentDate
	}
	return ""
}

func (x *ReceiptData) GetReceiptNumber() string {
	if x != nil {
		return x.ReceiptNumber
	}
	return ""
}

func (x *ReceiptData) GetPaymentMethod() string {
	if x != nil {
		return x.PaymentMethod
	}
	return ""
}

func (x *ReceiptData) GetDiscount() string {
	if x != nil {
		return x.Discount
	}
	return ""
}

func (x *ReceiptData) GetSubTotal() string {
	if x != nil {
		return x.SubTotal
	}
	return ""
}

func (x *ReceiptData) GetTotal() string {
	if x != nil {
		return x.Total
	}
	return ""
}

func (x *ReceiptData) GetPaymentItems() []*PaymentItem {
	if x != nil {
		return x.PaymentItems
	}
	return nil
}

func (x *ReceiptData) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

type RequestFormData struct {
	state              protoimpl.MessageState `protogen:"open.v1"`
	Id                 uint32                 `protobuf:"varint,12,opt,name=id,proto3" json:"id,omitempty"`
	PatientName        string                 `protobuf:"bytes,1,opt,name=patientName,proto3" json:"patientName,omitempty"`
	Sex                string                 `protobuf:"bytes,2,opt,name=sex,proto3" json:"sex,omitempty"`
	Date               string                 `protobuf:"bytes,3,opt,name=date,proto3" json:"date,omitempty"`
	Age                string                 `protobuf:"bytes,4,opt,name=age,proto3" json:"age,omitempty"`
	PhoneNumber        string                 `protobuf:"bytes,5,opt,name=phoneNumber,proto3" json:"phoneNumber,omitempty"`
	Address            string                 `protobuf:"bytes,6,opt,name=address,proto3" json:"address,omitempty"`
	RequestingDoctor   string                 `protobuf:"bytes,7,opt,name=requestingDoctor,proto3" json:"requestingDoctor,omitempty"`
	RequestingFacility string                 `protobuf:"bytes,8,opt,name=requestingFacility,proto3" json:"requestingFacility,omitempty"`
	Examination        string                 `protobuf:"bytes,9,opt,name=examination,proto3" json:"examination,omitempty"`
	Query              string                 `protobuf:"bytes,10,opt,name=query,proto3" json:"query,omitempty"`
	FileName           string                 `protobuf:"bytes,11,opt,name=fileName,proto3" json:"fileName,omitempty"`
	unknownFields      protoimpl.UnknownFields
	sizeCache          protoimpl.SizeCache
}

func (x *RequestFormData) Reset() {
	*x = RequestFormData{}
	mi := &file_pkg_proto_pdf_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RequestFormData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestFormData) ProtoMessage() {}

func (x *RequestFormData) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_pdf_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestFormData.ProtoReflect.Descriptor instead.
func (*RequestFormData) Descriptor() ([]byte, []int) {
	return file_pkg_proto_pdf_proto_rawDescGZIP(), []int{2}
}

func (x *RequestFormData) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *RequestFormData) GetPatientName() string {
	if x != nil {
		return x.PatientName
	}
	return ""
}

func (x *RequestFormData) GetSex() string {
	if x != nil {
		return x.Sex
	}
	return ""
}

func (x *RequestFormData) GetDate() string {
	if x != nil {
		return x.Date
	}
	return ""
}

func (x *RequestFormData) GetAge() string {
	if x != nil {
		return x.Age
	}
	return ""
}

func (x *RequestFormData) GetPhoneNumber() string {
	if x != nil {
		return x.PhoneNumber
	}
	return ""
}

func (x *RequestFormData) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *RequestFormData) GetRequestingDoctor() string {
	if x != nil {
		return x.RequestingDoctor
	}
	return ""
}

func (x *RequestFormData) GetRequestingFacility() string {
	if x != nil {
		return x.RequestingFacility
	}
	return ""
}

func (x *RequestFormData) GetExamination() string {
	if x != nil {
		return x.Examination
	}
	return ""
}

func (x *RequestFormData) GetQuery() string {
	if x != nil {
		return x.Query
	}
	return ""
}

func (x *RequestFormData) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

type ReportData struct {
	state             protoimpl.MessageState `protogen:"open.v1"`
	Id                uint32                 `protobuf:"varint,19,opt,name=id,proto3" json:"id,omitempty"`
	Name              string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Dob               string                 `protobuf:"bytes,2,opt,name=dob,proto3" json:"dob,omitempty"`
	Age               string                 `protobuf:"bytes,3,opt,name=age,proto3" json:"age,omitempty"`
	PatientId         string                 `protobuf:"bytes,4,opt,name=patientId,proto3" json:"patientId,omitempty"`
	ReferringDoctor   string                 `protobuf:"bytes,5,opt,name=referringDoctor,proto3" json:"referringDoctor,omitempty"`
	ReferringFacility string                 `protobuf:"bytes,6,opt,name=referringFacility,proto3" json:"referringFacility,omitempty"`
	Detail            string                 `protobuf:"bytes,7,opt,name=detail,proto3" json:"detail,omitempty"`
	Procedure         string                 `protobuf:"bytes,8,opt,name=procedure,proto3" json:"procedure,omitempty"`
	ExaminationDate   string                 `protobuf:"bytes,9,opt,name=examinationDate,proto3" json:"examinationDate,omitempty"`
	Radiologist       string                 `protobuf:"bytes,10,opt,name=radiologist,proto3" json:"radiologist,omitempty"`
	VerifiedAt        string                 `protobuf:"bytes,11,opt,name=verifiedAt,proto3" json:"verifiedAt,omitempty"`
	FileName          string                 `protobuf:"bytes,12,opt,name=fileName,proto3" json:"fileName,omitempty"`
	TitleOnReport     string                 `protobuf:"bytes,13,opt,name=titleOnReport,proto3" json:"titleOnReport,omitempty"`
	Impressions       string                 `protobuf:"bytes,14,opt,name=impressions,proto3" json:"impressions,omitempty"`
	Findings          string                 `protobuf:"bytes,15,opt,name=findings,proto3" json:"findings,omitempty"`
	ApprovedBy        string                 `protobuf:"bytes,16,opt,name=approvedBy,proto3" json:"approvedBy,omitempty"`
	ApproverTitle     string                 `protobuf:"bytes,17,opt,name=approverTitle,proto3" json:"approverTitle,omitempty"`
	WithHeader        string                 `protobuf:"bytes,18,opt,name=withHeader,proto3" json:"withHeader,omitempty"`
	unknownFields     protoimpl.UnknownFields
	sizeCache         protoimpl.SizeCache
}

func (x *ReportData) Reset() {
	*x = ReportData{}
	mi := &file_pkg_proto_pdf_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReportData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReportData) ProtoMessage() {}

func (x *ReportData) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_pdf_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReportData.ProtoReflect.Descriptor instead.
func (*ReportData) Descriptor() ([]byte, []int) {
	return file_pkg_proto_pdf_proto_rawDescGZIP(), []int{3}
}

func (x *ReportData) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *ReportData) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ReportData) GetDob() string {
	if x != nil {
		return x.Dob
	}
	return ""
}

func (x *ReportData) GetAge() string {
	if x != nil {
		return x.Age
	}
	return ""
}

func (x *ReportData) GetPatientId() string {
	if x != nil {
		return x.PatientId
	}
	return ""
}

func (x *ReportData) GetReferringDoctor() string {
	if x != nil {
		return x.ReferringDoctor
	}
	return ""
}

func (x *ReportData) GetReferringFacility() string {
	if x != nil {
		return x.ReferringFacility
	}
	return ""
}

func (x *ReportData) GetDetail() string {
	if x != nil {
		return x.Detail
	}
	return ""
}

func (x *ReportData) GetProcedure() string {
	if x != nil {
		return x.Procedure
	}
	return ""
}

func (x *ReportData) GetExaminationDate() string {
	if x != nil {
		return x.ExaminationDate
	}
	return ""
}

func (x *ReportData) GetRadiologist() string {
	if x != nil {
		return x.Radiologist
	}
	return ""
}

func (x *ReportData) GetVerifiedAt() string {
	if x != nil {
		return x.VerifiedAt
	}
	return ""
}

func (x *ReportData) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *ReportData) GetTitleOnReport() string {
	if x != nil {
		return x.TitleOnReport
	}
	return ""
}

func (x *ReportData) GetImpressions() string {
	if x != nil {
		return x.Impressions
	}
	return ""
}

func (x *ReportData) GetFindings() string {
	if x != nil {
		return x.Findings
	}
	return ""
}

func (x *ReportData) GetApprovedBy() string {
	if x != nil {
		return x.ApprovedBy
	}
	return ""
}

func (x *ReportData) GetApproverTitle() string {
	if x != nil {
		return x.ApproverTitle
	}
	return ""
}

func (x *ReportData) GetWithHeader() string {
	if x != nil {
		return x.WithHeader
	}
	return ""
}

type PdfEvent struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Event:
	//
	//	*PdfEvent_Receipt
	//	*PdfEvent_RequestForm
	//	*PdfEvent_Report
	Event         isPdfEvent_Event `protobuf_oneof:"event"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PdfEvent) Reset() {
	*x = PdfEvent{}
	mi := &file_pkg_proto_pdf_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PdfEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PdfEvent) ProtoMessage() {}

func (x *PdfEvent) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_pdf_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PdfEvent.ProtoReflect.Descriptor instead.
func (*PdfEvent) Descriptor() ([]byte, []int) {
	return file_pkg_proto_pdf_proto_rawDescGZIP(), []int{4}
}

func (x *PdfEvent) GetEvent() isPdfEvent_Event {
	if x != nil {
		return x.Event
	}
	return nil
}

func (x *PdfEvent) GetReceipt() *ReceiptData {
	if x != nil {
		if x, ok := x.Event.(*PdfEvent_Receipt); ok {
			return x.Receipt
		}
	}
	return nil
}

func (x *PdfEvent) GetRequestForm() *RequestFormData {
	if x != nil {
		if x, ok := x.Event.(*PdfEvent_RequestForm); ok {
			return x.RequestForm
		}
	}
	return nil
}

func (x *PdfEvent) GetReport() *ReportData {
	if x != nil {
		if x, ok := x.Event.(*PdfEvent_Report); ok {
			return x.Report
		}
	}
	return nil
}

type isPdfEvent_Event interface {
	isPdfEvent_Event()
}

type PdfEvent_Receipt struct {
	Receipt *ReceiptData `protobuf:"bytes,1,opt,name=receipt,proto3,oneof"`
}

type PdfEvent_RequestForm struct {
	RequestForm *RequestFormData `protobuf:"bytes,2,opt,name=requestForm,proto3,oneof"`
}

type PdfEvent_Report struct {
	Report *ReportData `protobuf:"bytes,3,opt,name=report,proto3,oneof"`
}

func (*PdfEvent_Receipt) isPdfEvent_Event() {}

func (*PdfEvent_RequestForm) isPdfEvent_Event() {}

func (*PdfEvent_Report) isPdfEvent_Event() {}

type FileUploaded struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	FileKey       string                 `protobuf:"bytes,1,opt,name=fileKey,proto3" json:"fileKey,omitempty"`
	Id            uint32                 `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FileUploaded) Reset() {
	*x = FileUploaded{}
	mi := &file_pkg_proto_pdf_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FileUploaded) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileUploaded) ProtoMessage() {}

func (x *FileUploaded) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_pdf_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileUploaded.ProtoReflect.Descriptor instead.
func (*FileUploaded) Descriptor() ([]byte, []int) {
	return file_pkg_proto_pdf_proto_rawDescGZIP(), []int{5}
}

func (x *FileUploaded) GetFileKey() string {
	if x != nil {
		return x.FileKey
	}
	return ""
}

func (x *FileUploaded) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

var File_pkg_proto_pdf_proto protoreflect.FileDescriptor

var file_pkg_proto_pdf_proto_rawDesc = string([]byte{
	0x0a, 0x13, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x64, 0x66, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x95, 0x01, 0x0a,
	0x0b, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x12, 0x0a, 0x04,
	0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65,
	0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x1c,
	0x0a, 0x09, 0x75, 0x6e, 0x69, 0x74, 0x50, 0x72, 0x69, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x75, 0x6e, 0x69, 0x74, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x6d,
	0x6f, 0x75, 0x6e, 0x74, 0x22, 0xdf, 0x02, 0x0a, 0x0b, 0x52, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74,
	0x44, 0x61, 0x74, 0x61, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x61, 0x74, 0x69, 0x65, 0x6e, 0x74, 0x49,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x61, 0x74, 0x69, 0x65, 0x6e, 0x74,
	0x49, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x44, 0x61, 0x74,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74,
	0x44, 0x61, 0x74, 0x65, 0x12, 0x24, 0x0a, 0x0d, 0x72, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x4e,
	0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x72, 0x65, 0x63,
	0x65, 0x69, 0x70, 0x74, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x24, 0x0a, 0x0d, 0x70, 0x61,
	0x79, 0x6d, 0x65, 0x6e, 0x74, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0d, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64,
	0x12, 0x1a, 0x0a, 0x08, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1a, 0x0a, 0x08,
	0x73, 0x75, 0x62, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x73, 0x75, 0x62, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x74, 0x61,
	0x6c, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x12, 0x36,
	0x0a, 0x0c, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x0a,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x61, 0x79,
	0x6d, 0x65, 0x6e, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x0c, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e,
	0x74, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x0b, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x02, 0x69, 0x64, 0x22, 0xe7, 0x02, 0x0a, 0x0f, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x46, 0x6f, 0x72, 0x6d, 0x44, 0x61, 0x74, 0x61, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x0c, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x69, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x70, 0x61,
	0x74, 0x69, 0x65, 0x6e, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x70, 0x61, 0x74, 0x69, 0x65, 0x6e, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03,
	0x73, 0x65, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x73, 0x65, 0x78, 0x12, 0x12,
	0x0a, 0x04, 0x64, 0x61, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61,
	0x74, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x61, 0x67, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x4e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x68, 0x6f, 0x6e, 0x65,
	0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x12, 0x2a, 0x0a, 0x10, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x44, 0x6f,
	0x63, 0x74, 0x6f, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x72, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x44, 0x6f, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x2e, 0x0a, 0x12,
	0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x46, 0x61, 0x63, 0x69, 0x6c, 0x69,
	0x74, 0x79, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x69, 0x6e, 0x67, 0x46, 0x61, 0x63, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x12, 0x20, 0x0a, 0x0b,
	0x65, 0x78, 0x61, 0x6d, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x65, 0x78, 0x61, 0x6d, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14,
	0x0a, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x71,
	0x75, 0x65, 0x72, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65,
	0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65,
	0x22, 0xd2, 0x04, 0x0a, 0x0a, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x44, 0x61, 0x74, 0x61, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x13, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x64, 0x6f, 0x62, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x64, 0x6f, 0x62, 0x12, 0x10, 0x0a, 0x03, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x61, 0x67, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x61, 0x74, 0x69, 0x65,
	0x6e, 0x74, 0x49, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x61, 0x74, 0x69,
	0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x28, 0x0a, 0x0f, 0x72, 0x65, 0x66, 0x65, 0x72, 0x72, 0x69,
	0x6e, 0x67, 0x44, 0x6f, 0x63, 0x74, 0x6f, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f,
	0x72, 0x65, 0x66, 0x65, 0x72, 0x72, 0x69, 0x6e, 0x67, 0x44, 0x6f, 0x63, 0x74, 0x6f, 0x72, 0x12,
	0x2c, 0x0a, 0x11, 0x72, 0x65, 0x66, 0x65, 0x72, 0x72, 0x69, 0x6e, 0x67, 0x46, 0x61, 0x63, 0x69,
	0x6c, 0x69, 0x74, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x72, 0x65, 0x66, 0x65,
	0x72, 0x72, 0x69, 0x6e, 0x67, 0x46, 0x61, 0x63, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x12, 0x16, 0x0a,
	0x06, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64,
	0x65, 0x74, 0x61, 0x69, 0x6c, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x64, 0x75,
	0x72, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x64,
	0x75, 0x72, 0x65, 0x12, 0x28, 0x0a, 0x0f, 0x65, 0x78, 0x61, 0x6d, 0x69, 0x6e, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x44, 0x61, 0x74, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x65, 0x78,
	0x61, 0x6d, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x61, 0x74, 0x65, 0x12, 0x20, 0x0a,
	0x0b, 0x72, 0x61, 0x64, 0x69, 0x6f, 0x6c, 0x6f, 0x67, 0x69, 0x73, 0x74, 0x18, 0x0a, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x72, 0x61, 0x64, 0x69, 0x6f, 0x6c, 0x6f, 0x67, 0x69, 0x73, 0x74, 0x12,
	0x1e, 0x0a, 0x0a, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64, 0x41, 0x74, 0x18, 0x0b, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64, 0x41, 0x74, 0x12,
	0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x24, 0x0a, 0x0d, 0x74,
	0x69, 0x74, 0x6c, 0x65, 0x4f, 0x6e, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x0d, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0d, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x4f, 0x6e, 0x52, 0x65, 0x70, 0x6f, 0x72,
	0x74, 0x12, 0x20, 0x0a, 0x0b, 0x69, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73,
	0x18, 0x0e, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x69, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x18,
	0x0f, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x12,
	0x1e, 0x0a, 0x0a, 0x61, 0x70, 0x70, 0x72, 0x6f, 0x76, 0x65, 0x64, 0x42, 0x79, 0x18, 0x10, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x61, 0x70, 0x70, 0x72, 0x6f, 0x76, 0x65, 0x64, 0x42, 0x79, 0x12,
	0x24, 0x0a, 0x0d, 0x61, 0x70, 0x70, 0x72, 0x6f, 0x76, 0x65, 0x72, 0x54, 0x69, 0x74, 0x6c, 0x65,
	0x18, 0x11, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x61, 0x70, 0x70, 0x72, 0x6f, 0x76, 0x65, 0x72,
	0x54, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x77, 0x69, 0x74, 0x68, 0x48, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x18, 0x12, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x77, 0x69, 0x74, 0x68, 0x48,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x22, 0xac, 0x01, 0x0a, 0x08, 0x50, 0x64, 0x66, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x12, 0x2e, 0x0a, 0x07, 0x72, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x63, 0x65,
	0x69, 0x70, 0x74, 0x44, 0x61, 0x74, 0x61, 0x48, 0x00, 0x52, 0x07, 0x72, 0x65, 0x63, 0x65, 0x69,
	0x70, 0x74, 0x12, 0x3a, 0x0a, 0x0b, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x46, 0x6f, 0x72,
	0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x46, 0x6f, 0x72, 0x6d, 0x44, 0x61, 0x74, 0x61, 0x48,
	0x00, 0x52, 0x0b, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x46, 0x6f, 0x72, 0x6d, 0x12, 0x2b,
	0x0a, 0x06, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x44, 0x61, 0x74,
	0x61, 0x48, 0x00, 0x52, 0x06, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x42, 0x07, 0x0a, 0x05, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x22, 0x38, 0x0a, 0x0c, 0x46, 0x69, 0x6c, 0x65, 0x55, 0x70, 0x6c, 0x6f,
	0x61, 0x64, 0x65, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x66, 0x69, 0x6c, 0x65, 0x4b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x66, 0x69, 0x6c, 0x65, 0x4b, 0x65, 0x79, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x69, 0x64, 0x42, 0x0f,
	0x5a, 0x0d, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x64, 0x66, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_pkg_proto_pdf_proto_rawDescOnce sync.Once
	file_pkg_proto_pdf_proto_rawDescData []byte
)

func file_pkg_proto_pdf_proto_rawDescGZIP() []byte {
	file_pkg_proto_pdf_proto_rawDescOnce.Do(func() {
		file_pkg_proto_pdf_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_pkg_proto_pdf_proto_rawDesc), len(file_pkg_proto_pdf_proto_rawDesc)))
	})
	return file_pkg_proto_pdf_proto_rawDescData
}

var file_pkg_proto_pdf_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_pkg_proto_pdf_proto_goTypes = []any{
	(*PaymentItem)(nil),     // 0: proto.PaymentItem
	(*ReceiptData)(nil),     // 1: proto.ReceiptData
	(*RequestFormData)(nil), // 2: proto.RequestFormData
	(*ReportData)(nil),      // 3: proto.ReportData
	(*PdfEvent)(nil),        // 4: proto.PdfEvent
	(*FileUploaded)(nil),    // 5: proto.FileUploaded
}
var file_pkg_proto_pdf_proto_depIdxs = []int32{
	0, // 0: proto.ReceiptData.paymentItems:type_name -> proto.PaymentItem
	1, // 1: proto.PdfEvent.receipt:type_name -> proto.ReceiptData
	2, // 2: proto.PdfEvent.requestForm:type_name -> proto.RequestFormData
	3, // 3: proto.PdfEvent.report:type_name -> proto.ReportData
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_pkg_proto_pdf_proto_init() }
func file_pkg_proto_pdf_proto_init() {
	if File_pkg_proto_pdf_proto != nil {
		return
	}
	file_pkg_proto_pdf_proto_msgTypes[4].OneofWrappers = []any{
		(*PdfEvent_Receipt)(nil),
		(*PdfEvent_RequestForm)(nil),
		(*PdfEvent_Report)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_pkg_proto_pdf_proto_rawDesc), len(file_pkg_proto_pdf_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_proto_pdf_proto_goTypes,
		DependencyIndexes: file_pkg_proto_pdf_proto_depIdxs,
		MessageInfos:      file_pkg_proto_pdf_proto_msgTypes,
	}.Build()
	File_pkg_proto_pdf_proto = out.File
	file_pkg_proto_pdf_proto_goTypes = nil
	file_pkg_proto_pdf_proto_depIdxs = nil
}
