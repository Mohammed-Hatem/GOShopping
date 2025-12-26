package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"bookstore-project/internal/models"
)

const (
	baseURL = "http://localhost:3001/api"
)

type TestUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Token    string
}

type TestResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
}

var (
	adminToken    string
	customerToken string
	testISBN      = "9780123456789"
)

func TestMain(m *testing.M) {
	// Wait a moment for server to be ready
	time.Sleep(2 * time.Second)

	// Run tests
	code := m.Run()

	os.Exit(code)
}

func TestAuthentication(t *testing.T) {
	t.Run("Admin Signup", func(t *testing.T) {
		user := TestUser{
			Username: "testadmin",
			Password: "adminpass123",
			Email:    "admin@test.com",
		}

		body, _ := json.Marshal(user)
		resp, err := http.Post(baseURL+"/auth/signup", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to signup admin: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", resp.StatusCode)
		}
	})

	t.Run("Admin Login", func(t *testing.T) {
		loginData := map[string]string{
			"username": "admin1", // Use existing admin from seed data
			"password": "admin_pass_1",
		}

		body, _ := json.Marshal(loginData)
		resp, err := http.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to login admin: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		if token, ok := response["token"].(string); ok {
			adminToken = token
		} else {
			t.Error("No token received in login response")
		}
	})

	t.Run("Customer Login", func(t *testing.T) {
		loginData := map[string]string{
			"username": "alice", // Use existing customer from seed data
			"password": "hashed_password_1",
		}

		body, _ := json.Marshal(loginData)
		resp, err := http.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to login customer: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		if token, ok := response["token"].(string); ok {
			customerToken = token
		} else {
			t.Error("No token received in login response")
		}
	})
}

func TestBookSearch(t *testing.T) {
	t.Run("Get All Books", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/books")
		if err != nil {
			t.Fatalf("Failed to get all books: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var books []models.Book
		json.NewDecoder(resp.Body).Decode(&books)

		if len(books) == 0 {
			t.Error("No books found")
		}
	})

	t.Run("Get Book by ISBN", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/books/9780132350884")
		if err != nil {
			t.Fatalf("Failed to get book by ISBN: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var book models.Book
		json.NewDecoder(resp.Body).Decode(&book)

		if book.Isbn != "9780132350884" {
			t.Errorf("Expected ISBN 9780132350884, got %s", book.Isbn)
		}
	})

	t.Run("Get Books by Category", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/books/category/Science")
		if err != nil {
			t.Fatalf("Failed to get books by category: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var books []models.Book
		json.NewDecoder(resp.Body).Decode(&books)

		if len(books) == 0 {
			t.Error("No Science books found")
		}
	})

	t.Run("Get Books by Author", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/books/author/Robert%20C.%20Martin")
		if err != nil {
			t.Fatalf("Failed to get books by author: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var books []models.Book
		json.NewDecoder(resp.Body).Decode(&books)

		if len(books) == 0 {
			t.Error("No books by Robert C. Martin found")
		}
	})
}

func TestAdminBookManagement(t *testing.T) {
	if adminToken == "" {
		t.Skip("Admin token not available")
	}

	t.Run("Add New Book", func(t *testing.T) {
		book := models.Book{
			Isbn:          testISBN,
			Title:         "Test Book for Go",
			PubYear:       2024,
			SellingPrice:  29.99,
			Category:      "Science",
			AuthorName:    "Test Author",
			StockQuantity: 15,
			Threshold:     5,
			PublisherId:   1,
		}

		body, _ := json.Marshal(book)
		req, _ := http.NewRequest("POST", baseURL+"/admin/books", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+adminToken)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to add book: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", resp.StatusCode)
		}
	})

	t.Run("Update Book", func(t *testing.T) {
		book := models.Book{
			Title:        "Updated Test Book",
			SellingPrice: 34.99,
			Category:     "Art",
		}

		body, _ := json.Marshal(book)
		req, _ := http.NewRequest("PUT", baseURL+"/admin/books/"+testISBN, bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+adminToken)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to update book: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("Update Book Stock", func(t *testing.T) {
		stockUpdate := map[string]int{
			"quantity": 12,
		}

		body, _ := json.Marshal(stockUpdate)
		req, _ := http.NewRequest("PATCH", baseURL+"/admin/books/"+testISBN+"/stock", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+adminToken)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to update book stock: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("Test Stock Below Threshold Trigger", func(t *testing.T) {
		// Update stock to below threshold to test automatic ordering trigger
		stockUpdate := map[string]int{
			"quantity": 3, // Below threshold of 5
		}

		body, _ := json.Marshal(stockUpdate)
		req, _ := http.NewRequest("PATCH", baseURL+"/admin/books/"+testISBN+"/stock", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+adminToken)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to update book stock: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		// Wait a moment for trigger to execute
		time.Sleep(100 * time.Millisecond)
	})
}

func TestPublisherOrderManagement(t *testing.T) {
	if adminToken == "" {
		t.Skip("Admin token not available")
	}

	t.Run("Get All Publisher Orders", func(t *testing.T) {
		req, _ := http.NewRequest("GET", baseURL+"/admin/publisher-orders", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to get publisher orders: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var orders []models.PublisherOrder
		json.NewDecoder(resp.Body).Decode(&orders)

		if len(orders) == 0 {
			t.Error("No publisher orders found")
		}
	})

	t.Run("Get Pending Publisher Orders", func(t *testing.T) {
		req, _ := http.NewRequest("GET", baseURL+"/admin/publisher-orders/pending", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to get pending publisher orders: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("Place Publisher Order", func(t *testing.T) {
		order := models.PublisherOrder{
			ISBN:     "9780132350884",
			Quantity: 25,
			AdminID:  1,
		}

		body, _ := json.Marshal(order)
		req, _ := http.NewRequest("POST", baseURL+"/admin/publisher-orders", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+adminToken)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to place publisher order: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", resp.StatusCode)
		}
	})

	t.Run("Confirm Publisher Order", func(t *testing.T) {
		// First get a pending order to confirm
		req, _ := http.NewRequest("GET", baseURL+"/admin/publisher-orders/pending", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to get pending orders: %v", err)
		}
		defer resp.Body.Close()

		var orders []models.PublisherOrder
		json.NewDecoder(resp.Body).Decode(&orders)

		if len(orders) == 0 {
			t.Skip("No pending orders to confirm")
		}

		// Confirm the first pending order
		orderID := orders[0].ID
		req, _ = http.NewRequest("PUT", baseURL+fmt.Sprintf("/admin/publisher-orders/%d/confirm", orderID), nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)

		resp, err = client.Do(req)
		if err != nil {
			t.Fatalf("Failed to confirm publisher order: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})
}

func TestAuthorization(t *testing.T) {
	t.Run("Customer Cannot Access Admin Endpoints", func(t *testing.T) {
		if customerToken == "" {
			t.Skip("Customer token not available")
		}

		req, _ := http.NewRequest("GET", baseURL+"/admin/publisher-orders", nil)
		req.Header.Set("Authorization", "Bearer "+customerToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("Unauthorized Access Without Token", func(t *testing.T) {
		req, _ := http.NewRequest("POST", baseURL+"/admin/books", nil)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})
}
