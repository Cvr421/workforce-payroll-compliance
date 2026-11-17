package auth

var RolePermissions = map[string][]string{
	"admin": {
		"org.read", "org.write",
		"employees.read", "employees.write",
		"attendance.read", "attendance.write", "attendance.lock",
		"payroll.read", "payroll.create", "payroll.approve", "payroll.finalize",
		"compliance.read", "compliance.write", "compliance.review",
		"audit.read", "audit.write", "audit.verify",
		"users.manage",
		"reports.read",
	},
	"hr_admin": {
		"org.read", "org.write",
		"employees.read", "employees.write",
		"attendance.read", "attendance.write", "attendance.lock",
		"reports.read",
	},
	"payroll_manager": {
		"org.read",
		"employees.read",
		"attendance.read",
		"payroll.read", "payroll.create", "payroll.approve", "payroll.finalize",
		"reports.read",
	},
	"compliance_officer": {
		"org.read",
		"employees.read",
		"compliance.read", "compliance.write", "compliance.review",
		"audit.read",
		"reports.read",
	},
	"auditor": {
		"org.read",
		"employees.read",
		"compliance.read",
		"audit.read", "audit.write", "audit.verify",
		"reports.read",
	},
}

func GetPermissionsForRole(role string) []string {
	if perms, ok := RolePermissions[role]; ok {
		return perms
	}
	return nil
}

func HasPermission(role, permission string) bool {
	perms := GetPermissionsForRole(role)
	for _, p := range perms {
		if p == permission {
			return true
		}
	}
	return false
}

const (
	PermOrgRead        = "org.read"
	PermOrgWrite       = "org.write"
	PermEmployeesRead  = "employees.read"
	PermEmployeesWrite = "employees.write"
	PermAttendanceRead = "attendance.read"
	PermAttendanceWrite = "attendance.write"
	PermAttendanceLock = "attendance.lock"
	PermPayrollRead    = "payroll.read"
	PermPayrollCreate  = "payroll.create"
	PermPayrollApprove = "payroll.approve"
	PermPayrollFinalize = "payroll.finalize"
	PermComplianceRead = "compliance.read"
	PermComplianceWrite = "compliance.write"
	PermComplianceReview = "compliance.review"
	PermAuditRead      = "audit.read"
	PermAuditWrite     = "audit.write"
	PermAuditVerify    = "audit.verify"
	PermUsersManage    = "users.manage"
	PermReportsRead    = "reports.read"
)