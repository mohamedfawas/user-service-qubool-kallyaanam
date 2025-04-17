package constants

// Role constants
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// Header constants for propagating user identity
const (
	HeaderUserID   = "X-User-ID"
	HeaderUserRole = "X-User-Role"
	HeaderUsername = "X-Username"
)
