package export

var mapping = map[string][]string{
	"patients": {
		"patient_id", "first_name", "last_name", "gender", "dob", "age", "phone_number", "email", "address", "created_at",
	},
	"appointments": {
		"appointment_number", "appointment_time", "patient_id", "patient_name", "patient_phone", "procedure", "modality", "center", "status", "is_emergency", "referring_doctor", "reporting_doctor", "accession_number", "total_amount", "amount_paid", "payment_method", "note", "created_at",
	},
	"receipts": {
		"appointment_number", "appointment_time", "patient_id", "patient_name", "patient_phone", "procedure", "modality", "center", "status", "is_emergency", "referring_doctor", "reporting_doctor", "total_amount", "amount_paid", "balance",
	},
	"paid-appointments": {
		"appointment_number", "appointment_time", "patient_id", "patient_name", "patient_phone", "procedure", "modality", "center", "status", "is_emergency", "referring_doctor", "reporting_doctor", "receipt_number", "total_amount", "amount_paid", "discount_amount", "discount_authorizer", "payment_method", "paid_at",
	},
	"reports": {
		"appointment_number", "patient_id", "patient_name", "procedure", "modality", "status", "report_type", "version_number", "radiologist", "referring_doctor", "referring_hospital", "is_emergency", "created_at", "confirmed_at", "impressions",
	},
	"referring-centers": {
		"name", "phone_number", "address", "region", "doctor_count", "has_organization", "organization_name",
	},
	"organizations": {
		"name", "domain", "phone_number", "address", "is_active", "staff_count", "appointment_count", "created_at",
	},
	"organization-staff": {
		"first_name", "last_name", "email", "role", "is_disabled", "last_login", "created_at",
	},
	"procedures": {
		"name", "modality", "price", "flat_commision",
	},
	"modalities": {
		"name", "short_name", "estimated_reporting_time_minutes", "procedure_count",
	},
	"centers": {
		"name", "address", "phone_number", "email", "region", "is_main", "modality_count",
	},
	"users": {
		"first_name", "last_name", "email", "phone_number", "role", "title_on_report", "center", "is_admin", "can_verify", "is_disabled", "last_login", "created_at",
	},
	"templates": {
		"template_name", "report_type", "clinical_detail", "findings", "impression",
	},
	"insurances": {
		"name", "telephone", "address", "payment_count", "created_at",
	},
	"discounts": {
		"id", "patient_id", "patient_name", "procedure", "total_amount", "amount", "status", "created_by_name", "created_at", "approved_at", "rejection_reason",
	},
	"adjustments": {
		"id", "appointment_id", "patient_id", "patient_name", "adjustment_type", "previous_amount", "new_amount", "amount", "status", "reason", "created_at", "processed_at",
	},
	"payment-providers": {
		"name", "type", "direction", "is_active", "created_at",
	},
	"campaigns": {
		"name", "medium", "target_audience", "status", "recipient_count", "scheduled_date", "created_at",
	},
	"announcements": {
		"title", "message", "is_published", "created_at",
	},
	"transactions": {
		"date", "patient_id", "patient_name", "modality", "procedure", "payment_method", "cost",
	},
	"invoices": {
		"invoice_number", "date", "organization", "patient_name", "procedure", "status", "total_amount",
	},
}
