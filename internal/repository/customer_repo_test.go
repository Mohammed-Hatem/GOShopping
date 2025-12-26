package repository

import (
	"bookstore-project/internal/config"
	"bookstore-project/internal/database"
	"testing"
)

func setupCustomerTestDB(t *testing.T) (*CustomerRepo, func()) {
	err := config.NewConfig()
	if err != nil {
		t.Fatal(err)
	}

	db, err := database.New()
	if err != nil {
		t.Fatal(err)
	}

	repo := NewCustomerRepo(db)
	cleanup := func() {
		db.Close()
	}

	return repo, cleanup
}

func TestCreateCustomer(t *testing.T) {
	repo, cleanup := setupCustomerTestDB(t)
	defer cleanup()

	username := "test_customer_user"
	email := "test_customer_email@example.com"

	// Ensure a clean state
	_, _ = repo.db.Exec("DELETE FROM customer WHERE username = $1 OR email = $2", username, email)

	// 1) Successful creation
	err := repo.CreateCustomer(username, "password123", "John", "Doe", email, "123-456-7890", "123 Main St")
	if err != nil {
		t.Fatalf("expected no error on first create, got %v", err)
	}

	// 2) Duplicate username (different email)
	err = repo.CreateCustomer(username, "password456", "Jane", "Smith", "other_"+email, "000-000-0000", "Other St")
	if err == nil {
		t.Fatal("expected error for duplicate username, got nil")
	}
	if err.Error() != "username already exists" {
		t.Errorf("expected 'username already exists', got %v", err)
	}

	// 3) Duplicate email (different username)
	err = repo.CreateCustomer("other_"+username, "password789", "Bob", "Jones", email, "999-999-9999", "Another St")
	if err == nil {
		t.Fatal("expected error for duplicate email, got nil")
	}
	if err.Error() != "email already exists" {
		t.Errorf("expected 'email already exists', got %v", err)
	}

	// Cleanup
	_, _ = repo.db.Exec("DELETE FROM customer WHERE username = $1 OR email = $2 OR username = $3 OR email = $4", username, email, "other_"+username, "other_"+email)
}

func TestGetCustomerByUsername(t *testing.T) {
	repo, cleanup := setupCustomerTestDB(t)
	defer cleanup()

	username := "get_by_username_user"
	email := "get_by_username@example.com"
	_, _ = repo.db.Exec("DELETE FROM customer WHERE username = $1 OR email = $2", username, email)

	// Create customer
	err := repo.CreateCustomer(username, "password123", "Alice", "Wonder", email, "555-1234", "Wonderland")
	if err != nil {
		t.Fatalf("failed to create test customer: %v", err)
	}
	defer func() {
		_, _ = repo.db.Exec("DELETE FROM customer WHERE username = $1", username)
	}()

	// Fetch existing
	cust, err := repo.GetCustomerByUsername(username)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cust == nil {
		t.Fatal("expected customer, got nil")
	}
	if cust.Username != username {
		t.Errorf("expected username %q, got %q", username, cust.Username)
	}
}

func TestGetCustomerByEmail(t *testing.T) {
	repo, cleanup := setupCustomerTestDB(t)
	defer cleanup()

	username := "get_by_email_user"
	email := "get_by_email@example.com"
	_, _ = repo.db.Exec("DELETE FROM customer WHERE username = $1 OR email = $2", username, email)

	// Create customer
	err := repo.CreateCustomer(username, "password123", "Eve", "Adams", email, "555-5678", "Garden")
	if err != nil {
		t.Fatalf("failed to create test customer: %v", err)
	}
	defer func() {
		_, _ = repo.db.Exec("DELETE FROM customer WHERE username = $1", username)
	}()

	// Fetch existing
	cust, err := repo.GetCustomerByEmail(email)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cust == nil {
		t.Fatal("expected customer, got nil")
	}
	if cust.Email != email {
		t.Errorf("expected email %q, got %q", email, cust.Email)
	}
}

func TestGetOrdersByCustomer(t *testing.T) {
	repo, cleanup := setupCustomerTestDB(t)
	defer cleanup()

	// Uses seeded customer "alice" from 003_seed_data.sql
	orders, err := repo.GetOrdersByCustomer("alice")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if orders == nil {
		t.Error("orders should not be nil")
	}
}

func TestGetOrderItemsByOrderID(t *testing.T) {
	repo, cleanup := setupCustomerTestDB(t)
	defer cleanup()

	// Uses seeded order_id 1 from 003_seed_data.sql
	items, err := repo.GetOrderItemsByOrderID(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if items == nil {
		t.Error("items should not be nil")
	}
}
