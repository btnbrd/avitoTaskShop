package storage

import (
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/btnbrd/avitoshop/internal/storage/dbconfig"
	_ "github.com/lib/pq"
)

func createTestDatabase() (*sql.DB, error) {
	conf, err := dbconfig.NewDBConfig()
	if err != nil {
		return nil, err
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable",
		conf.Host, conf.Port, conf.User, conf.Password)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	testDBName := "test_avitoshop_db_" + strconv.FormatInt(time.Now().UnixNano(), 10)

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", testDBName))
	if err != nil {
		return nil, fmt.Errorf("failed to create test database: %w", err)
	}

	psqlTestInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.Host, conf.Port, conf.User, conf.Password, testDBName)

	testDB, err := sql.Open("postgres", psqlTestInfo)
	if err != nil {
		return nil, err
	}

	return testDB, nil
}

func dropTestDatabase(db *sql.DB, dbName string) error {
	conf, err := dbconfig.NewDBConfig()
	if err != nil {
		return err
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable",
		conf.Host, conf.Port, conf.User, conf.Password)

	adminDB, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}
	defer adminDB.Close()

	_, err = adminDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
	if err != nil {
		return fmt.Errorf("failed to drop test database: %w", err)
	}

	return nil
}

func TestNewDBConnection(t *testing.T) {
	db, err := NewDBConnection()
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping the database: %v", err)
	}
}

func TestInitDB(t *testing.T) {
	testDB, err := createTestDatabase()
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer func() {
		err := dropTestDatabase(testDB, "test_avitoshop_db_"+strconv.FormatInt(time.Now().UnixNano(), 10))
		if err != nil {
			t.Fatalf("Failed to drop test database: %v", err)
		}
	}()

	err = InitDB(testDB)
	if err != nil {
		t.Fatalf("Failed to initialize tables: %v", err)
	}

	rows, err := testDB.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'")
	if err != nil {
		t.Fatalf("Failed to query tables: %v", err)
	}
	defer rows.Close()

	expectedTables := map[string]bool{
		"users":          true,
		"coin_transfers": true,
		"purchases":      true,
	}

	tablesFound := 0
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			t.Errorf("Error scanning table name: %v", err)
			continue
		}

		if expectedTables[tableName] {
			tablesFound++
		}
	}

	if tablesFound != len(expectedTables) {
		t.Errorf("Expected %d tables, but found %d", len(expectedTables), tablesFound)
	}
}
