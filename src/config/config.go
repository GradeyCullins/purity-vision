package config

var (
	// Default port to expose the API server on.
	DefaultPort int = 8080
)

// DB defines connection string details for the postgres connection.
var DB = struct {
	Name     string
	User     string
	Password string
	SSLMode  string
}{
	"purity",
	"postgres",
	"password",
	"disable",
}
