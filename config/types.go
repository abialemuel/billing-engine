package config

var Service = "billing-engine"
var Version = "v1.0.0"
var GitCommit string
var OSBuildName string
var BuildDate string

type MainConfig struct {
	Log struct {
		Level  string `yaml:"level" validate:"oneof=trace debug info warn error fatal panic"`
		Format string `yaml:"format" validate:"oneof=text json"`
	} `yaml:"log"`
	APM struct {
		Enabled     bool     `yaml:"enabled"`
		Host        string   `yaml:"host"`
		Port        int      `yaml:"port" validate:"required,min=1,max=65535"`
		Rate        *float64 `yaml:"rate" validate:"omitempty,min=0.1,max=1"`
		DDAgentHost string   `yaml:"ddAgentHost" env:"DD_AGENT_HOST"`
		DDAgentPort int      `yaml:"ddAgentPort"`
	} `yaml:"apm"`
	App struct {
		Name    string `yaml:"name" validate:"required"`
		Port    int    `yaml:"port" validate:"required,min=1,max=65535"`
		Version string `yaml:"version" validate:"required"`
		Env     string `yaml:"env" validate:"required"`
	} `yaml:"app"`
	Postgres struct {
		Host     string `yaml:"host" validate:"required"`
		Port     int    `yaml:"port" validate:"required"`
		User     string `yaml:"user" validate:"required"`
		Password string `yaml:"password" validate:"required"`
		DB       string `yaml:"db" validate:"required"`
	} `yaml:"postgres"`
	Billing struct {
		DelinquentThreshold int `yaml:"delinquent_threshold" validate:"required,min=1"`
	} `yaml:"billing"`
}
