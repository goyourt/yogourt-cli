package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)

// Config Structure du fichier config
type Config struct {
	AppName  string   `yaml:"app_name"`
	Version  string   `yaml:"version"`
	Mode     string   `yaml:"mode"`
	EnvFiles EnvFiles `yaml:"env_files"`

	Server struct {
		Port int  `yaml:"port"`
		CORS bool `yaml:"cors"`
	} `yaml:"server"`

	Database struct {
		Type     string `yaml:"type"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`

	Paths struct {
		ModelFolder string `yaml:"model_folder"`
		ProjectName string `yaml:"project_name"`
		MainFile    string `yaml:"main_file"`
		RouteFolder string `yaml:"route_folder"`
		APIFolder   string `yaml:"api_folder"`
	} `yaml:"paths"`
}

type EnvFiles []string

func (e *EnvFiles) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var single string
	if err := unmarshal(&single); err == nil {
		single = strings.TrimSpace(single)
		if single == "" {
			*e = nil
			return nil
		}
		*e = []string{single}
		return nil
	}

	var multiple []string
	if err := unmarshal(&multiple); err != nil {
		return fmt.Errorf("env_files must be a string or a list of strings")
	}

	var envFiles []string
	for _, envFile := range multiple {
		envFile = strings.TrimSpace(envFile)
		if envFile != "" {
			envFiles = append(envFiles, envFile)
		}
	}
	*e = envFiles

	return nil
}

type envFileConfig struct {
	EnvFiles EnvFiles `yaml:"env_files"`
}

// LoadConfig Lecture et parse du fichier config
func LoadConfig(path string) (*Config, error) {

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("❌ Impossible de lire config.yaml : %v", err)
	}

	if err := loadEnvFiles(file); err != nil {
		return nil, err
	}

	replaced := os.ExpandEnv(string(file))

	var cfg Config
	err = yaml.Unmarshal([]byte(replaced), &cfg)
	if err != nil {
		return nil, fmt.Errorf("❌ Erreur de parsing YAML : %v", err)
	}

	return &cfg, nil
}

func loadEnvFiles(configContent []byte) error {
	cfg := &envFileConfig{}
	if err := yaml.Unmarshal(configContent, cfg); err != nil {
		return fmt.Errorf("❌ Erreur de parsing env_files : %v", err)
	}
	if len(cfg.EnvFiles) == 0 {
		return nil
	}
	if err := godotenv.Load(cfg.EnvFiles...); err != nil {
		return fmt.Errorf("❌ Impossible de charger env_files %v : %v", cfg.EnvFiles, err)
	}
	return nil
}
