package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

type Product struct {
	ID          string
	ImageURL    string
	Name        string
	Description string
	Price       float64
}

type Order struct {
	ID        string
	ProductID string
	Quantity  int
	Total     float64
	CreatedAt time.Time
}

func initDB() {
	var err error
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"))
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		panic(err)
	}
	fmt.Println("Database connection initialized")
}

func getProductsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getProductsHandler called")
	rows, err := db.Query("SELECT id, image_url, name, description, price FROM products")
	if err != nil {
		fmt.Printf("Error querying products: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.ImageURL, &p.Name, &p.Description, &p.Price); err != nil {
			fmt.Printf("Error scanning product: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
	fmt.Println("getProductsHandler completed successfully")
}

func createProductHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("createProductHandler called")
	if r.Method != http.MethodPost {
		fmt.Println("Invalid request method for createProductHandler")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var newProduct Product
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		fmt.Printf("Error decoding new product: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newProduct.ID = fmt.Sprintf("%d", time.Now().Unix())
	_, err = db.Exec("INSERT INTO products (id, image_url, name, description, price) VALUES ($1, $2, $3, $4, $5)",
		newProduct.ID, newProduct.ImageURL, newProduct.Name, newProduct.Description, newProduct.Price)
	if err != nil {
		fmt.Printf("Error inserting new product: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newProduct)
	fmt.Println("createProductHandler completed successfully")
}

func getProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getProductByIDHandler called")
	idStr := r.URL.Path[len("/products/"):]

	var p Product
	err := db.QueryRow("SELECT id, image_url, name, description, price FROM products WHERE id = $1", idStr).Scan(&p.ID, &p.ImageURL, &p.Name, &p.Description, &p.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("Product not found: %v\n", idStr)
			http.Error(w, "Product not found", http.StatusNotFound)
		} else {
			fmt.Printf("Error querying product by ID: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
	fmt.Println("getProductByIDHandler completed successfully")
}

func deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("deleteProductHandler called")
	if r.Method != http.MethodDelete {
		fmt.Println("Invalid request method for deleteProductHandler")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Path[len("/products/delete/"):]

	_, err := db.Exec("DELETE FROM products WHERE id = $1", idStr)
	if err != nil {
		fmt.Printf("Error deleting product: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	fmt.Println("deleteProductHandler completed successfully")
}

func updateProductHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("updateProductHandler called")
	if r.Method != http.MethodPut {
		fmt.Println("Invalid request method for updateProductHandler")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Path[len("/products/update/"):]

	var updatedProduct Product
	err := json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		fmt.Printf("Error decoding updated product: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE products SET image_url = $1, name = $2, description = $3, price = $4 WHERE id = $5",
		updatedProduct.ImageURL, updatedProduct.Name, updatedProduct.Description, updatedProduct.Price, idStr)
	if err != nil {
		fmt.Printf("Error updating product: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedProduct)
	fmt.Println("updateProductHandler completed successfully")
}

func createOrdersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("createOrdersHandler called")
	if r.Method != http.MethodPost {
		fmt.Println("Invalid request method for createOrdersHandler")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var newOrders []Order
	fmt.Printf("Create orders request body: %+v\n", r.Body)
	err := json.NewDecoder(r.Body).Decode(&newOrders)
	if err != nil {
		fmt.Printf("Error decoding new orders: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Decoded new orders: %+v\n", newOrders)

	tx, err := db.Begin()
	if err != nil {
		fmt.Printf("Error beginning transaction: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Transaction started")

	for _, newOrder := range newOrders {
		var productExists bool
		err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE id=$1)", newOrder.ProductID).Scan(&productExists)
		if err != nil {
			fmt.Printf("Error checking if product exists: %v\n", err)
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !productExists {
			fmt.Printf("Product ID %s does not exist\n", newOrder.ProductID)
			tx.Rollback()
			http.Error(w, fmt.Sprintf("Product ID %s does not exist", newOrder.ProductID), http.StatusBadRequest)
			return
		}

		timestamp := fmt.Sprintf("%d", time.Now().Unix())

		// Инициализируем генератор случайных чисел
		rand.Seed(time.Now().UnixNano())

		// Генерируем случайную цифру от 0 до 9
		randomDigit := rand.Intn(10)

		newOrder.ID = fmt.Sprintf("%s%d", timestamp, randomDigit)
		newOrder.CreatedAt = time.Now()
		_, err = tx.Exec("INSERT INTO orders (id, product_id, quantity, total, created_at) VALUES ($1, $2, $3, $4, $5)",
			newOrder.ID, newOrder.ProductID, newOrder.Quantity, newOrder.Total, newOrder.CreatedAt)
		if err != nil {
			fmt.Printf("Error inserting new order: %v\n", err)
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Printf("Inserted new order: %+v\n", newOrder)
	}

	err = tx.Commit()
	if err != nil {
		fmt.Printf("Error committing transaction: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Transaction committed")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newOrders)
	fmt.Println("createOrdersHandler completed successfully")
}

// обработчик для получения всех заказов
func getOrdersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getOrdersHandler called")
	rows, err := db.Query("SELECT id, product_id, quantity, total, created_at FROM orders")
	if err != nil {
		fmt.Printf("Error querying orders: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.ProductID, &o.Quantity, &o.Total, &o.CreatedAt); err != nil {
			fmt.Printf("Error scanning order: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		orders = append(orders, o)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
	fmt.Println("getOrdersHandler completed successfully")
}

// обработчик для получения заказа по ID
func getOrderByIDHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getOrderByIDHandler called")
	idStr := r.URL.Path[len("/orders/"):]

	var o Order
	err := db.QueryRow("SELECT id, product_id, quantity, total, created_at FROM orders WHERE id = $1", idStr).Scan(&o.ID, &o.ProductID, &o.Quantity, &o.Total, &o.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("Order not found: %v\n", idStr)
			http.Error(w, "Order not found", http.StatusNotFound)
		} else {
			fmt.Printf("Error querying order by ID: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(o)
	fmt.Println("getOrderByIDHandler completed successfully")
}

func main() {
	fmt.Println("Initializing database connection")
	initDB()
	defer db.Close()

	http.HandleFunc("/products", getProductsHandler)
	http.HandleFunc("/products/create", createProductHandler)
	http.HandleFunc("/products/", getProductByIDHandler)
	http.HandleFunc("/products/update/", updateProductHandler)
	http.HandleFunc("/products/delete/", deleteProductHandler)

	http.HandleFunc("/orders", getOrdersHandler)
	http.HandleFunc("/orders/create", createOrdersHandler)
	http.HandleFunc("/orders/", getOrderByIDHandler)

	fmt.Println("Server is running on port 8080!")
	http.ListenAndServe(":8080", nil)
}
