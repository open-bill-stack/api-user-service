package config

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	HttpPort int `mapstructure:"APP_HTTP_PORT"`
	GrpcPort int `mapstructure:"APP_GRPC_PORT"`
}

type DatabaseConfig struct {
	Url string `mapstructure:"DATABASE_URL"`
}
type AMQPConfig struct {
	Url string `mapstructure:"AMQP_URL"`
}

type JWTConfig struct {
	PrivateKeyBase64 string `mapstructure:"JWT_PRIVATE_KEY_BASE64"`
	PublicKeyBase64  string `mapstructure:"JWT_PUBLIC_KEY_BASE64"`

	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
}

type Config struct {
	App      *AppConfig
	Database *DatabaseConfig
	JWT      *JWTConfig
	AMQP     *AMQPConfig
}

func newAppConfig(cmd *cobra.Command) (*AppConfig, error) {
	var cfg AppConfig
	appViper := viper.New()

	for _, v := range os.Environ() {
		pair := strings.SplitN(v, "=", 2)
		appViper.SetDefault(pair[0], pair[1])
	}

	for key, value := range viper.AllSettings() {
		appViper.Set(key, value)
	}
	appViper.AutomaticEnv()
	appViper.SetConfigName(".env")
	appViper.SetConfigType("env")
	appViper.AddConfigPath(".")

	if value, err := strconv.Atoi(cmd.PersistentFlags().Lookup("http-port").Value.String()); err == nil {
		appViper.SetDefault("APP_HTTP_PORT", value)
	}
	if value, err := strconv.Atoi(cmd.PersistentFlags().Lookup("grpc-port").Value.String()); err == nil {
		appViper.SetDefault("APP_GRPC_PORT", value)
	}

	if err := appViper.ReadInConfig(); err != nil {
	}

	if err := appViper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("не вдалося розпарсити конфігурацію: %v", err)
	}
	return &cfg, nil
}

func newDatabaseConfig() (*DatabaseConfig, error) {
	var cfg DatabaseConfig
	databaseViper := viper.New()

	for _, v := range os.Environ() {
		pair := strings.SplitN(v, "=", 2)
		databaseViper.SetDefault(pair[0], pair[1])
	}

	for key, value := range viper.AllSettings() {
		databaseViper.Set(key, value)
	}
	databaseViper.AutomaticEnv()
	databaseViper.SetConfigName(".env")
	databaseViper.SetConfigType("env")
	databaseViper.AddConfigPath(".")

	if err := databaseViper.ReadInConfig(); err != nil {
	}

	if err := databaseViper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("не вдалося розпарсити конфігурацію: %v", err)
	}
	return &cfg, nil
}

func newAMQPConfig() (cfg *AMQPConfig, err error) {
	cfgViper := viper.New()

	for _, v := range os.Environ() {
		pair := strings.SplitN(v, "=", 2)
		cfgViper.SetDefault(pair[0], pair[1])
	}

	for key, value := range viper.AllSettings() {
		cfgViper.Set(key, value)
	}
	cfgViper.AutomaticEnv()
	cfgViper.SetConfigName(".env")
	cfgViper.SetConfigType("env")
	cfgViper.AddConfigPath(".")

	if err := cfgViper.ReadInConfig(); err != nil {
	}

	if err := cfgViper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("не вдалося розпарсити конфігурацію: %v", err)
	}
	return
}

func newJWTConfig() (cfg *JWTConfig, err error) {
	cfgViper := viper.New()

	for _, v := range os.Environ() {
		pair := strings.SplitN(v, "=", 2)
		cfgViper.SetDefault(pair[0], pair[1])
	}

	for key, value := range viper.AllSettings() {
		cfgViper.Set(key, value)
	}
	cfgViper.AutomaticEnv()
	cfgViper.SetConfigName(".env")
	cfgViper.SetConfigType("env")
	cfgViper.AddConfigPath(".")

	if err := cfgViper.ReadInConfig(); err != nil {
	}

	if err := cfgViper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("не вдалося розпарсити конфігурацію: %v", err)
	}

	if decoded, err := base64.StdEncoding.DecodeString(cfg.PrivateKeyBase64); err != nil {
		log.Fatal("Error decoding base64: ", err)
	} else {
		cfg.PrivateKey = ed25519.PrivateKey(decoded)
	}

	if decoded, err := base64.StdEncoding.DecodeString(cfg.PublicKeyBase64); err != nil {
		log.Fatal("Error decoding base64: ", err)
	} else {
		cfg.PublicKey = ed25519.PublicKey(decoded)
	}

	//cfg.PublicKey = publicKey
	//cfg.PrivateKey = privateKey
	return
}

func newService(cmd *cobra.Command) (*Config, error) {
	appConfig, err := newAppConfig(cmd)
	if err != nil {
		return nil, err
	}

	databaseConfig, err := newDatabaseConfig()
	if err != nil {
		return nil, err
	}

	amqpConfig, err := newAMQPConfig()
	if err != nil {
		return nil, err
	}

	jwtConfig, err := newJWTConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		App:      appConfig,
		Database: databaseConfig,
		JWT:      jwtConfig,
		AMQP:     amqpConfig,
	}, nil

}
