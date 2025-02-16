package dbconfig

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDBConfig(t *testing.T) {
	tempConfigFile := "config.yaml"
	content := `
db:
  name: "test_db"
  user: "test_user"
  password: "test_password"
  host: "localhost"
  port: 5432
`
	err := os.WriteFile(tempConfigFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp config file: %v", err)
	}
	defer os.Remove(tempConfigFile)

	dbConfig, err := NewDBConfig()
	if err != nil {
		t.Fatalf("Failed to load database config: %v", err)
	}

	assert.Equal(t, "test_db", dbConfig.Name, "Database name mismatch")
	assert.Equal(t, "test_user", dbConfig.User, "Database user mismatch")
	assert.Equal(t, "test_password", dbConfig.Password, "Database password mismatch")
	assert.Equal(t, "localhost", dbConfig.Host, "Database host mismatch")
	assert.Equal(t, 5432, dbConfig.Port, "Database port mismatch")
}

func TestNewDBConfig_FileNotFound(t *testing.T) {
	tempConfigFile := "config.yaml"
	if _, err := os.Stat(tempConfigFile); !os.IsNotExist(err) {
		os.Remove(tempConfigFile)
	}
	_, err := NewDBConfig()
	if err == nil {
		t.Fatalf("Expected error when config file is not found, but got nil")
	}

	assert.Contains(t, err.Error(), "error reading config file", "Error message mismatch")
}
