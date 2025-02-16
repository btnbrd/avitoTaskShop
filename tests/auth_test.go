package main

import (
	"bytes"
	crrnd "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/btnbrd/avitoshop/internal/models"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"testing"
)

func RandomString(n int) string {
	b := make([]byte, n)
	_, err := crrnd.Read(b)
	if err != nil {
		panic(err) // В реальном коде лучше обработать ошибку
	}
	return base64.URLEncoding.EncodeToString(b)[:n] // Обрезаем до нужной длины
}

const baseURL = "http://localhost:8080"

var users = []models.AuthRequest{
	{"Alex", "a1"},
	{"Bernard", "a11"},
	{"Craig", "a1"},
	{"Empty", ""},
	{"", "a1"},
}
var Bankrupt = models.AuthRequest{"Bankrupt4", "querty"}

var tokens = []string{}

func getAuthToken(t *testing.T, authData models.AuthRequest) string {

	// Преобразуем данные в JSON
	authJSON, err := json.Marshal(authData)
	assert.NoError(t, err)

	// Отправка запроса на аутентификацию
	resp, err := http.Post(fmt.Sprintf("%s/auth", baseURL), "application/json", bytes.NewBuffer(authJSON))
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Проверка успешности аутентификации
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Парсим ответ, чтобы получить токен
	var authResponse struct {
		Token string `json:"token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&authResponse)
	assert.NoError(t, err)

	return authResponse.Token
}

func getInfo(token string) (*models.InfoResponse, error) {
	//token := tokens[id]
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/info", baseURL), nil)
	if err != nil {
		log.Println("Ошибка при создании запроса:", err)
		return nil, err
	}

	// Добавляем заголовок Authorization
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Ошибка при выполнении запроса:", err)
		return nil, err
	}
	defer resp.Body.Close()

	infresp := new(models.InfoResponse)
	err = json.NewDecoder(resp.Body).Decode(infresp)
	if err != nil {
		return nil, err
	}
	return infresp, nil
}

func TestSignIN(t *testing.T) {
	tokens = make([]string, 4)
	for i, authData := range users[:5] {
		if i < 3 {
			token := getAuthToken(t, authData)
			tokens[i] = token
			//log.Println(tokens[i])
			_, err := getInfo(token)
			assert.NoError(t, err)

		} else {
			authJSON, err := json.Marshal(authData)
			assert.NoError(t, err)

			// Отправка запроса на аутентификацию
			resp, err := http.Post(fmt.Sprintf("%s/auth", baseURL), "application/json", bytes.NewBuffer(authJSON))
			assert.NoError(t, err)
			assert.Equal(t, resp.StatusCode, http.StatusBadRequest)

		}
	}
	au := models.AuthRequest{Username: RandomString(10), Password: RandomString(4)}
	token := getAuthToken(t, au)
	inf, err := getInfo(token)
	assert.NoError(t, err)

	assert.Equal(t, inf.Coins, 1000)
	assert.Empty(t, inf.CoinHistory)
	assert.Empty(t, inf.Inventory)
}
