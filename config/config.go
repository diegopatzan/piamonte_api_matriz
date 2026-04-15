package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	SftpUser               string
	SftpPassword           string
	SftpHost               string
	SftpDestDir            string
	ApiPort                string
	SftpInsecureSkipVerify string
	SftpKnownHostsFile     string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, falling back to system environment variables.")
	}

	homeDir, _ := os.UserHomeDir()
	defaultKnownHosts := filepath.Join(homeDir, ".ssh", "known_hosts")

	return &Config{
		SftpUser:               getEnv("sftpUser", "TU_USUARIO_SFTP"),
		SftpPassword:           getEnv("sftpPassword", "TU_PASSWORD_SFTP"),
		SftpHost:               getEnv("sftpHost", "IP_PIAMONTE:22"),
		SftpDestDir:            getEnv("sftpDestDir", "/ruta/destino/en/piamonte/"),
		ApiPort:                getEnv("apiPort", ":8080"),
		SftpInsecureSkipVerify: getEnv("sftpInsecureSkipVerify", "false"),
		SftpKnownHostsFile:     getEnv("sftpKnownHostsFile", defaultKnownHosts),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
