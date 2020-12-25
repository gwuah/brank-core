package core

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Environment string

const (
	Development Environment = "dev"
	Staging                 = "stage"
	Production              = "prod"
	Test                    = "test"
)

type Config struct {
	WORKER_POOL_SIZE          int
	PG_HOST                   string
	PG_PORT                   string
	PG_NAME                   string
	PG_USER                   string
	PG_PASS                   string
	PG_SSLMODE                string
	REDIS_ADDRESS             string
	REDIS_PASSWORD            string
	REDIS_DB                  int
	REDIS_URL                 string
	DATABASE_URL              string
	PORT                      int
	JWT_SIGNING_KEY           string
	KAFKA_GROUP_ID            string
	CLOUDKARAFKA_BROKERS      string
	CLOUDKARAFKA_USERNAME     string
	CLOUDKARAFKA_PASSWORD     string
	CLOUDKARAFKA_TOPIC_PREFIX string
	SSL_CA                    string
	RUN_SEEDS                 bool

	ENVIRONMENT Environment
}

func Get(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func GetInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		i, err := strconv.Atoi(v)
		if err != nil {
			log.Printf("%s: %s", key, err)
			return fallback
		}
		return i
	}
	return fallback
}

func GetEnvironment() Environment {
	if env := Get("ENV", ""); env == "" {
		return Development
	} else {
		return Environment(env)
	}
}

func load() {
	if env := GetEnvironment(); env == Production {
		return
	}
	err := godotenv.Load()
	if err != nil {
		log.Printf("error loading env file")
	}
}

