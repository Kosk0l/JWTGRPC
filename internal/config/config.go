package config

import (
	"flag"
	"os"
	"time"
	"github.com/ilyakaznacheev/cleanenv"
)

//===================================================================================================================//

// Струкртура для YAML //local.yaml
type Config struct {
	Env 		string 		`yaml:"env" env-default:"local"`
	StoragePath string 		`yaml:"storage_path" env-required:"true"`
	TokenTTL time.Duration 	`yaml:"token_ttl" env-required:"true"`
	GRPC 		GRPCConfig 	`yaml:"grpc"`
}

type GRPCConfig struct {
	Port 	int 			`yaml:"port"`
	TimeOut time.Duration 	`yaml:"timeout"`
}

//===================================================================================================================//

// Функция загрузки конфига
func MustLoad() *Config {
	// Must - Функция не будет возвращать ошибку, если такая произошла 
	// АнтиПаттерн бизнесЛогики
	
	path := fetchConfigPath()
	if path == "" {
		panic("Config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist" + path)
	}

	var config Config
	if err := cleanenv.ReadConfig(path, &config); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &config
}

// Извлекает путь конфигурации из Флага Командной строки или переменной среды.
// Приоритет: Flag > env > default.
// Значение по умолчанию - Пустая строка.
func fetchConfigPath() string {
	var res string

	// BeckEnd --config="path/to/config.yaml" // Запуск // "config=" - имя флага
	// Переменная Флага // Имя флага // значение // Подсказака командной строки
	flag.StringVar(&res, "config", "", "path to config file") // Привязывает флаг к существующей переменной:
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}

// TODO: Переменная Окружения
// CONFIG_PATH=./path/to/config/file.yaml BeckEnd