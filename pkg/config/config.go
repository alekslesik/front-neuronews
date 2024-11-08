package config

import (
	"errors"
	"io/fs"
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

var (
	ErrEnvNotExists error = errors.New(".env file is not exists")
)

// AppConfig is general config for application
type AppConfig struct {
	Port      int    `env:"PORT" env-default:"443"`
	Host      string `env:"HOST" env-default:""`
	Env       string `env:"ENV" env-default:"development"`
	Delta     int    `env:"DELTA" env-default:"15"`
	AdminUser struct {
		Email    string `env:"ADMIN_EMAIL" env-default:"admin"`
		Password string `env:"ADMIN_PWD" env-default:"admin"`
	}
}

// KandinskyConfig is config for Kandinsky API
type KandinskyConfig struct {
	Key     string `env:"KAND_API_KEY"`
	Secret  string `env:"KAND_API_SECRET"`
	URL     string `env:"KAND_API_URL"`
	AuthURL string `env:"KAND_API_AUTH_URL"`
	GenURL  string `env:"KAND_API_GEN_URL"`
}

// LoggerConfig is config for logging part
type LoggerConfig struct {
	LogLevel    string `env:"LOG_LEVEL" env-default:"development"`
	LogFilePath string `env:"LOG_FILE" env-default:"./tmp/log.log"`
	MaxSize     int    `env:"LOG_MAXSIZE" env-default:"100"`
	MaxBackups  int    `env:"LOG_MAXBACKUP" env-default:"3"`
	MaxAge      int    `env:"LOG_MAXAGE" env-default:"24"`
	Compress    bool   `env:"LOG_COMPRESS" env-default:"true"`
}

// MySQLConfig is config fo MySQL part
type MySQLConfig struct {
	Driver string `env:"MYSQL_DRIVER" env-default:"mysql"`
	DSN    string `env:"MYSQL_DSN" env-default:"root:486464@tcp(localhost:3306)/neuronews?parseTime=true"`
}

// TLSConfig is config for TLS part
type TLSConfig struct {
	KeyPath  string `env:"TLS_KEY_PATH" env-default:"./tls/key.pem"`
	CertPath string `env:"TLS_CERT_PATH" env-default:"./tls/cert.pem"`
}

// SessionConfig is config for session part
type SessionConfig struct {
	Secret string `env:"SESSION_SECRET" env-default:"s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge"`
}

// SMTPConfig is config for SMTP part
type SMTPConfig struct {
	Host     string `env:"SMTP_HOST"`
	Port     int    `env:"SMTP_PORT"`
	Username string `env:"SMTP_USERNAME"`
	Password string `env:"SMTP_PASSWORD"`
	Sender   string `env:"SMTP_SENDER"`
}

// Config fo application
type Config struct {
	App     AppConfig
	Kand    KandinskyConfig
	Logger  LoggerConfig
	MySQL   MySQLConfig
	Session SessionConfig
	SMTP    SMTPConfig
	TLS     TLSConfig
}

var (
	instance *Config
	once     sync.Once
)

// New return instance of config (singleton)
func New() (*Config, error) {
	var err error

	once.Do(func() {
		instance, err = loadEnv()
	})

	return instance, err
}

// loadEnv load environments from .env file
func loadEnv() (*Config, error) {
	// load environments from .env file
	err := godotenv.Load("./configs/.env")

	if err != nil {
		//  if .env is not exists in ./configs/? check root path
		err := godotenv.Load()

		// if .env file not exists continue code and take default environments
		if err != nil {
			var pathErr *fs.PathError

			if !errors.As(err, &pathErr) {
				return nil, err
			} else {
				log.Println(err)
				log.Println(ErrEnvNotExists)
			}
		}
	}

	// read environments in structure
	cfg := &Config{}
	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

//TODO Планы по улучшению пакета конфигурации:
//
// 1. Валидация:
//    - Реализовать валидацию конфигурации для обеспечения корректности строк подключения, портов, форматов электронной почты и т.д.
//
// 2. Горячая Перезагрузка:
//    - Добавить поддержку горячей перезагрузки конфигураций без необходимости перезапуска приложения, например, при изменениях в файле `.env`.
//
// 3. Поддержка Множественных Форматов:
//    - Включить поддержку различных форматов файлов конфигурации, таких как YAML, JSON, TOML, для большей гибкости.
//
// 4. Интеграция с Внешними Системами:
//    - Интегрировать с внешними системами управления конфигурациями, такими как Consul, etcd, для динамического управления конфигурациями в распределенных системах.
//
// 5. Улучшенное Логирование:
//    - Улучшить логирование ошибок конфигурации для облегчения отладки и устранения неполадок.
//
// 6. Использование Контекста:
//    - Рассмотреть возможность использования `context.Context` для управления длительными или отменяемыми операциями во время загрузки или обновления конфигураций.
//
// 7. Конфигурации По Умолчанию:
//    - Включить возможность указания полного набора значений конфигурации по умолчанию непосредственно в коде, что полезно при отсутствии файла `.env`.
//
// 8. Обработка Ошибок:
//    - Пересмотреть обработку ошибок, чтобы обеспечить более гибкую логику в случаях, когда файл `.env` отсутствует или поврежден.
