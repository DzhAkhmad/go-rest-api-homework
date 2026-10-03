package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта
// 1 Обработчик для получения всех задач
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	// сериализуем мапу в байты json
	data, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}

	// записываем заголовок, статус и тело
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

// 2 Обработчик для отправки задачи на сервер
func taskPostHandler(w http.ResponseWriter, r *http.Request) {
	// создаем экземпляр структуры Task
	var task Task
	// создаем контейнер хранения байтов task из клиентского запроса
	var buf bytes.Buffer	
	// читаем данные из тела запроса, складываем их в буффер (buf)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// десериализуем буфферные данные из байтов в данные структуры Task
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// добавляем данные в мапу
	tasks[task.ID] = task
	// записываем заголовок и код ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// 3 Обработчик для получения задачи по ID
func getTaskIDHandler (w http.ResponseWriter, r *http.Request) {
	// получаем id из URL клиентского запроса
	id := chi.URLParam(r, "id")
	// проверяем есть ли такой элемент в мапе
	task, ok := tasks[id]
	if !ok {
		http.Error(w, "Элемент с таким id отсутствует", http.StatusBadRequest)
		return
	}
	// сериализуем task в байты json
	data, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// записываем заголовок статус и тело ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// 4 Обработчик для удаления записи task по ID
func deleteTaskIDHandler(w http.ResponseWriter, r *http.Request) {
	// получаем id из URL клиентского запроса
	id := chi.URLParam(r, "id")
	// проверяем есть ли такой элемент в мапе
	_, ok := tasks[id]
	if !ok {
		http.Error(w, "Элемент с таким id отсутствует", http.StatusBadRequest)
		return
	}
	// удаляем элемент с полученным id из мапы
	delete(tasks, id)
	// записываем заголовок и статус ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	// 1 обработчик
	r.Get("/tasks", getTaskHandler)
	// 2 обработчик
	r.Post("/tasks", taskPostHandler)
	// 3 обработчик
	r.Get("/tasks/{id}", getTaskIDHandler)
	// 4 обработчик
	r.Delete("/tasks/{id}", deleteTaskIDHandler)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
