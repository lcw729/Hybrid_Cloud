package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type User struct {
	// annotation
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreateAt  time.Time `json:"created_at"`
}
type fooHandler struct{}

func (f *fooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := new(User)
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Bad Request", err)
		return
	}
	user.CreateAt = time.Now()

	data, _ := json.Marshal(user)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(data))
}

// request에서 필요한 argument값 넣기
// http:10.0.5.86:3000/bar?name=<my_name>
// http://localhost?id=1&name=abcd
func barHandler(w http.ResponseWriter, r *http.Request) {
	// URL에서 argument를 뽑아냄
	values := r.URL.Query()
	name := values.Get("name")
	if name == "" {
		name = "world"
	}
	id, _ := strconv.Atoi(values.Get("id"))
	fmt.Fprintf(w, "Hello %s! id:%d", name, id)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello World")
	})

	http.HandleFunc("/bar", barHandler)

	http.Handle("/foo", &fooHandler{})

	http.ListenAndServe(":3000", nil)
}
