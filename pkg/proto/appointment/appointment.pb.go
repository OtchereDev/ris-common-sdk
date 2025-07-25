// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.12.4
// source: pkg/proto/appointment.proto

package appointment

import (
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

type Appointment struct {
	state             protoimpl.MessageState `protogen:"open.v1"`
	Id                uint32                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	PatientId         uint32                 `protobuf:"varint,2,opt,name=patient_id,json=patientId,proto3" json:"patient_id,omitempty"`
	IsEmergency       bool                   `protobuf:"varint,3,opt,name=is_emergency,json=isEmergency,proto3" json:"is_emergency,omitempty"`
	ProcedureId       uint32                 `protobuf:"varint,4,opt,name=procedure_id,json=procedureId,proto3" json:"procedure_id,omitempty"`
	ReportId          uint32                 `protobuf:"varint,5,opt,name=report_id,json=reportId,proto3" json:"report_id,omitempty"`
	ReferringDoctorId uint32                 `protobuf:"varint,6,opt,name=referring_doctor_id,json=referringDoctorId,proto3" json:"referring_doctor_id,omitempty"`
	ReportingDoctorId uint32                 `protobuf:"varint,7,opt,name=reporting_doctor_id,json=reportingDoctorId,proto3" json:"reporting_doctor_id,omitempty"`
	AppointmentTime   *timestamp.Timestamp   `protobuf:"bytes,8,opt,name=appointment_time,json=appointmentTime,proto3" json:"appointment_time,omitempty"`
	PaymentId         uint32                 `protobuf:"varint,9,opt,name=payment_id,json=paymentId,proto3" json:"payment_id,omitempty"`
	Note              string                 `protobuf:"bytes,10,opt,name=note,proto3" json:"note,omitempty"`
	Status            string                 `protobuf:"bytes,11,opt,name=status,proto3" json:"status,omitempty"`
	HasPaidCommission bool                   `protobuf:"varint,12,opt,name=has_paid_commission,json=hasPaidCommission,proto3" json:"has_paid_commission,omitempty"`
	OrganizationId    uint32                 `protobuf:"varint,13,opt,name=organization_id,json=organizationId,proto3" json:"organization_id,omitempty"`
	unknownFields     protoimpl.UnknownFields
	sizeCache         protoimpl.SizeCache
}

func (x *Appointment) Reset() {
	*x = Appointment{}
	mi := &file_pkg_proto_appointment_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Appointment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Appointment) ProtoMessage() {}

func (x *Appointment) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_appointment_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Appointment.ProtoReflect.Descriptor instead.
func (*Appointment) Descriptor() ([]byte, []int) {
	return file_pkg_proto_appointment_proto_rawDescGZIP(), []int{0}
}

func (x *Appointment) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Appointment) GetPatientId() uint32 {
	if x != nil {
		return x.PatientId
	}
	return 0
}

func (x *Appointment) GetIsEmergency() bool {
	if x != nil {
		return x.IsEmergency
	}
	return false
}

func (x *Appointment) GetProcedureId() uint32 {
	if x != nil {
		return x.ProcedureId
	}
	return 0
}

func (x *Appointment) GetReportId() uint32 {
	if x != nil {
		return x.ReportId
	}
	return 0
}

func (x *Appointment) GetReferringDoctorId() uint32 {
	if x != nil {
		return x.ReferringDoctorId
	}
	return 0
}

func (x *Appointment) GetReportingDoctorId() uint32 {
	if x != nil {
		return x.ReportingDoctorId
	}
	return 0
}

func (x *Appointment) GetAppointmentTime() *timestamp.Timestamp {
	if x != nil {
		return x.AppointmentTime
	}
	return nil
}

func (x *Appointment) GetPaymentId() uint32 {
	if x != nil {
		return x.PaymentId
	}
	return 0
}

func (x *Appointment) GetNote() string {
	if x != nil {
		return x.Note
	}
	return ""
}

func (x *Appointment) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *Appointment) GetHasPaidCommission() bool {
	if x != nil {
		return x.HasPaidCommission
	}
	return false
}

