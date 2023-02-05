package model

type TokenPayload struct {
	Email       string
	Exp         int64
	RoleID      uint64
	Permissions []string
}
