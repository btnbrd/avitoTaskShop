# avitoTaskShop

### Сервис реализующий работу магазина мерча

## Инструкция по запуску

Проект запускается с помощью docker-compose. Для этого перейдите в корневую папку проекта и в консоли запустите команду:

```
    docker-compose up --build
```

Api апускается на 8080 порту по умолчанию

## Инструкция по E2E тестам

Тесты на различные пользовательские сценарии доступны в папки tests. Запускаются через

```
go test -v github.com/btnbrd/avitoshop/tests -run '' tests/e2e_test.go
```

## Взаимодействие

После того, как проект запустится вам доступны следующие запросы
## POST

### `/api/auth`

```http request
POST http://localhost:8080/auth
```

Аутентификация пользователя. Принимает JSON:

```json
{
  "username": "user1",
  "password": "password123"
}
```

#### Ответ:

```json
{
  "token": "your_jwt_token"
}
```

После аутентификации нужно использовать полученный `token` для доступа к защищённым маршрутам, добавляя его в `Authorization` заголовок:  
`Authorization: Bearer your_jwt_token`

---

### `/api/sendCoin`

```http request
POST http://localhost:8080/sendCoin
```

Отправляет монеты другому пользователю. Принимает JSON:

```json
{
  "toUser": "user2",
  "amount": 100
}
```

Отправка монет самому себе запрещена и возврщает Ошибку **400** -- неверный запрос.

#### Возможные ответы:

- **200** – Монеты успешно отправлены.
- **400** – Неверный запрос.
- **401** – Неавторизованный пользователь.
- **500** – Внутренняя ошибка сервера.

---

## GET

### `/api/info`

```http request
GET http://localhost:8080/info
```

Возвращает информацию о пользователе: баланс монет, инвентарь и историю транзакций.

#### Ответ:

```json
{
  "coins": 500,
  "inventory": [
    {
      "type": "sword",
      "quantity": 1
    },
    {
      "type": "shield",
      "quantity": 2
    }
  ],
  "coinHistory": {
    "received": [
      {
        "fromUser": "user3",
        "amount": 200
      }
    ],
    "sent": [
      {
        "toUser": "user2",
        "amount": 100
      }
    ]
  }
}
```

---

### `/api/buy/{item}`

```http request
GET http://localhost:8080/buy/sword
```

Покупает предмет за монеты.

#### Возможные ответы:

- **200** – "Успешный ответ."
- **400** – "Неверный запрос."
- **401** – "Неавторизован."
- **500** – "Внутренняя ошибка сервера."