func (x *Appointment) GetOrganizationId() uint32 {
	if x != nil {
		return x.OrganizationId
	}
	return 0
}

type Report struct {
	state             protoimpl.MessageState `protogen:"open.v1"`
	Id                uint32                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	AppointmentId     uint32                 `protobuf:"varint,2,opt,name=appointment_id,json=appointmentId,proto3" json:"appointment_id,omitempty"`
	ReferringHospital string                 `protobuf:"bytes,3,opt,name=referring_hospital,json=referringHospital,proto3" json:"referring_hospital,omitempty"`
	Findings          string                 `protobuf:"bytes,4,opt,name=findings,proto3" json:"findings,omitempty"`
	Impressions       string                 `protobuf:"bytes,5,opt,name=impressions,proto3" json:"impressions,omitempty"`
	ClinicalDetails   string                 `protobuf:"bytes,6,opt,name=clinical_details,json=clinicalDetails,proto3" json:"clinical_details,omitempty"`
	Url               string                 `protobuf:"bytes,7,opt,name=url,proto3" json:"url,omitempty"`
	Status            string                 `protobuf:"bytes,8,opt,name=status,proto3" json:"status,omitempty"`
	unknownFields     protoimpl.UnknownFields
	sizeCache         protoimpl.SizeCache
}

func (x *Report) Reset() {
	*x = Report{}
	mi := &file_pkg_proto_appointment_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Report) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Report) ProtoMessage() {}

func (x *Report) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_appointment_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Report.ProtoReflect.Descriptor instead.
func (*Report) Descriptor() ([]byte, []int) {
	return file_pkg_proto_appointment_proto_rawDescGZIP(), []int{1}
}

func (x *Report) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Report) GetAppointmentId() uint32 {
	if x != nil {
		return x.AppointmentId
	}
	return 0
}

func (x *Report) GetReferringHospital() string {
	if x != nil {
		return x.ReferringHospital
	}
	return ""
}

func (x *Report) GetFindings() string {
	if x != nil {
		return x.Findings
	}
	return ""
}

func (x *Report) GetImpressions() string {
	if x != nil {
		return x.Impressions
	}
	return ""
}

func (x *Report) GetClinicalDetails() string {
	if x != nil {
		return x.ClinicalDetails
	}
	return ""
}

