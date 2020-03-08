package web

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"my-app2/db"
	"net/http"
)

type App struct {
	d        db.DB
	handlers map[string]http.HandlerFunc
}

// 追加
type TodoForm struct {
	Text string `json:"text"`
}

func NewApp(d db.DB, cors bool) App {
	app := App{
		d:        d,
		handlers: make(map[string]http.HandlerFunc),
	}
	techHandler := app.GetTechnologies
	// 追加
	todoHandler := app.Todos
	if !cors {
		techHandler = disableCors(techHandler)
		// 追加
		todoHandler = disableCors(todoHandler)
	}
	app.handlers["/api/technologies"] = techHandler
	// 追加
	app.handlers["/api/todos"] = todoHandler
	app.handlers["/"] = http.FileServer(http.Dir("/webapp")).ServeHTTP
	return app
}

func (a *App) Serve() error {
	for path, handler := range a.handlers {
		http.Handle(path, handler)
	}
	log.Println("Web server is available on port 8080")
	return http.ListenAndServe(":8080", nil)
}

func (a *App) GetTechnologies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	technologies, err := a.d.GetTechnologies()
	if err != nil {
		sendErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = json.NewEncoder(w).Encode(technologies)
	if err != nil {
		sendErr(w, http.StatusInternalServerError, err.Error())
	}
}

// 追加
func (a *App) Todos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		todos, err := a.d.GetTodos()
		if err != nil {
			sendErr(w, http.StatusInternalServerError, err.Error())
			return
		}
		err = json.NewEncoder(w).Encode(todos)
		if err != nil {
			sendErr(w, http.StatusInternalServerError, err.Error())
		}
	case http.MethodPost:
		body := r.Body
		defer body.Close()
		buf := new(bytes.Buffer)
		io.Copy(buf, body)

		var todo TodoForm
		json.Unmarshal(buf.Bytes(), &todo)
		err := a.d.PostTodos(todo.Text)
		if err != nil {
			sendErr(w, http.StatusInternalServerError, err.Error())
			return
		}

	case http.MethodDelete:
		err := a.d.DeleteTodos(r.FormValue("id"))
		if err != nil {
			sendErr(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}

func sendErr(w http.ResponseWriter, code int, message string) {
	resp, _ := json.Marshal(map[string]string{"error": message})
	http.Error(w, string(resp), code)
}

// Needed in order to disable CORS for local developazsment
func disableCors(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		h(w, r)
	}
}
