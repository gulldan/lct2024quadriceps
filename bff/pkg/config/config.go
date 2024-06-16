package config

import (
	"fmt"
	"io/fs"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog"
)

type Config struct {
	LogLevel string `yaml:"LOG_LEVEL" env:"LOG_LEVEL" env-default:"info"`

	Grpc          GrpcConfig `yaml:"http"`
	Minio         MinioConfig
	Postgres      PostgresConfig
	HTTPPort      string `env:"HTTP_PORT" env-default:"8888"`
	Wav2VecAddr   string `env:"WAV2VEC_ADDR" env-default:"http://10.35.56.4:18010"`
	VideocopyAddr string `env:"VIDEOCOPY_ADDR" env-default:"http://10.35.56.10:32323"`
}

type GrpcConfig struct {
	Address string `yaml:"address" env:"HTTP_ADDRESS" env-default:":7083"`
}

type PostgresConfig struct {
	Addr string `yaml:"pg_addr" env:"PG_ADDR"`
}

type MinioConfig struct {
	Endpoint          string `yaml:"minio_addr" env:"MINIO_ADDR"`
	AccessKey         string `yaml:"minio_access_key" env:"MINIO_ACCESS_KEY"`
	SecretAccessKey   string `yaml:"secret_access_key" env:"MINIO_SECRET_ACCESS_KEY"`
	IsUseSsl          bool   `yaml:"is_use_ssl" env:"MINIO_IS_USE_SSL"`
	VideoBucket       string `yaml:"video_bucket" env:"VIDEO_BUCKET" env-default:"video"`
	PreviewBucket     string `yaml:"preview_bucket" env:"PREVIEW_BUCKET" env-default:"preview"`
	OriginVideoBucket string `yaml:"orig_video_bucket" env:"ORIG_VIDEO_BUCKET" env-default:"origvideo"`
}

func InitConfig() (*Config, *zerolog.Level, error) {
	cnf := Config{}

	err := cleanenv.ReadConfig("config.yml", &cnf)
	if err != nil {
		_, ok := err.(*fs.PathError)
		if ok {
			err = cleanenv.ReadEnv(&cnf)
		}

		if err != nil {
			return nil, nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	if cnf.LogLevel == "" {
		cnf.LogLevel = "info"
	}

	logLevel, err := zerolog.ParseLevel(cnf.LogLevel)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid log level: %w", err)
	}

	return &cnf, &logLevel, nil
}
