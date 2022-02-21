package user

type RoleID int

const (
	_ RoleID = iota
	RoleCashier
	RoleAdmin
)

type Role struct {
	RoleID   RoleID
	RoleName string
}

func GetRoleNameByRoleID(id RoleID) string {
	switch id {
	case RoleCashier:
		return "Kasir"
	case RoleAdmin:
		return "Admin"
	default:
		return "Invalid Role ID"
	}
}
