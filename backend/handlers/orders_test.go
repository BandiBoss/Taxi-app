package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"Taxi-app/backend/repository"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) CreateOrder(customerID int, origin, destination string) (int, error) {
	args := m.Called(customerID, origin, destination)
	return args.Int(0), args.Error(1)
}
func (m *MockOrderRepository) GetOrders(customerID, limit, offset int, sortField, orderDir string) ([]repository.Order, error) {
	args := m.Called(customerID, limit, offset, sortField, orderDir)
	return args.Get(0).([]repository.Order), args.Error(1)
}
func (m *MockOrderRepository) GetOrderByID(orderID, customerID int) (*repository.Order, error) {
	args := m.Called(orderID, customerID)
	return args.Get(0).(*repository.Order), args.Error(1)
}
func (m *MockOrderRepository) GetOrderDetailsByID(orderID, customerID int) (*repository.OrderDetails, error) {
	args := m.Called(orderID, customerID)
	return args.Get(0).(*repository.OrderDetails), args.Error(1)
}
func (m *MockOrderRepository) GetOrderLocationHistory(orderID, customerID int, limit int) ([]repository.Location, error) {
	args := m.Called(orderID, customerID, limit)
	return args.Get(0).([]repository.Location), args.Error(1)
}
func (m *MockOrderRepository) UpdateOrderStatus(orderID int, status string) error     { return nil }
func (m *MockOrderRepository) GetOrderStatus(orderID, customerID int) (string, error) { return "", nil }
func (m *MockOrderRepository) GetRandomActiveDriver() (int, error)                    { return 0, nil }
func (m *MockOrderRepository) AssignDriverAndStart(orderID, driverID int) error       { return nil }
func (m *MockOrderRepository) ListActiveDrivers(limit int) ([]int, error)             { return nil, nil }
func (m *MockOrderRepository) ListCreatedOrders(limit int) ([]int, error)             { return nil, nil }
func (m *MockOrderRepository) SeedCreatedOrders(limit int) error                      { return nil }

func TestCreateOrder_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockOrderRepository)
	h := CreateOrder(mockRepo)
	mockRepo.On("CreateOrder", 1, "Origin", "Destination").Return(123, nil)

	body, _ := json.Marshal(CreateOrderRequest{Origin: "Origin", Destination: "Destination"})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/orders", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("userID", 1)

	h(c)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Order created")
	mockRepo.AssertExpectations(t)
}

func TestCreateOrder_DBError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockOrderRepository)
	h := CreateOrder(mockRepo)
	mockRepo.On("CreateOrder", 1, "Origin", "Destination").Return(0, errors.New("db error"))

	body, _ := json.Marshal(CreateOrderRequest{Origin: "Origin", Destination: "Destination"})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/orders", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("userID", 1)

	h(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to create order")
	mockRepo.AssertExpectations(t)
}

func TestGetOrders_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockOrderRepository)
	h := GetOrders(mockRepo)
	orders := []repository.Order{{ID: 1, CustomerID: 1, Status: "created", Origin: "A", Destination: "B", CreatedAt: "now"}}
	mockRepo.On("GetOrders", 1, 10, 0, "created_at", "DESC").Return(orders, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/orders?page=1&limit=10&sort=created_at&order=DESC", nil)
	c.Set("userID", 1)

	h(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "orders")
	mockRepo.AssertExpectations(t)
}

func TestGetOrders_DBError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockOrderRepository)
	h := GetOrders(mockRepo)
	mockRepo.On("GetOrders", 1, 10, 0, "created_at", "DESC").Return([]repository.Order{}, errors.New("db error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/orders?page=1&limit=10&sort=created_at&order=DESC", nil)
	c.Set("userID", 1)

	h(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Query failed")
	mockRepo.AssertExpectations(t)
}

func TestGetOrderByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockOrderRepository)
	h := GetOrderByID(mockRepo)
	order := &repository.OrderDetails{ID: 1, CustomerID: 1, Status: "created", Origin: "A", Destination: "B", CreatedAt: "now"}
	mockRepo.On("GetOrderDetailsByID", 1, 1).Return(order, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Set("userID", 1)

	h(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "created")
	mockRepo.AssertExpectations(t)
}

func TestGetOrderByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockOrderRepository)
	h := GetOrderByID(mockRepo)
	mockRepo.On("GetOrderDetailsByID", 1, 1).Return((*repository.OrderDetails)(nil), errors.New("not found"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Set("userID", 1)

	h(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Order not found")
	mockRepo.AssertExpectations(t)
}

func TestGetOrderLocationHistory_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockOrderRepository)
	h := GetOrderLocationHistory(mockRepo)
	hist := []repository.Location{{Latitude: 1.1, Longitude: 2.2, GeneratedTime: "now"}}
	mockRepo.On("GetOrderLocationHistory", 1, 1, 50).Return(hist, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Set("userID", 1)

	h(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "history")
	mockRepo.AssertExpectations(t)
}

func TestGetOrderLocationHistory_DBError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockOrderRepository)
	h := GetOrderLocationHistory(mockRepo)
	mockRepo.On("GetOrderLocationHistory", 1, 1, 50).Return([]repository.Location{}, errors.New("db error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Set("userID", 1)

	h(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Query failed")
	mockRepo.AssertExpectations(t)
}
