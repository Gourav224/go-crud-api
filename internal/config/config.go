package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string `yaml:"address" env:"HTTP_SERVER_ADDR" env-default:":8080"`
}

type Config struct {
	Env         string     `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string     `yaml:"storage_path" env:"STORAGE_PATH" env-required:"true"`
	HTTPServer  HTTPServer `yaml:"http_server"`
}

// MustLoad reads configuration from file or environment variables
// and panics if something goes wrong.
func MustLoad() *Config {
	var configPath string

	// 1️⃣ Priority 1: ENV variable
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		configPath = path
	} else {
		// 2️⃣ Priority 2: Command-line flag
		flag.StringVar(&configPath, "config", "", "path to config file")
		flag.Parse()

		if configPath == "" {
			log.Fatal("Config path is not set (use --config or CONFIG_PATH env var)")
		}
	}

	// 3️⃣ Check existence
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configPath)
	}

	// 4️⃣ Parse YAML into struct
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config file: %v", err)
	}

	log.Printf("✅ Config loaded from %s", configPath)
	return &cfg
}
