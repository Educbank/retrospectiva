package handlers

import (
	"educ-retro/internal/auth"
)

// Global middleware instance
var authMiddleware = auth.AuthMiddleware()
