package config

import (
	"flag"
	"fmt"
	"github.com/SETTER2000/prove/pkg/log/logger"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"sync"
)

type (
	Config struct {
		App     `json:"app" yaml:"app"`
		Storage `json:"storage" yaml:"storage"`
		Cookie  `json:"cookie" yaml:"cookie"`
		JWT     `json:"jwt" yml:"jwt"`
		GRPC    `json:"grpc" yml:"grpc"`
		Log     `json:"log" yml:"log"`
		HTTP    `json:"http" yaml:"http"`
	}
	App struct {
		Name           string `env-required:"true" json:"name"  yaml:"name"    env:"APP_NAME"`
		ConfigFileName string `env:"CONFIG"`
		Author         string `json:"author" yml:"author" env:"AUTHOR" envDefault:"Soviet Union"`
	}
	HTTP struct {
		TrustedSubnet        string `json:"trusted_subnet" yml:"trusted_subnet"  env:"TRUSTED_SUBNET"`
		BaseURL              string `json:"base_url" yml:"base_url" env:"BASE_URL"`
		ServerAddress        string `json:"server_address" yml:"server_address" env:"RUN_ADDRESS"`
		ServerDomain         string `env-required:"true" json:"server_domain" yml:"server_domain" env:"SERVER_DOMAIN"`
		CertsDir             string `json:"certs_dir" yml:"certs_dir" env:"CERTS_DIR"`
		EnableHTTPS          bool   `env:"ENABLE_HTTPS"`
		ResolveIPUsingHeader bool   `env:"RESOLVE_IP_USING_HEADER"`
	}
	Storage struct {
		// FILE_STORAGE_PATH путь до файла с сокращёнными URL (директории не создаёт)
		FileStorage string `json:"file_storage" yml:"file_storage" env:"FILE_STORAGE_PATH"`
		// Строка с адресом подключения к БД, например для PostgreSQL (драйвер pgx): postgres://username:password@localhost:5432/database_name
		ConnectDB string `json:"connect_db" yml:"connect_db" env:"DATABASE_URI"`
	}
	Cookie struct {
		AccessTokenName string `env-required:"true" yaml:"access_token_name" env:"ACCESS_TOKEN_NAME" envDefault:"access_token"`
		SecretKey       string `env-required:"true" yaml:"secret_key" env:"SECRET_KEY" envDefault:"RtsynerpoGIYdab_s234r"` // cookie encrypt application
	}
	JWT struct {
		// SECRET_KEY ключ шифрования для JWT токена авторизации
		Secret   string `env-required:"true" json:"secret" yml:"secret" env:"SECRET_JWT"`
		HashSalt string `env-required:"true" json:"hash_salt" yml:"hash_salt" env:"HASH_SALT_JWT"`
	}
	Log struct {
		// LOG_LEVEL переменная окружения, содержит значение уровня логирования проекта
		Level string `env-required:"true" json:"log_level"  yaml:"log_level"  env:"LOG_LEVEL"`
	}
	GRPC struct {
		Port string `env-required:"true" json:"port" yaml:"port" env:"GRPC_PORT"`
		Host string `env-required:"true" json:"host" yaml:"host" env:"GRPC_HOST"`
	}
)

// Config .
var cfg *Config
var once sync.Once

// GetConfig (singleton) возвращает инициализированную структуру конфига.
func GetConfig() *Config {
	once.Do(func() {
		l := logger.GetLogger()
		l.Info("read application configuration")
		cfg = &Config{}
		flag.StringVar(&cfg.ConfigFileName, "c", "", "configuration file name")
		flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "the base address of the resulting shortened URL")
		flag.StringVar(&cfg.ConnectDB, "d", "", "dsn connect string urlExample PostgreSQL: postgres://prove:DBprove-2023@127.0.0.1:5432/prove?sslmode=disable")
		flag.StringVar(&cfg.ServerAddress, "a", "", "host to listen on")
		flag.BoolVar(&cfg.EnableHTTPS, "s", false, "start server with https protocol")
		flag.StringVar(&cfg.ServerDomain, "dm", "", "server domain name")
		flag.StringVar(&cfg.TrustedSubnet, "t", "", "you can pass Classless Addressing Representation (CIDR) strings")
		flag.BoolVar(&cfg.ResolveIPUsingHeader, "resolve_ip_using_header", false, "Разрешение на проверку заголовка X-Real-IP и X-Forwarded-For")
		flag.StringVar(&cfg.CertsDir, "cd", "", "certificate directory")
		flag.StringVar(&cfg.FileStorage, "f", "", "path to file with abbreviated URLs")
		flag.StringVar(&cfg.SecretKey, "sk", "RtsynerpoGIYdab_s234r", "cookie secret key")
		flag.StringVar(&cfg.AccessTokenName, "at", "access_token", "Access Token Name")
		flag.StringVar(&cfg.Secret, "skj", "", "jwt secret key")
		flag.StringVar(&cfg.HashSalt, "hs", "", "salt for hash")
		flag.StringVar(&cfg.Port, "g", "", "server grpc port")
		flag.StringVar(&cfg.Host, "h", "", "server grpc host")
		flag.StringVar(&cfg.Author, "author", "", "author project")
		flag.BoolVar(&cfg.EnableHTTPS, "ps", false, "start server with https protocol")
		flag.Usage = func() {
			fmt.Fprintf(flag.CommandLine.Output(), "Prove Version %s %v\nUsage : Project Prove - URL Shortener Server\n", os.Args[0], cfg.Version)
			flag.PrintDefaults()
		}

		if err := cleanenv.ReadConfig("./config/config.json", cfg); err != nil {
			l = logger.GetLogger()
			help, _ := cleanenv.GetDescription(cfg, nil)
			l.Info(help)
			l.Fatal(err)
		}

		// Parse flags
		flag.Parse()

		// Parse environ
		err := env.Parse(cfg) // caarlos0
		if err != nil {
			panic(err)
		}
	})

	return cfg
}
