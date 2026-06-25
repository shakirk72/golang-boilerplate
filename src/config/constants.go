package config

const (
	// Environments
	EnvDev  = "dev"
	EnvQA   = "qa"
	EnvProd = "prod"

	// HTTP status messages
	MsgSuccess            = "success"
	MsgCreated            = "created"
	MsgUpdated            = "updated"
	MsgDeleted            = "deleted"
	MsgNotFound           = "resource not found"
	MsgBadRequest         = "bad request"
	MsgUnauthorized       = "unauthorized"
	MsgForbidden          = "forbidden"
	MsgInternal           = "internal server error"
	MsgConflict           = "resource already exists"
	MsgValidationFailed   = "validation failed"
	MsgTooManyRequests    = "too many requests"
	MsgRequestTimeout     = "request timeout"
	MsgInvalidCredentials = "invalid email or password"
	MsgTokenInvalid       = "invalid or expired token"
	MsgLoginSuccess       = "login successful"

	// Roles
	RoleAdmin = "admin"
	RoleUser  = "user"

	// Pagination defaults
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 100

	// Rate limiting (requests per minute per IP)
	DefaultRateLimitPerMin = 60

	// Context keys
	CtxUserID    = "userID"
	CtxUserRole  = "userRole"
	CtxRequestID = "requestID"

	// Header keys
	HeaderRequestID     = "X-Request-ID"
	HeaderAuthorization = "Authorization"
	HeaderContentType   = "Content-Type"

	// API versions
	APIV1 = "/api/v1"
)
