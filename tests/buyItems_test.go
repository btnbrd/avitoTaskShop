package main

import (
	"encoding/json"
	"fmt"
	"github.com/btnbrd/avitoshop/internal/inventory"
	"log"
	"net/http"
	"testing"
	//"time"

	"github.com/stretchr/testify/assert"
	"math/rand"
)

func TestBuyMerch(t *testing.T) {
	user := users[rand.Intn(3)]
	log.Println(user)
	token := getAuthToken(t, user)
	info, err := getInfo(token)
	assert.NoError(t, err)
	var item string
	for item, _ = range inventory.MerchItems {
		break
	}
	item_cost := inventory.MerchItems[item]
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

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseBody map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	_, ok := responseBody["message"]
	assert.True(t, ok, "Ответ не содержит ожидаемое сообщение о успешной покупке")

	info2, err := getInfo(token)
	assert.NoError(t, err)
	assert.Equal(t, info2.Coins+item_cost, info.Coins)
	if len(info.Inventory) == 0 {
		assert.Equal(t, len(info2.Inventory), 1)
	} else {
		q, qq := 0, 0
		for _, v := range info.Inventory {
			if v.Type == item {
				q = v.Quantity
			}
		}
		for _, v := range info2.Inventory {
			if v.Type == item {
				qq = v.Quantity
			}
		}
		assert.Equal(t, q+1, qq)
	}
}

func TestBuyUnavailableMerch(t *testing.T) {
	user := users[rand.Intn(3)]
	log.Println(user)
	token := getAuthToken(t, user)
	info, err := getInfo(token)
	assert.NoError(t, err)
	var item string = "jfljnl"
	item_cost := inventory.MerchItems[item]
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

	var responseBody map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	v, ok := responseBody["errors"]
	assert.True(t, ok, "Ответ не содержит ожидаемое сообщение о покупке")
	assert.Equal(t, "Неверный тип мерча", v)
	info2, err := getInfo(token)
	assert.NoError(t, err)
	assert.Equal(t, info2.Coins+item_cost, info.Coins)
	if len(info.Inventory) == 0 {
		assert.Equal(t, len(info2.Inventory), 0)
	} else {
		q, qq := 0, 0
		for _, v := range info.Inventory {
			if v.Type == item {
				q = v.Quantity
			}
		}
		for _, v := range info2.Inventory {
			if v.Type == item {
				qq = v.Quantity
			}
		}
		assert.Equal(t, q, qq)
	}
}
