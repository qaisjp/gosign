package main

type Config struct {
	LogLevel string `default:"debug"`

	CoSign CoSignConfig

	// This is how long it should take for a login token to time out
	Timeout int `default:"43200"`

	Address  string `default:"0.0.0.0:8080"`
	Insecure bool   `default:"false"`

	Tokens []Token
}

// CoSignConfig contains all configuration data for a CoSign connection
type CoSignConfig struct {
	DaemonHost string `required:"true"`
	DaemonPort string `default:"6663"`
	CGIAddress string `required:"true"`
	ServerName string `required:"true"`
	Service    string `required:"true"`

	Insecure bool   `default:"false"`
	CAFile   string `required:"true"`
	KeyFile  string `required:"true"`
	CertFile string `required:"true"`
	// User             string `required:"true"`
	// Database         string `required:"true"`
	// Password         string
}

// Token is an "API" user
type Token struct {
	Name string `required:"true"`
	Key  string `required:"true"`
}
