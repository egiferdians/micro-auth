package transport

import (
	"github.com/egiferdians/micro-auth/auth"
)
type (
	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	LoginResponse struct {
		Status string `json:"status"`
		Err    error  `json:"errcode"`
		Data *auth.Authenticated `json:"data"`
	}
)