func (x *Report) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *Report) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type RequestForm struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint32                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	AppointmentId uint32                 `protobuf:"varint,2,opt,name=appointment_id,json=appointmentId,proto3" json:"appointment_id,omitempty"`
	Url           string                 `protobuf:"bytes,3,opt,name=url,proto3" json:"url,omitempty"`
	IsFile        bool                   `protobuf:"varint,4,opt,name=is_file,json=isFile,proto3" json:"is_file,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RequestForm) Reset() {
	*x = RequestForm{}
	mi := &file_pkg_proto_appointment_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RequestForm) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestForm) ProtoMessage() {}

func (x *RequestForm) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_appointment_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestForm.ProtoReflect.Descriptor instead.
func (*RequestForm) Descriptor() ([]byte, []int) {
	return file_pkg_proto_appointment_proto_rawDescGZIP(), []int{2}
}

func (x *RequestForm) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *RequestForm) GetAppointmentId() uint32 {
	if x != nil {
		return x.AppointmentId
	}
	return 0
}

func (x *RequestForm) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *RequestForm) GetIsFile() bool {
	if x != nil {
		return x.IsFile
	}
	return false
}

var File_pkg_proto_appointment_proto protoreflect.FileDescriptor

const file_pkg_proto_appointment_proto_rawDesc = "" +
	"\n" +
	"\x1bpkg/proto/appointment.proto\x12\x05proto\x1a\x1fgoogle/protobuf/timestamp.proto\"\xea\x03\n" +
	"\vAppointment\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\rR\x02id\x12\x1d\n" +
	"\n" +
	"patient_id\x18\x02 \x01(\rR\tpatientId\x12!\n" +
	"\fis_emergency\x18\x03 \x01(\bR\visEmergency\x12!\n" +
	"\fprocedure_id\x18\x04 \x01(\rR\vprocedureId\x12\x1b\n" +
	"\treport_id\x18\x05 \x01(\rR\breportId\x12.\n" +
	"\x13referring_doctor_id\x18\x06 \x01(\rR\x11referringDoctorId\x12.\n" +
	"\x13reporting_doctor_id\x18\a \x01(\rR\x11reportingDoctorId\x12E\n" +
	"\x10appointment_time\x18\b \x01(\v2\x1a.google.protobuf.TimestampR\x0fappointmentTime\x12\x1d\n" +
	"\n" +
	"payment_id\x18\t \x01(\rR\tpaymentId\x12\x12\n" +
	"\x04note\x18\n" +
	" \x01(\tR\x04note\x12\x16\n" +
	"\x06status\x18\v \x01(\tR\x06status\x12.\n" +
	"\x13has_paid_commission\x18\f \x01(\bR\x11hasPaidCommission\x12'\n" +
	"\x0forganization_id\x18\r \x01(\rR\x0eorganizationId\"\x81\x02\n" +
	"\x06Report\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\rR\x02id\x12%\n" +
	"\x0eappointment_id\x18\x02 \x01(\rR\rappointmentId\x12-\n" +
	"\x12referring_hospital\x18\x03 \x01(\tR\x11referringHospital\x12\x1a\n" +
	"\bfindings\x18\x04 \x01(\tR\bfindings\x12 \n" +
	"\vimpressions\x18\x05 \x01(\tR\vimpressions\x12)\n" +
	"\x10clinical_details\x18\x06 \x01(\tR\x0fclinicalDetails\x12\x10\n" +
	"\x03url\x18\a \x01(\tR\x03url\x12\x16\n" +
	"\x06status\x18\b \x01(\tR\x06status\"o\n" +
	"\vRequestForm\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\rR\x02id\x12%\n" +
	"\x0eappointment_id\x18\x02 \x01(\rR\rappointmentId\x12\x10\n" +
	"\x03url\x18\x03 \x01(\tR\x03url\x12\x17\n" +
	"\ais_file\x18\x04 \x01(\bR\x06isFileB\x17Z\x15pkg/proto/appointmentb\x06proto3"

var (
	file_pkg_proto_appointment_proto_rawDescOnce sync.Once
	file_pkg_proto_appointment_proto_rawDescData []byte
)

func file_pkg_proto_appointment_proto_rawDescGZIP() []byte {
	file_pkg_proto_appointment_proto_rawDescOnce.Do(func() {
		file_pkg_proto_appointment_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_pkg_proto_appointment_proto_rawDesc), len(file_pkg_proto_appointment_proto_rawDesc)))
	})
	return file_pkg_proto_appointment_proto_rawDescData
}

var file_pkg_proto_appointment_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_pkg_proto_appointment_proto_goTypes = []any{
	(*Appointment)(nil),         // 0: proto.Appointment
	(*Report)(nil),              // 1: proto.Report
	(*RequestForm)(nil),         // 2: proto.RequestForm
	(*timestamp.Timestamp)(nil), // 3: google.protobuf.Timestamp
}
var file_pkg_proto_appointment_proto_depIdxs = []int32{
	3, // 0: proto.Appointment.appointment_time:type_name -> google.protobuf.Timestamp
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_pkg_proto_appointment_proto_init() }
func file_pkg_proto_appointment_proto_init() {
	if File_pkg_proto_appointment_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_pkg_proto_appointment_proto_rawDesc), len(file_pkg_proto_appointment_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_proto_appointment_proto_goTypes,
		DependencyIndexes: file_pkg_proto_appointment_proto_depIdxs,
		MessageInfos:      file_pkg_proto_appointment_proto_msgTypes,
	}.Build()
	File_pkg_proto_appointment_proto = out.File
	file_pkg_proto_appointment_proto_goTypes = nil
	file_pkg_proto_appointment_proto_depIdxs = nil
}
