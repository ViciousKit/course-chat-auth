package models

import "github.com/golang/protobuf/ptypes/timestamp"

type User struct {
	Id        int64               `json:"id"`
	Name      string              `json:"name"`
	Email     string              `json:"email"`
	Password  string              `json:"password"`
	Role      int                 `json:"role"`
	CreatedAt timestamp.Timestamp `json:"created_at"`
	UpdatedAt timestamp.Timestamp `json:"updated_at"`
}
