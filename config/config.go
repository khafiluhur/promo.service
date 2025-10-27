package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AppName     string `envconfig:"APP_NAME" default:"demo-api"`
	AppVersion  string `envconfig:"APP_VERSION" default:"v1.0.0"`
	AppApiPath  string `envconfig:"APP_API_PATH" default:"/api/v1/promo"`
	AppEnv      string `envconfig:"APP_ENV" default:"development"`
	AppPort     int    `envconfig:"APP_PORT" default:"1214"`
	AppDebug    bool   `envconfig:"APP_DEBUG" default:"true"`
	AppLocale   string `envconfig:"APP_LOCALE" default:"id"`
	AppEnvTopic string `envconfig:"APP_ENV_TOPIC" default:"staging"`

	RedisHost     string `envconfig:"REDIS_HOST" default:"localhost"`
	RedisPort     int    `envconfig:"REDIS_PORT" default:"6365"`
	RedisPassword string `envconfig:"REDIS_PASSWORD" default:""`
	RedisDbIndex  int    `envconfig:"REDIS_DBINDEX" default:"0"`

	DBMySQLHost              string `envconfig:"MYSQL_HOST" default:"localhost"`
	DBMySQLPort              int    `envconfig:"MYSQL_PORT" default:"3335"`
	DBMySQLDbUsername        string `envconfig:"MYSQL_USERNAME" default:"root"`
	DBMySQLPassword          string `envconfig:"MYSQL_PASSWORD" default:"password12345"`
	DBMySQLDbName            string `envconfig:"MYSQL_DBNAME" default:"promo"`
	DBMySQLMaxIdleConnection int    `envconfig:"MYSQL_MAX_IDLE_CONNECTION" default:"10"`
	DBMySQLMaxOpenConnection int    `envconfig:"MYSQL_MAX_OPEN_CONNECTION" default:"100"`
	DBMySQLConnMaxLifetime   int    `envconfig:"MYSQL_CONN_MAX_LIFETIME" default:"5"`

	ServiceTimeout int    `envconfig:"SERVICE_TIMEOUT" default:"30"`
	ExpiredService int    `envconfig:"EXPIRED_SERVICE" default:"3600"`
	UserService    string `envconfig:"USER_SERVICE" default:"https://sandbox-api.tripdeals.id"`
	PlatformConfig string `envconfig:"PROMO_CREDENTIAL" default:"W3sicGxhdGZvcm0iOiJ0cmlwZGVhbHMiLCJzeXN0ZW0iOiJ0cmlwZGVhbHMiLCJzZWNyZXQiOiJhMjgyMmIwNi02ODc0LTQ1ODktOWZmZS00MDdhMTdkOTI5MTEifSx7InBsYXRmb3JtIjoidHJpcGRlYWxzY21zIiwic3lzdGVtIjoidHJpcGRlYWxzIiwic2VjcmV0IjoiMTJkMmNjZTAtODQ1Yi00YTQyLThhYzEtNTlkZWRiNzM5ZGZmIn0seyJwbGF0Zm9ybSI6ImdvbGRlbnJhbWEiLCJzeXN0ZW0iOiJnb2xkZW5yYW1hIiwic2VjcmV0IjoiYTg5OGE4YjktNmQ1My00ODI5LWEzN2EtYmQzOTBmMmQxMzA5In1d"`
}

var data *Config

func Get() *Config {
	if data == nil {
		data = &Config{}
		envconfig.MustProcess("", data)
	}

	return data
}
