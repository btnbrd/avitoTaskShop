package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/btnbrd/avitoshop/internal/inventory"
	"github.com/btnbrd/avitoshop/internal/models"
	"log"
	"net/http"
	"testing"
	//"time"

	crrnd "crypto/rand"
	"github.com/stretchr/testify/assert"
	"math/rand"
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
var Bankrupt models.AuthRequest = models.AuthRequest{"Bankrupt4", "querty"}

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

//func TestBecameBankrupt(t *testing.T) {
//	bankrupt := models.AuthRequest{fmt.Sprintf("Bankrupt%d", rand.Intn(100)), "querty"}
//	Bankrupt = bankrupt
//
//	token := getAuthToken(t, bankrupt)
//	_, err := getInfo(token)
//	assert.NoError(t, err)
//	jsonData := `{"toUser": "Bernard", "amount": 1000}`
//	reqBody := bytes.NewBuffer([]byte(jsonData))
//
//	req, err := http.NewRequest("POST", fmt.Sprintf("%s/sendCoin", baseURL), reqBody)
//	assert.NoError(t, err)
//	req.Header.Set("Authorization", "Bearer "+token)
//	req.Header.Set("Content-Type", "application/json")
//
//	client := &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		fmt.Println("Ошибка отправки запроса:", err)
//		return
//	}
//	defer resp.Body.Close()
//
//	var responseBody map[string]interface{}
//	err = json.NewDecoder(resp.Body).Decode(&responseBody)
//	assert.NoError(t, err)
//	v, ok := responseBody["message"]
//	log.Println(v)
//	assert.True(t, ok)
//	log.Println(responseBody["message"])
//
//	info, err := getInfo(token)
//	assert.NoError(t, err)
//	assert.Equal(t, info.Coins, 0)
//
//}

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

// Главная функция для запуска тестов
func main() {
	testing.RunTests(func(pat, str string) (bool, error) {
		return true, nil
	}, []testing.InternalTest{{"TestCreate", TestSignIN},
		{"TestBuyMerch", TestBuyMerch},
		//{"TestBankruptBecome", TestBecameBankrupt},
		{"BankruptSend", TestBankruptSend},
		{"Bankrupt buy", TestBankruptBuy},
		{"Send", TestSend},
		{"Send to yourself", TestSendToYouself},
		{"Send to absent", TestSendToAbsent}})

}
