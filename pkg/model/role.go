package model

import "errors"

const (
	RoleUser  string = "ROLE_USER"
	RoleAdmin string = "ROLE_ADMIN"
)

var (
	NoRoleError = errors.New("no role")
)
