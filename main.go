package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Note представляет заметку
type Note struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Пример списка заметок
var notes = []Note{
	{ID: 1, Title: "First Note", Content: "This is the content of the first note."},
	{ID: 2, Title: "Second Note", Content: "This is the content of the second note."},
	{ID: 3, Title: "Third Note", Content: "This is the content of the third note."},
}

// обработчик для GET-запроса, возвращает список заметок
func getNotesHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовки для правильного формата JSON
	w.Header().Set("Content-Type", "application/json")
	// Преобразуем список заметок в JSON
	json.NewEncoder(w).Encode(notes)
}

// обработчик для POST-запроса, добавляет заметку
func createNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var newNote Note
	err := json.NewDecoder(r.Body).Decode(&newNote)
	if err != nil {
		fmt.Println("Error decoding request body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Received new note: %+v\n", newNote)

	newNote.ID = len(notes) + 1
	notes = append(notes, newNote)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newNote)
}

//Добавление маршрута для получения одной заметки

func getNoteByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из URL
	idStr := r.URL.Path[len("/notes/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	// Ищем заметку с данным ID
	for _, note := range notes {
		if note.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(note)
			return
		}
	}

	// Если заметка не найдена
	http.Error(w, "Note not found", http.StatusNotFound)
}

// удаление заметок по id
func deleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID из URL
	idStr := r.URL.Path[len("/notes/delete/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	// Ищем и удаляем заметку с данным ID
	for i, note := range notes {
		if note.ID == id {
			// Удаляем заметку из среза
			notes = append(notes[:i], notes[i+1:]...)
			w.WriteHeader(http.StatusNoContent) // Успешное удаление, нет содержимого
			return
		}
	}

	// Если заметка не найдена
	http.Error(w, "Note not found", http.StatusNotFound)
}

// Обновление заметок по id
func updateNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID из URL
	idStr := r.URL.Path[len("/notes/update/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	// Декодируем обновлённые данные заметки
	var updatedNote Note
	err = json.NewDecoder(r.Body).Decode(&updatedNote)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Ищем заметку для обновления
	for i, note := range notes {
		if note.ID == id {
			notes[i].Title = updatedNote.Title
			notes[i].Content = updatedNote.Content
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(notes[i])
			return
		}
	}

	// Если заметка не найдена
	http.Error(w, "Note not found", http.StatusNotFound)
}

func main() {
	http.HandleFunc("/notes", getNotesHandler)           // Получить все заметки
	http.HandleFunc("/notes/create", createNoteHandler)  // Создать заметку
	http.HandleFunc("/notes/", getNoteByIDHandler)       // Получить заметку по ID
	http.HandleFunc("/notes/update/", updateNoteHandler) // Обновить заметку
	http.HandleFunc("/notes/delete/", deleteNoteHandler) // Удалить заметку

	fmt.Println("Server is running on port 8080!")
	http.ListenAndServe(":8080", nil)
}
