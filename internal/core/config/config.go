package config

type Config struct {
	RestServer RestServer  `mapstructure:"restServer"`
	Database   MongoConfig `mapstructure:"database"`
}

type RestServer struct {
	Port int       `mapsturcture:"port"`
	Jwt  JwtConfig `mapstructure:"jwt"`
}

type JwtConfig struct {
	Secret   string `mapstructure:"secret"`
	ExpireIn int    `mapstructure:"expiresIn"`
	Issuer   string `mapstructure:"issuer"`
}

type MongoConfig struct {
	Host              string `mapstructure:"host"`
	Port              int    `mapstructure:"port"`
	DatabaseName      string `mapstructure:"databaseName"`
	MaxPoolSize       int    `mapstructure:"maxPoolSize"`
	ConnectionTimeout int    `mapstructure:"connectionTimeout"`
}
