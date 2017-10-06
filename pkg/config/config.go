package config

type Config struct {
	LogLevel string `default:"debug"`

	CoSign CoSignConfig

	// This is how long it should take for a login token to time out
	Timeout int `default:"43200"`

	Address  string `default:"0.0.0.0:8080"`
	Insecure bool   `default:"false"`
}

// CoSignConfig contains all configuration data for a CoSign connection
type CoSignConfig struct {
	DaemonAddress string `required:"true"`
	CGIAddress    string `required:"true"`
	Service       string `required:"true"`
	KeyFile       string `required:"true"`
	CertFile      string `required:"true"`
	Insecure      bool   `default:"false"`
	// User             string `required:"true"`
	// Database         string `required:"true"`
	// Password         string
}
