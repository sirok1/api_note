package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Product представляет продукт
type Product struct {
	ID          int
	ImageURL    string
	Name        string
	Description string
	Price       float64
}

// Пример списка продуктов
var products = []Product{
	{ID: 1, ImageURL: "https://example.com/image1.jpg", Name: "Product One", Description: "This is the first product description.", Price: 19.99},
	{ID: 2, ImageURL: "https://example.com/image2.jpg", Name: "Product Two", Description: "This is the second product description.", Price: 29.99},
	{ID: 3, ImageURL: "https://example.com/image3.jpg", Name: "Product Three", Description: "This is the third product description.", Price: 39.99},
	{ID: 4, ImageURL: "https://example.com/image4.jpg", Name: "Product Four", Description: "This is the fourth product description.", Price: 49.99},
	{ID: 5, ImageURL: "https://example.com/image5.jpg", Name: "Product Five", Description: "This is the fifth product description.", Price: 59.99},
	{ID: 6, ImageURL: "https://example.com/image6.jpg", Name: "Product Six", Description: "This is the sixth product description.", Price: 15.99},
	{ID: 7, ImageURL: "https://example.com/image7.jpg", Name: "Product Seven", Description: "This is the seventh product description.", Price: 25.99},
	{ID: 8, ImageURL: "https://example.com/image8.jpg", Name: "Product Eight", Description: "This is the eighth product description.", Price: 35.99},
	{ID: 9, ImageURL: "https://example.com/image9.jpg", Name: "Product Nine", Description: "This is the ninth product description.", Price: 45.99},
	{ID: 10, ImageURL: "https://example.com/image10.jpg", Name: "Product Ten", Description: "This is the tenth product description.", Price: 55.99},
	{ID: 11, ImageURL: "https://example.com/image11.jpg", Name: "Product Eleven", Description: "This is the eleventh product description.", Price: 22.49},
	{ID: 12, ImageURL: "https://example.com/image12.jpg", Name: "Product Twelve", Description: "This is the twelfth product description.", Price: 32.49},
	{ID: 13, ImageURL: "https://example.com/image13.jpg", Name: "Product Thirteen", Description: "This is the thirteenth product description.", Price: 42.49},
	{ID: 14, ImageURL: "https://example.com/image14.jpg", Name: "Product Fourteen", Description: "This is the fourteenth product description.", Price: 52.49},
	{ID: 15, ImageURL: "https://example.com/image15.jpg", Name: "Product Fifteen", Description: "This is the fifteenth product description.", Price: 62.49},
}

// обработчик для GET-запроса, возвращает список продуктов
func getProductsHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовки для правильного формата JSON
	w.Header().Set("Content-Type", "application/json")
	// Преобразуем список заметок в JSON
	json.NewEncoder(w).Encode(products)
}

// обработчик для POST-запроса, добавляет продукт
func createProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var newProduct Product
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		fmt.Println("Error decoding request body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Received new Product: %+v\n", newProduct)

	newProduct.ID = len(products) + 1
	products = append(products, newProduct)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newProduct)
}

//Добавление маршрута для получения одного продукта

func getProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из URL
	idStr := r.URL.Path[len("/Products/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Product ID", http.StatusBadRequest)
		return
	}

	// Ищем продукт с данным ID
	for _, Product := range products {
		if Product.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Product)
			return
		}
	}

	// Если продукт не найден
	http.Error(w, "Product not found", http.StatusNotFound)
}

// удаление продукта по id
func deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID из URL
	idStr := r.URL.Path[len("/Products/delete/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Product ID", http.StatusBadRequest)
		return
	}

	// Ищем и удаляем продукт с данным ID
	for i, Product := range products {
		if Product.ID == id {
			// Удаляем продукт из среза
			products = append(products[:i], products[i+1:]...)
			w.WriteHeader(http.StatusNoContent) // Успешное удаление, нет содержимого
			return
		}
	}

	// Если продукт не найден
	http.Error(w, "Product not found", http.StatusNotFound)
}

// Обновление продукта по id
func updateProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID из URL
	idStr := r.URL.Path[len("/Products/update/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Product ID", http.StatusBadRequest)
		return
	}

	// Декодируем обновлённые данные продукта
	var updatedProduct Product
	err = json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Ищем продукт для обновления
	for i, Product := range products {
		if Product.ID == id {

			products[i].ImageURL = updatedProduct.ImageURL
			products[i].Name = updatedProduct.Name
			products[i].Description = updatedProduct.Description
			products[i].Price = updatedProduct.Price

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(products[i])
			return
		}
	}

	// Если продукт не найден
	http.Error(w, "Product not found", http.StatusNotFound)
}

func main() {
	http.HandleFunc("/Products", getProductsHandler)           // Получить все продукты
	http.HandleFunc("/Products/create", createProductHandler)  // Создать продукт
	http.HandleFunc("/Products/", getProductByIDHandler)       // Получить продукт по ID
	http.HandleFunc("/Products/update/", updateProductHandler) // Обновить продукт
	http.HandleFunc("/Products/delete/", deleteProductHandler) // Удалить продукт

	fmt.Println("Server is running on port 8080!")
	http.ListenAndServe(":8080", nil)
}
