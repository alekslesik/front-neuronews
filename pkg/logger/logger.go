package logger

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Level string

const (
	DEVELOPMENT Level = "development"
	PRODUCTION  Level = "production"
)

var (
	once            sync.Once
	ErrCreateFile   error = errors.New("create log file error")
	ErrLogFilePath  error = errors.New("log file path is not specified")
	ErrLevelMissing error = errors.New("logging level missing")
)

// Logger with zerolog logger instance and log file
type Logger struct {
	*zerolog.Logger
	File *os.File
}

// New create new logger instance with level. File string must be like "./path/logname.log"
func New(l Level, path string) (*Logger, error) {
	if l == "" {
		return nil, ErrLevelMissing
	}

	if path == "" {
		return nil, ErrLogFilePath
	}

	logger := new(Logger)

	// create log file
	err := logger.createLogFile(path)
	if err != nil {
		return nil, errors.Join(ErrCreateFile, err)
	}

	// new logger depends on log level
	switch l {
	case DEVELOPMENT:
		logger.setDevLogger()
	case PRODUCTION:
		logger.setProdLogger()
	}

	// wrong log level
	return logger, nil
}

// createLogFile create log file in specified filePath
func (l *Logger) createLogFile(path string) error {
	// Get dir where log file must be
	logDir := filepath.Dir(path)

	// Check existing dir, and create if not exists
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return err
		}
	}

	// Create or open log file for writing
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	l.File = f

	return nil
}

// setProdLogger return logger with logging in file only
func (l *Logger) setProdLogger() {
	zerolog.TimeFieldFormat = time.RFC1123
	z := zerolog.New(l.File).
		Level(zerolog.WarnLevel).
		With().
		Stack().
		Timestamp().
		Caller().
		Logger()

	l.Logger = &z
}

// setDevLogger return logger with logging in file and console
func (l *Logger) setDevLogger() {
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC1123}
	multi := zerolog.MultiLevelWriter(consoleWriter, l.File)

	z := zerolog.New(multi).
		Level(zerolog.DebugLevel).
		With().
		Stack().
		Timestamp().
		Caller().
		Logger()

	l.Logger = &z
}

// SetGlobalLog set global logger
func SetGlobalLog() {
	// set up once
	once.Do(setOnceGlobalLog)
}

// setOnceGlobalLog set once global logger in application
func setOnceGlobalLog() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    false,
		TimeFormat: "15:04:05",
		PartsOrder: []string{
			zerolog.TimestampFieldName,
			zerolog.LevelFieldName,
			zerolog.CallerFieldName,
			zerolog.MessageFieldName,
		},
	})
}

// TODO
// Потенциальные Улучшения
// Документация:
//     Добавьте комментарии к экспортируемым функциям и типам в соответствии с конвенциями Go.
//     Это улучшит понимание кода и его использование.
// Обработка Ошибок:
//     В функции createLogFile, если ошибка при создании директории не относится к os.IsNotExist,
//     она будет проигнорирована. Рекомендуется обрабатывать все возможные ошибки.
// Уровни Логирования:
//     Уровни логирования (DEVELOPMENT и PRODUCTION) могут быть расширены для большей гибкости.
//     Например, добавление уровней DEBUG, INFO, WARN, ERROR позволит более детально контролировать логирование.
// Глобальная Конфигурация:
//     Ваша функция SetGlobalLog устанавливает глобальный логгер с фиксированными настройками.
//     Хотя это и удобно, это также ограничивает гибкость. Рассмотрите возможность передачи параметров в SetGlobalLog,
//     чтобы позволить настройку глобального логгера в зависимости от потребностей приложения.

// Замечания
// Возврат ошибок:
//     Рассмотрите возможность добавления более специфичных ошибок вместо ErrCreateFile и ErrLevelMissing.
//     Это может помочь при отладке и обработке ошибок.

// Закрытие Файлов:
//     Важно убедиться, что файлы, открытые для логирования, корректно закрываются при завершении программы или при переключении логгера.
//     В текущей реализации файлы остаются открытыми.

// Конфигурация Пути к Файлу:
//     Путь к файлу лога задается напрямую в функции New. Рассмотрите использование конфигурационного файла
//     или переменных окружения для более гибкой конфигурации пути к файлу лога.

// Расширяемость:
//     Ваш пакет предоставляет заранее определенные конфигурации логгера.
//     Возможно, стоит предусмотреть механизм для более тонкой настройки логгера
//     (например, изменение формата вывода, выбор уровней логирования) во время выполнения.

// Использование Констант:
//     Рассмотрите использование встроенных в zerolog констант для уровней логирования вместо собственных
//     строковых представлений.
