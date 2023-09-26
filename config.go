package config

import (
  "os"
  "strconv"
)

// Configuration structs
type Config struct {
  Port int
  PythonAPIURL string
}

// Load configuration
func LoadConfig() Config {

  // Set defaults
  cfg := Config{
    Port: 8080,
    PythonAPIURL: "http://localhost:5000", 
  }

  // Override with environment variables if set
  if os.Getenv("PORT") != "" {
    port, _ := strconv.Atoi(os.Getenv("PORT"))
    cfg.Port = port
  }

  if os.Getenv("PYTHON_API_URL") != "" {
    cfg.PythonAPIURL = os.Getenv("PYTHON_API_URL")
  }

  return cfg
}

// Exported config
var Config = LoadConfig()