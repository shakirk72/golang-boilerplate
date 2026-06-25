package config

// Common holds values that are identical across all environments.
// Never read from .env — update here when releasing a new version.
var Common = struct {
	AppName        string
	AppVersion     string
	DefaultPort    string
	DefaultDBPort  string
	DefaultSSLMode string
	DefaultTZ      string
}{
	AppName:        "go-mvc-app",
	AppVersion:     "1.0.0",
	DefaultPort:    "8080",
	DefaultDBPort:  "5432",
	DefaultSSLMode: "disable",
	DefaultTZ:      "Asia/Kolkata",
}
