// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package sqlcore

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type RoleType string

const (
	RoleTypeUserAdmin  RoleType = "user_admin"
	RoleTypeUserGov    RoleType = "user_gov"
	RoleTypeUserVerify RoleType = "user_verify"
	RoleTypeUser       RoleType = "user"
)

func (e *RoleType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = RoleType(s)
	case string:
		*e = RoleType(s)
	default:
		return fmt.Errorf("unsupported scan type for RoleType: %T", src)
	}
	return nil
}

type NullRoleType struct {
	RoleType RoleType
	Valid    bool // Valid is true if RoleType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullRoleType) Scan(value interface{}) error {
	if value == nil {
		ns.RoleType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.RoleType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullRoleType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.RoleType), nil
}

type Account struct {
	ID        uuid.UUID
	Email     string
	Role      RoleType
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Session struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Token     string
	Client    string
	IpFirst   string
	IpLast    string
	CreatedAt time.Time
	LastUsed  time.Time
}
