package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"bookstore-project/internal/models"
)

const (
	apiURL = "http://localhost:3001/api"
)

type AuthResponse struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Token    string `json:"token"`
	Error    string `json:"error"`
}

type BookResponse struct {
	Message string      `json:"message"`
	Error   string      `json:"error"`
	Book    models.Book `json:"book"`
}

func main() {
	// Test all endpoints sequentially
	fmt.Println("🚀 Testing Bookstore API Endpoints")
	fmt.Println("==================================")

	// Test public book search endpoints
	testPublicEndpoints()

	// Test admin login and admin endpoints
	testAdminEndpoints()

	fmt.Println("✅ All tests completed!")
}

func testPublicEndpoints() {
	fmt.Println("\n📚 Testing Public Book Search Endpoints")
	fmt.Println("-------------------------------------")

	// Test 1: Get all books
	fmt.Print("1. GET /api/books - ")
	resp, err := http.Get(apiURL + "/books")
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Println("✅ Success")
		body, _ := io.ReadAll(resp.Body)
		var books []models.Book
		json.Unmarshal(body, &books)
		fmt.Printf("   Found %d books\n", len(books))
	} else {
		fmt.Printf("❌ Failed with status %d\n", resp.StatusCode)
	}

	// Test 2: Get book by ISBN
	fmt.Print("2. GET /api/books/9780132350884 - ")
	resp, err = http.Get(apiURL + "/books/9780132350884")
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Println("✅ Success")
	} else {
		fmt.Printf("❌ Failed with status %d\n", resp.StatusCode)
	}

	// Test 3: Get books by category
	fmt.Print("3. GET /api/books/category/Science - ")
	resp, err = http.Get(apiURL + "/books/category/Science")
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Println("✅ Success")
		body, _ := io.ReadAll(resp.Body)
		var books []models.Book
		json.Unmarshal(body, &books)
		fmt.Printf("   Found %d Science books\n", len(books))
	} else {
		fmt.Printf("❌ Failed with status %d\n", resp.StatusCode)
	}

	// Test 4: Get books by author
	fmt.Print("4. GET /api/books/author/Robert%20C.%20Martin - ")
	resp, err = http.Get(apiURL + "/books/author/Robert%20C.%20Martin")
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Println("✅ Success")
		body, _ := io.ReadAll(resp.Body)
		var books []models.Book
		json.Unmarshal(body, &books)
		fmt.Printf("   Found %d books by Robert C. Martin\n", len(books))
	} else {
		fmt.Printf("❌ Failed with status %d\n", resp.StatusCode)
	}
}

func testAdminEndpoints() {
	fmt.Println("\n🔐 Testing Admin Endpoints")
	fmt.Println("----------------------------")

	// Test admin login
	fmt.Print("1. Admin Login - ")
	loginData := map[string]string{
		"username": "admin1",
		"password": "admin_pass_1",
		"role":     "admin",
	}

	body, _ := json.Marshal(loginData)
	resp, err := http.Post(apiURL+"/auth/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var authResp AuthResponse
	body, _ = io.ReadAll(resp.Body)
	json.Unmarshal(body, &authResp)

	if resp.StatusCode == 200 && authResp.Token != "" {
		fmt.Println("✅ Success")
		fmt.Printf("   Got token for user: %s (role: %s)\n", authResp.Username, authResp.Role)

		// Test admin endpoints with token
		testProtectedEndpoints(authResp.Token)
	} else {
		fmt.Printf("❌ Failed: %s\n", authResp.Error)
	}
}

func testProtectedEndpoints(token string) {
	fmt.Println("\n🔒 Testing Protected Admin Endpoints")
	fmt.Println("-----------------------------------")

	// Test 1: Add new book
	fmt.Print("1. POST /api/admin/books - ")
	book := models.Book{
		Isbn:          "9780123456789",
		Title:         "Test Book for API",
		PubYear:       2024,
		SellingPrice:  29.99,
		Category:      "Science",
		AuthorName:    "API Test Author",
		StockQuantity: 15,
		Threshold:     5,
		PublisherId:   1,
	}

	bookBody, _ := json.Marshal(book)
	req, _ := http.NewRequest("POST", apiURL+"/admin/books", bytes.NewBuffer(bookBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 201 {
		fmt.Println("✅ Success")
		body, _ := io.ReadAll(resp.Body)
		var bookResp BookResponse
		json.Unmarshal(body, &bookResp)
		fmt.Printf("   Created book: %s\n", bookResp.Book.Title)
	} else {
		fmt.Printf("❌ Failed with status %d\n", resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("   Error: %s\n", string(body))
	}

	// Test 2: Get publisher orders
	fmt.Print("2. GET /api/admin/publisher-orders - ")
	req, _ = http.NewRequest("GET", apiURL+"/admin/publisher-orders", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Println("✅ Success")
		body, _ := io.ReadAll(resp.Body)
		var orders []models.PublisherOrder
		json.Unmarshal(body, &orders)
		fmt.Printf("   Found %d publisher orders\n", len(orders))
	} else {
		fmt.Printf("❌ Failed with status %d\n", resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("   Error: %s\n", string(body))
	}

	// Test 3: Update book stock (trigger test)
	fmt.Print("3. PATCH /api/admin/books/9780553103540/stock - ")
	stockUpdate := map[string]int{"quantity": 3} // Below threshold to test trigger
	stockBody, _ := json.Marshal(stockUpdate)
	req, _ = http.NewRequest("PATCH", apiURL+"/admin/books/9780553103540/stock", bytes.NewBuffer(stockBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Println("✅ Success (Stock updated, trigger should create automatic order)")
	} else {
		fmt.Printf("❌ Failed with status %d\n", resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("   Error: %s\n", string(body))
	}

	// Test 4: Confirm publisher order
	fmt.Print("4. PUT /api/admin/publisher-orders/3/confirm - ")
	req, _ = http.NewRequest("PUT", apiURL+"/admin/publisher-orders/3/confirm", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Println("✅ Success (Order confirmed, stock should be updated)")
	} else {
		fmt.Printf("❌ Failed with status %d\n", resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("   Error: %s\n", string(body))
	}
}