func NewConfig() *Config {
	load()
	return &Config{
		WORKER_POOL_SIZE:          GetInt("WORKER_POOL_SIZE", 10),
		PG_HOST:                   Get("PG_HOST", "localhost"),
		PG_PORT:                   Get("PG_PORT", "5432"),
		PG_NAME:                   Get("PG_NAME", "oka_dev"),
		PG_USER:                   Get("PG_USER", "user"),
		PG_PASS:                   Get("PG_PASS", "password"),
		PG_SSLMODE:                Get("PG_SSLMODE", "disable"),
		REDIS_ADDRESS:             Get("REDIS_ADDRESS", "localhost:6379"),
		REDIS_PASSWORD:            Get("REDIS_PASSWORD", ""),
		REDIS_DB:                  GetInt("REDIS_DB", 1),
		REDIS_URL:                 Get("REDIS_URL", ""),
		DATABASE_URL:              Get("DATABASE_URL", ""),
		PORT:                      GetInt("PORT", 5454),
		JWT_SIGNING_KEY:           Get("JWT_SIGNING_KEY", "4RH93HFFUYFBY384V9F3Fbhbe3y34g36w7v37273v2tv32v73v2"),
		KAFKA_GROUP_ID:            Get("KAFKA_GROUP_ID", "brank_mq"),
		CLOUDKARAFKA_BROKERS:      Get("CLOUDKARAFKA_BROKERS", ""),
		CLOUDKARAFKA_USERNAME:     Get("CLOUDKARAFKA_USERNAME", ""),
		CLOUDKARAFKA_PASSWORD:     Get("CLOUDKARAFKA_PASSWORD", ""),
		CLOUDKARAFKA_TOPIC_PREFIX: Get("CLOUDKARAFKA_TOPIC_PREFIX", "6nnq3d64-"),
		SSL_CA: Get("SSL_CA", `-----BEGIN CERTIFICATE-----
		MIIF2DCCA8CgAwIBAgIQTKr5yttjb+Af907YWwOGnTANBgkqhkiG9w0BAQwFADCB
		hTELMAkGA1UEBhMCR0IxGzAZBgNVBAgTEkdyZWF0ZXIgTWFuY2hlc3RlcjEQMA4G
		A1UEBxMHU2FsZm9yZDEaMBgGA1UEChMRQ09NT0RPIENBIExpbWl0ZWQxKzApBgNV
		BAMTIkNPTU9ETyBSU0EgQ2VydGlmaWNhdGlvbiBBdXRob3JpdHkwHhcNMTAwMTE5
		MDAwMDAwWhcNMzgwMTE4MjM1OTU5WjCBhTELMAkGA1UEBhMCR0IxGzAZBgNVBAgT
		EkdyZWF0ZXIgTWFuY2hlc3RlcjEQMA4GA1UEBxMHU2FsZm9yZDEaMBgGA1UEChMR
		Q09NT0RPIENBIExpbWl0ZWQxKzApBgNVBAMTIkNPTU9ETyBSU0EgQ2VydGlmaWNh
		dGlvbiBBdXRob3JpdHkwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQCR
		6FSS0gpWsawNJN3Fz0RndJkrN6N9I3AAcbxT38T6KhKPS38QVr2fcHK3YX/JSw8X
		pz3jsARh7v8Rl8f0hj4K+j5c+ZPmNHrZFGvnnLOFoIJ6dq9xkNfs/Q36nGz637CC
		9BR++b7Epi9Pf5l/tfxnQ3K9DADWietrLNPtj5gcFKt+5eNu/Nio5JIk2kNrYrhV
		/erBvGy2i/MOjZrkm2xpmfh4SDBF1a3hDTxFYPwyllEnvGfDyi62a+pGx8cgoLEf
		Zd5ICLqkTqnyg0Y3hOvozIFIQ2dOciqbXL1MGyiKXCJ7tKuY2e7gUYPDCUZObT6Z
		+pUX2nwzV0E8jVHtC7ZcryxjGt9XyD+86V3Em69FmeKjWiS0uqlWPc9vqv9JWL7w
		qP/0uK3pN/u6uPQLOvnoQ0IeidiEyxPx2bvhiWC4jChWrBQdnArncevPDt09qZah
		SL0896+1DSJMwBGB7FY79tOi4lu3sgQiUpWAk2nojkxl8ZEDLXB0AuqLZxUpaVIC
		u9ffUGpVRr+goyhhf3DQw6KqLCGqR84onAZFdr+CGCe01a60y1Dma/RMhnEw6abf
		Fobg2P9A3fvQQoh/ozM6LlweQRGBY84YcWsr7KaKtzFcOmpH4MN5WdYgGq/yapiq
		crxXStJLnbsQ/LBMQeXtHT1eKJ2czL+zUdqnR+WEUwIDAQABo0IwQDAdBgNVHQ4E
		FgQUu69+Aj36pvE8hI6t7jiY7NkyMtQwDgYDVR0PAQH/BAQDAgEGMA8GA1UdEwEB
		/wQFMAMBAf8wDQYJKoZIhvcNAQEMBQADggIBAArx1UaEt65Ru2yyTUEUAJNMnMvl
		wFTPoCWOAvn9sKIN9SCYPBMtrFaisNZ+EZLpLrqeLppysb0ZRGxhNaKatBYSaVqM
		4dc+pBroLwP0rmEdEBsqpIt6xf4FpuHA1sj+nq6PK7o9mfjYcwlYRm6mnPTXJ9OV
		2jeDchzTc+CiR5kDOF3VSXkAKRzH7JsgHAckaVd4sjn8OoSgtZx8jb8uk2Intzna
		FxiuvTwJaP+EmzzV1gsD41eeFPfR60/IvYcjt7ZJQ3mFXLrrkguhxuhoqEwWsRqZ
		CuhTLJK7oQkYdQxlqHvLI7cawiiFwxv/0Cti76R7CZGYZ4wUAc1oBmpjIXUDgIiK
		boHGhfKppC3n9KUkEEeDys30jXlYsQab5xoq2Z0B15R97QNKyvDb6KkBPvVWmcke
		jkk9u+UJueBPSZI9FoJAzMxZxuY67RIuaTxslbH9qh17f4a+Hg4yRvv7E491f0yL
		S0Zj/gA0QHDBw7mh3aZw4gSzQbzpgJHqZJx64SIDqZxubw5lT2yHh17zbqD5daWb
		QOhTsiedSrnAdyGN/4fy3ryM7xfft0kL0fJuMAsaDk527RH89elWsn2/x20Kk4yl
		0MC2Hb46TpSi125sC8KKfPog88Tk5c0NqMuRkrF8hey1FGlmDoLnzc7ILaZRfyHB
		NVOFBkpdn627G190
		-----END CERTIFICATE-----`),
		ENVIRONMENT: GetEnvironment(),
		RUN_SEEDS:   true,
	}
}
