package application

import (
	"database/sql"
	"fmt"
	"github.com/btnbrd/avitoshop/internal/inventory"
	"github.com/btnbrd/avitoshop/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (app *APIServer) infoHandler(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Нет информации о пользователе"})
		return
	}
	userID := userIDInterface.(int)

	var coins int
	err := app.db.QueryRow("SELECT coins FROM users WHERE id=$1", userID).Scan(&coins)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Ошибка получения баланса"})
		return
	}

	// Получаем историю полученных монет
	receivedRows, err := app.db.Query(`
		SELECT u.username, SUM(ct.amount) FROM coin_transfers ct
		JOIN users u ON ct.from_user_id = u.id
		WHERE ct.to_user_id = $1
		GROUP BY u.username`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Ошибка получения истории переводов"})
		return
	}
	defer receivedRows.Close()
	var received []models.ReceivedCoin
	for receivedRows.Next() {
		var rc models.ReceivedCoin
		if err := receivedRows.Scan(&rc.FromUser, &rc.Amount); err == nil {
			received = append(received, rc)
		}
	}

	// Получаем историю отправленных монет
	sentRows, err := app.db.Query(`
		SELECT u.username, SUM(ct.amount) AS total_amount
		FROM coin_transfers ct
		JOIN users u ON ct.to_user_id = u.id
		WHERE ct.from_user_id = $1
		GROUP BY u.username`, userID)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Ошибка получения истории переводов"})
		return
	}
	defer sentRows.Close()
	var sent []models.SentCoin
	for sentRows.Next() {
		var sc models.SentCoin
		if err := sentRows.Scan(&sc.ToUser, &sc.Amount); err == nil {
			sent = append(sent, sc)
		}
	}

	// Получаем купленные товары и группируем по типу
	rows, err := app.db.Query(`SELECT item, COUNT(*) FROM purchases WHERE user_id=$1 GROUP BY item`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Ошибка получения списка покупок"})
		return
	}
	defer rows.Close()
	var inventory []models.InventoryItem
	for rows.Next() {
		var item models.InventoryItem
		if err := rows.Scan(&item.Type, &item.Quantity); err == nil {
			inventory = append(inventory, item)
		}
	}

	resp := models.InfoResponse{
		Coins:     coins,
		Inventory: inventory,
	}
	resp.CoinHistory.Received = received
	resp.CoinHistory.Sent = sent

	c.JSON(http.StatusOK, resp)
}

func (app *APIServer) sendCoinHandler(c *gin.Context) {
	var req models.SendCoinRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.ToUser == "" || req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "Неверный запрос"})
		return
	}

	// Получаем отправителя из контекста
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Нет информации о пользователе"})
		return
	}
	fromUserID := userIDInterface.(int)

	receiver, err := app.getUserByUsername(req.ToUser)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"errors": "Пользователь-получатель не найден"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Ошибка сервера"})
		return
	}

	if receiver.ID == fromUserID {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "Отправлять можно только другим пользователям"})
		return
	}

	tx, err := app.db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Ошибка начала транзакции"})
		return
	}
	defer tx.Rollback()

	var senderCoins int
	err = tx.QueryRow("SELECT coins FROM users WHERE id=$1 FOR UPDATE", fromUserID).Scan(&senderCoins)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Ошибка получения баланса отправителя"})
		return
	}
	if senderCoins < req.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "Недостаточно монет"})
		return
	}

	_, err = tx.Exec("UPDATE users SET coins = coins - $1 WHERE id=$2", req.Amount, fromUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Ошибка обновления баланса отправителя"})
		return
	}

	_, err = tx.Exec("UPDATE users SET coins = coins + $1 WHERE id=$2", req.Amount, receiver.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Ошибка обновления баланса получателя"})
		return
	}

	_, err = tx.Exec("INSERT INTO coin_transfers(from_user_id, to_user_id, amount) VALUES($1, $2, $3)", fromUserID, receiver.ID, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Ошибка записи перевода"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Ошибка коммита транзакции"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Монеты успешно отправлены"})
}

func (app *APIServer) buyHandler(c *gin.Context) {
	item := c.Param("item")
	price, exists := inventory.MerchItems[item]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "Неверный тип мерча"})
		return
	}

	userIDInterface, existsCtx := c.Get("userID")
	if !existsCtx {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Нет информации о пользователе"})
		return
	}
	userID := userIDInterface.(int)

	tx, err := app.db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Ошибка начала транзакции"})
		return
	}
	defer tx.Rollback()

	var coins int
	err = tx.QueryRow("SELECT coins FROM users WHERE id=$1 FOR UPDATE", userID).Scan(&coins)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Ошибка получения баланса"})
		return
	}
	if coins < price {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "Недостаточно монет для покупки"})
		return
	}

	_, err = tx.Exec("UPDATE users SET coins = coins - $1 WHERE id=$2", price, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Ошибка обновления баланса"})
		return
	}

	_, err = tx.Exec("INSERT INTO purchases(user_id, item, price) VALUES($1, $2, $3)", userID, item, price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Ошибка записи покупки"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Ошибка коммита транзакции"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Покупка '%s' за %d монет прошла успешно", item, price)})
}
