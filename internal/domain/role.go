package domain

const (
	RoleAdmin    = "ADMIN"
	RoleCustomer = "CUSTOMER"
)

type Role struct {
	ID       uint
	RoleName string
}
