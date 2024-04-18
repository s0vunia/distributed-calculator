package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env                      string                    `yaml:"env" env-default:"local"`
	UrlRabbit                string                    `yaml:"url_rabbit" env-required:"true"`
	Queue                    QueueConfig               `yaml:"queue"`
	CalculationTimeouts      CalculationTimeoutsConfig `yaml:"calculation_timeouts"`
	GRPC                     GRPCConfig                `yaml:"grpc"`
	HTTP                     HTTPConfig                `yaml:"http"`
	Postgres                 PostgresConfig            `yaml:"postgres"`
	TokenTTL                 time.Duration             `yaml:"token_ttl" env-default:"1h"`
	RetrySubExpressionTimout time.Duration             `yaml:"retry_sub_expression_timout" env-default:"40s"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type HTTPConfig struct {
	Port int `yaml:"port"`
}
type QueueConfig struct {
	NameQueueWithTasks         string `yaml:"name_queue_with_tasks"`
	NameQueueWithFinishedTasks string `yaml:"name_queue_with_finished_tasks"`
	NameQueueWithHeartbeats    string `yaml:"name_queue_with_heartbeats"`
	NameQueueWithRPC           string `yaml:"name_queue_with_rpc"`
}

type CalculationTimeoutsConfig struct {
	TimeCalculatePlus   time.Duration `yaml:"time_calculate_plus"`
	TimeCalculateMinus  time.Duration `yaml:"time_calculate_minus"`
	TimeCalculateMult   time.Duration `yaml:"time_calculate_mult"`
	TimeCalculateDivide time.Duration `yaml:"time_calculate_divide"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DbName   string `yaml:"db_name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
