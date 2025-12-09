package config

import (
	"github.com/spf13/viper"
)

type App struct {
	AppPort string `json:"app_port"`
	AppEnv  string `json:"app_env"`

	JwtSecretKey string `json:"jwt_secret_key"`
	JwtIssuer    string `json:"jwt_issuer"`
}

type PsqlDB struct {
	Host      string `json:"host"`
	Port      string `json:"port"`
	User      string `json:"user"`
	Password  string `json:"password"`
	DBName    string `json:"db_name"`
	SSLMode   string `json:"ssl_mode"`
	DBMaxOpen int    `json:"db_max_open"`
	DBMaxIdle int    `json:"db_max_idle"`
}

type CloudFlare struct {
	Name      string `json:"name"`
	ApiKey    string `json:"api_key"`
	ApiSecret string `json:"api_secret"`
	Token     string `json:"token"`
	AccountId string `json:"account_id"`
	PublicUrl string `json:"public_url"`
}

type Cloudinary struct {
	Cloudname  string `json:"cloud_name"`
	Apikey     string `json:"api_key"`
	ApiSecret  string `json:"api_secret"`
	UploadFile string `json:"upload_file"`
}

type Config struct {
	App    App
	PsqlDB PsqlDB
	CF     CloudFlare
	CD     Cloudinary
}

func NewConfig() *Config {
	return &Config{
		App: App{
			AppPort: viper.GetString("APP_PORT"),
			AppEnv:  viper.GetString("APP_ENV"),

			JwtSecretKey: viper.GetString("JWT_SECRET_KEY"),
			JwtIssuer:    viper.GetString("JWT_ISSUER"),
		},

		PsqlDB: PsqlDB{
			Host:      viper.GetString("DATABASE_HOST"),
			Port:      viper.GetString("DATABASE_PORT"),
			User:      viper.GetString("DATABASE_USER"),
			Password:  viper.GetString("DATABASE_PASSWORD"),
			SSLMode:   viper.GetString("DATABASE_SSL_MODE"),
			DBName:    viper.GetString("DATABASE_NAME"),
			DBMaxOpen: viper.GetInt("DATABASE_MAX_CONNECTION"),
			DBMaxIdle: viper.GetInt("DATABASE_MAX_IDLE_CONNECTION"),
		},
		CF: CloudFlare{
			Name:      viper.GetString("CLOUDFLARE_BUCKET_NAME"),
			ApiKey:    viper.GetString("CLOUDFLARE_API_KEY"),
			ApiSecret: viper.GetString("CLOUDFLARE_API_SECRET"),
			Token:     viper.GetString("CLOUDFLARE_TOKEN"),
			PublicUrl: viper.GetString("CLOUDFLARE_PUBLIC_URL"),
			AccountId: viper.GetString("CLOUDFLARE_ACCOUNT_ID"),
		},
		CD: Cloudinary{
			Cloudname:  viper.GetString("CLOUDINARY_CLOUD_NAME"),
			Apikey:     viper.GetString("CLOUDINARY_API_KEY"),
			ApiSecret:  viper.GetString("CLOUDINARY_API_SECRET"),
			UploadFile: viper.GetString("CLOUDINARY_UPLOAD_FILE"),
		},
	}
}
