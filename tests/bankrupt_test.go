package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/btnbrd/avitoshop/internal/inventory"
	"github.com/btnbrd/avitoshop/internal/models"
	"github.com/stretchr/testify/assert"
	"log"
	"math/rand"
	"net/http"
	"testing"
)

func TestBecameBankrupt(t *testing.T) {
	bankrupt := models.AuthRequest{fmt.Sprintf("Bankrupt%d", rand.Intn(100)), "querty"}
	Bankrupt = bankrupt

	token := getAuthToken(t, bankrupt)
	_, err := getInfo(token)
	assert.NoError(t, err)
	jsonData := `{"toUser": "Bernard", "amount": 1000}`
	reqBody := bytes.NewBuffer([]byte(jsonData))

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/sendCoin", baseURL), reqBody)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка отправки запроса:", err)
		return
	}
	defer resp.Body.Close()

	var responseBody map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	log.Println(responseBody)
	assert.NoError(t, err)
	v, ok := responseBody["message"]
	log.Println(v)
	assert.True(t, ok)
	log.Println(responseBody["message"])

	info, err := getInfo(token)
	assert.NoError(t, err)
	assert.Equal(t, info.Coins, 0)

}

func TestBankruptSend(t *testing.T) {
	log.Println(Bankrupt)
	token := getAuthToken(t, Bankrupt)

	jsonData := `{"toUser": "Bernard", "amount": 1000}`
	reqBody := bytes.NewBuffer([]byte(jsonData))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/sendCoin", baseURL), reqBody)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка отправки запроса:", err)
		return
	}
	defer resp.Body.Close()
	var responseBody map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	v, ok := responseBody["errors"]
	assert.True(t, ok)
	assert.Equal(t, v, "Недостаточно монет")
	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
}

func TestBankruptBuy(t *testing.T) {
	log.Println(Bankrupt)
	token := getAuthToken(t, Bankrupt)

	var item string
	for item, _ = range inventory.MerchItems {
		break
	}
	_ = inventory.MerchItems[item]
	reqURL := fmt.Sprintf("%s/buy/%s", baseURL, item)

	req, err := http.NewRequest("GET", reqURL, nil)
	assert.NoError(t, err)

	client := &http.Client{}
	resp1, err := client.Do(req)
	defer resp1.Body.Close()
	assert.NoError(t, err)
	assert.Equal(t, resp1.StatusCode, http.StatusUnauthorized)

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
