package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type User struct {
	// annotation
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreateAt  time.Time `json:"created_at"`
}

func main() {
	user := new(User)
	user.FirstName = "chaewon"
	user.LastName = "lee"
	user.Email = "chaewon@gmail.com"
	byte, _ := json.Marshal(user)
	buff := bytes.NewBuffer(byte)

	// Request 객체 생성
	req, err := http.NewRequest("POST", "http://10.0.5.86:3000/foo", buff)
	if err != nil {
		panic(err)
	}

	//Content-Type 헤더 추가
	req.Header.Add("Content-Type", "application/json")

	// Client객체에서 Request 실행
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Response 체크
	response, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		str := string(response)
		println(str)
	}

}
