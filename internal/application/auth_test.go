package application

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/btnbrd/avitoshop/internal/models"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	//"time"

	"github.com/gin-gonic/gin"
	//"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/assert"
	//"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		panic(fmt.Sprintf("Could not connect to docker: %s", err))
	}

	resource, err := pool.Run("postgres", "13-alpine", []string{
		"POSTGRES_USER=test",
		"POSTGRES_PASSWORD=test",
		"POSTGRES_DB=testdb",
	})
	if err != nil {
		panic(fmt.Sprintf("Could not start resource: %s", err))
	}

	_ = resource.Expire(60)

	if err := pool.Retry(func() error {
		var err error
		db, err = sql.Open("postgres",
			fmt.Sprintf("postgres://test:test@localhost:%s/testdb?sslmode=disable",
				resource.GetPort("5432/tcp")))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		panic(fmt.Sprintf("Could not connect to docker: %s", err))
	}

	// Применяем миграции
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			coins INTEGER NOT NULL
		)`)
	if err != nil {
		panic(fmt.Sprintf("Failed to create tables: %s", err))
	}

	code := m.Run()

	// Очистка
	if err := pool.Purge(resource); err != nil {
		panic(fmt.Sprintf("Could not purge resource: %s", err))
	}

	os.Exit(code)
}

func TestAuthHandler_Integration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	app := &APIServer{db: db}

	t.Run("Successful registration and login", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(
			http.MethodPost,
			"/auth",
			strings.NewReader(`{"username":"newuser","password":"testpass"}`),
		)
		c.Request.Header.Set("Content-Type", "application/json")

		app.authHandler(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var authResponse models.AuthResponse
		err := json.Unmarshal(w.Body.Bytes(), &authResponse)
		assert.NoError(t, err)
		assert.NotEmpty(t, authResponse.Token)

		var user models.User
		err = db.QueryRow("SELECT username, coins FROM users WHERE username = $1", "newuser").Scan(&user.Username, &user.Coins)
		assert.NoError(t, err)
		assert.Equal(t, "newuser", user.Username)
		assert.Equal(t, 1000, user.Coins)
	})
}

func cleanupDB(t *testing.T) {
	_, err := db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	assert.NoError(t, err)
}
