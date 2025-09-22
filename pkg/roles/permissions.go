package roles

// Permission definitions
type Permission string

// Convert Permission to string
func (p Permission) String() string {
	return string(p)
}

const (
	// Report permissions
	ViewReports   Permission = "view_reports"
	CreateReports Permission = "create_reports"
	EditReports   Permission = "edit_reports"
	VerifyReports Permission = "verify_reports"
	DeleteReports Permission = "delete_reports"

	// Appointment permissions
	ViewAppointments   Permission = "view_appointments"
	CreateAppointments Permission = "create_appointments"
	EditAppointments   Permission = "edit_appointments"
	DeleteAppointments Permission = "delete_appointments"

	// Financial permissions
	ViewFinancials  Permission = "view_financials"
	EditFinancials  Permission = "edit_financials"
	CreateInvoices  Permission = "create_invoices"
	ProcessPayments Permission = "process_payments"

	// User management permissions
	ViewUsers   Permission = "view_users"
	CreateUsers Permission = "create_users"
	EditUsers   Permission = "edit_users"
	DeleteUsers Permission = "delete_users"

	// Admin permissions
	ManageSystem  Permission = "manage_system"
	ViewAuditLogs Permission = "view_audit_logs"
	ManageRoles   Permission = "manage_roles"
	BackupRestore Permission = "backup_restore"
)

var DefaultRolePermissions = map[string][]Permission{
	UserTypes.FrontDesk: {
		ViewAppointments,
		CreateAppointments,
		EditAppointments,
		ViewReports,
	},
	UserTypes.Cashier: {
		ViewAppointments,
		ViewFinancials,
		CreateInvoices,
		ProcessPayments,
	},
	UserTypes.Radiographer: {
		ViewAppointments,
		ViewReports,
		CreateReports,
		EditReports,
	},
	UserTypes.ReportingAssistant: {
		ViewReports,
		CreateReports,
		EditReports,
	},
	UserTypes.Radiologist: {
		ViewReports,
		CreateReports,
		EditReports,
		VerifyReports,
		ViewAppointments,
	},
	UserTypes.Accountant: {
		ViewFinancials,
		EditFinancials,
		CreateInvoices,
		ProcessPayments,
		ViewAuditLogs,
	},
	UserTypes.Admin: {
		ViewUsers,
		CreateUsers,
		EditUsers,
		DeleteUsers,
		ManageSystem,
		ViewAuditLogs,
		ManageRoles,
		BackupRestore,
		ViewReports,
		ViewAppointments,
		ViewFinancials,
	},
}
