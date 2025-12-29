package model

import (
	"time"

	"github.com/arf-dev/mekari-test/pkg/httputils/response"
)

type LoginRequest struct {
	Email string `json:"email" example:"test@gmail.com"`
}

type LoginResponse struct {
	response.BaseResponse
	Data string `json:"data"`
}

type User struct {
	Id        int32     `db:"id"`
	Email     string    `db:"email"`
	Name      string    `db:"name"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
}
