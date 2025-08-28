package config

import (
	"log"
	"os"
)

type Config struct {
	Address   string
	PFCPAddr  string
}

func Load() Config {
	cfg := Config{
		Address:  getenv("MB_SMF_ADDR", ":8080"),
		PFCPAddr: getenv("MB_SMF_PFCP_ADDR", "127.0.0.1:8805"),
	}
	log.Printf("config: addr=%s pfcp=%s", cfg.Address, cfg.PFCPAddr)
	return cfg
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" { return v }
	return def
}
