package models

type AuthData struct {
	InitData string   `json:"initData"`
	Role     RoleType `json:"role"`
}
