package application

//
//import (
//	"database/sql"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//	"testing"
//)
//
//// MockDB реализует интерфейс storage.DBInterface
//type MockDB struct {
//	mock.Mock
//}
//
//// Реализация методов DBInterface
//func (m *MockDB) Exec(query string, args ...any) (sql.Result, error) {
//	args = append([]any{query}, args...)
//	return m.Called(args...).Get(0).(sql.Result), m.Called(args...).Error(1)
//}
//
//func (m *MockDB) QueryRow(query string, args ...any) *sql.Row {
//	args = append([]any{query}, args...)
//	m.Called(args...)
//	return &sql.Row{}
//}
//
//func (m *MockDB) Query(query string, args ...any) (*sql.Rows, error) {
//	args = append([]any{query}, args...)
//	return m.Called(args...).Get(0).(*sql.Rows), m.Called(args...).Error(1)
//}
//
//func (m *MockDB) Begin() (*sql.Tx, error) {
//	return m.Called().Get(0).(*sql.Tx), m.Called().Error(1)
//}
//
//func TestNewServer(t *testing.T) {
//	// Создаем мок для DBInterface
//	mockDB := new(MockDB)
//
//	// Вызываем тестируемую функцию
//	server := NewServer(mockDB)
//
//	// Проверяем, что возвращаемый объект не nil
//	assert.NotNil(t, server, "Возвращенный APIServer не должен быть nil")
//
//	// Проверяем, что поле r (маршрутизатор Gin) инициализировано
//	assert.NotNil(t, server.r, "Поле r (маршрутизатор Gin) должно быть инициализировано")
//
//	// Проверяем, что поле db содержит переданный мок
//	assert.Equal(t, mockDB, server.db, "Поле db должно содержать переданный объект DBInterface")
//}
//
////func TestAPIServer_Run(t *testing.T) {
////	// Создаем mock-объект для базы данных
////	mockDB := new(MockDB)
////
////	// Создаем новый сервер с mock-базой данных
////	server := NewServer(mockDB)
////
////	// Настройка Gin в режиме тестирования
////	gin.SetMode(gin.TestMode)
////
////	// Проверяем маршрут /auth
////	authReq, _ := http.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer([]byte(`{"username":"test","password":"test"}`)))
////	authW := httptest.NewRecorder()
////	server.r.ServeHTTP(authW, authReq)
////
////	// Проверяем, что маршрут /auth работает
////	assert.Equal(t, http.StatusOK, authW.Code, "Route /auth should respond with 200 OK")
////
////	// Проверяем защищенный маршрут /info без аутентификации
////	infoReq, _ := http.NewRequest(http.MethodGet, "/info", nil)
////	infoW := httptest.NewRecorder()
////	server.r.ServeHTTP(infoW, infoReq)
////
////	// Проверяем, что маршрут /info требует аутентификации
////	assert.Equal(t, http.StatusUnauthorized, infoW.Code, "Route /info should require authentication")
////
////	// Проверяем защищенный маршрут /sendCoin без аутентификации
////	sendCoinReq, _ := http.NewRequest(http.MethodPost, "/sendCoin", bytes.NewBuffer([]byte(`{"amount":100}`)))
////	sendCoinW := httptest.NewRecorder()
////	server.r.ServeHTTP(sendCoinW, sendCoinReq)
////
////	// Проверяем, что маршрут /sendCoin требует аутентификации
////	assert.Equal(t, http.StatusUnauthorized, sendCoinW.Code, "Route /sendCoin should require authentication")
////
////	// Проверяем защищенный маршрут /buy/:item без аутентификации
////	buyReq, _ := http.NewRequest(http.MethodGet, "/buy/item1", nil)
////	buyW := httptest.NewRecorder()
////	server.r.ServeHTTP(buyW, buyReq)
////
////	// Проверяем, что маршрут /buy/:item требует аутентификации
////	assert.Equal(t, http.StatusUnauthorized, buyW.Code, "Route /buy/:item should require authentication")
////}
