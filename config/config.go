package config

import (
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type (
	Config struct {
		Server *Server
		Db     *Db
		Auth   *Auth
	}

	Server struct {
		Port int
	}

	Db struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
		TimeZone string
	}

	// Auth agrupa la configuración de JWT. En producción el secreto se
	// setea vía la variable de entorno AUTH_JWT_SECRET (Viper mapea
	// auth.jwt_secret <-> AUTH_JWT_SECRET por el EnvKeyReplacer de abajo).
	Auth struct {
		JWTSecret      string `mapstructure:"jwt_secret"`
		JWTExpiryHours int    `mapstructure:"jwt_expiry_hours"`
	}
)

var (
	once           sync.Once
	configInstance *Config
)

func GetConfig() *Config {
	// usamos viper para setear la config y usarla through nuestra app
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("../../")
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}

		if err := viper.Unmarshal(&configInstance); err != nil {
			panic(err)
		}
	})

	// devolvemos nuestra instancia de configuracion
	return configInstance
}
