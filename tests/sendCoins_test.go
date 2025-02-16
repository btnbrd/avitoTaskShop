package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/btnbrd/avitoshop/internal/models"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"testing"
)

func TestSend(t *testing.T) {
	receiiverName := "Recy"
	senderName := "Sendy"
	delta := 2
	sender := models.AuthRequest{senderName, "12"}
	receiver := models.AuthRequest{receiiverName, "12"}
	token1 := getAuthToken(t, sender)
	token2 := getAuthToken(t, receiver)
	info1, err := getInfo(token1)
	assert.NoError(t, err)
	info2, err := getInfo(token2)
	assert.NoError(t, err)

	jsonData := fmt.Sprintf(`{"toUser": %q, "amount": %d}`, receiiverName, delta)
	reqBody := bytes.NewBuffer([]byte(jsonData))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/sendCoin", baseURL), reqBody)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token1)
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

	assert.Equal(t, resp.StatusCode, http.StatusOK)

	info11, err := getInfo(token1)
	assert.NoError(t, err)
	info22, err := getInfo(token2)
	assert.NoError(t, err)
	log.Println(info2)
	assert.Equal(t, info11.Coins+delta, info1.Coins)
	assert.Equal(t, info2.Coins+delta, info22.Coins)
	log.Println(info22.CoinHistory.Received)
	weresent, sent, wererec, rec := 0, 0, 0, 0
	for _, v := range info1.CoinHistory.Sent {
		if v.ToUser == receiiverName {
			weresent = v.Amount
		}
	}
	for _, v := range info11.CoinHistory.Sent {
		if v.ToUser == receiiverName {
			sent = v.Amount
		}
	}
	for _, v := range info22.CoinHistory.Received {
		if v.FromUser == senderName {
			rec = v.Amount
		}
	}
	for _, v := range info2.CoinHistory.Received {
		if v.FromUser == senderName {
			wererec = v.Amount
		}
	}
	assert.Equal(t, weresent+delta, sent)
	assert.Equal(t, wererec+delta, rec)

}

func TestSendToYouself(t *testing.T) {
	receiiverName := "Recy"
	senderName := "Recy"
	delta := 2
	sender := models.AuthRequest{senderName, "12"}
	token1 := getAuthToken(t, sender)

	jsonData := fmt.Sprintf(`{"toUser": %q, "amount": %d}`, receiiverName, delta)
	reqBody := bytes.NewBuffer([]byte(jsonData))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/sendCoin", baseURL), reqBody)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token1)
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
	assert.Equal(t, v, "Отправлять можно только другим пользователям")

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
}

func TestSendToAbsent(t *testing.T) {
	receiverName := "Recy121221"
	senderName := "Sendy"
	delta := 2
	sender := models.AuthRequest{senderName, "12"}
	token1 := getAuthToken(t, sender)

	jsonData := fmt.Sprintf(`{"toUser": %q, "amount": %d}`, receiverName, delta)
	reqBody := bytes.NewBuffer([]byte(jsonData))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/sendCoin", baseURL), reqBody)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token1)
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
	assert.Equal(t, v, "Пользователь-получатель не найден")

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
}
