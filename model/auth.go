package model

import (
	"time"

	"github.com/arf-dev/mekari-test/pkg/httputils/response"
)

type LoginRequest struct {
	Email string `json:"email"`
}

type LoginResponse struct {
	response.BaseResponse
	Data string `json:"data"`
}

type User struct {
	Id        string    `db:"id"`
	Email     string    `db:"email"`
	Name      string    `db:"name"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
}
