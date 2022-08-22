package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Question struct {
	UUID string `json:"uuid"`
	A    int    `json:"a"`
	B    int    `json:"b"`
}

type Answer struct {
	UUID string `json:"uuid"`
	Sum  int    `json:"sum"`
}

type Response struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

var timeout = time.Minute
var DB sync.Map

func handler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		q := Question{
			UUID: uuid.NewString(),
			A:    rand.Intn(1 << 30),
			B:    rand.Intn(1 << 30),
		}
		qRes, _ := json.Marshal(q)

		DB.Store(q.UUID, q.A+q.B)
		time.AfterFunc(timeout, func() { DB.Delete(q.UUID) })
		w.Write(qRes)

	case "POST":
		A := Answer{}
		err := json.NewDecoder(req.Body).Decode(&A)
		if err != nil {
			errRes, _ := json.Marshal(Response{false, err.Error()})
			w.Write(errRes)
			return
		}

		sum, ok := DB.LoadAndDelete(A.UUID)
		if !ok {
			errRes, _ := json.Marshal(Response{false, "uuid not found"})
			w.Write(errRes)
			return
		}

		if sum.(int) != A.Sum {
			res, _ := json.Marshal(Response{false, "wrong answer"})
			w.Write(res)
		} else {
			res, _ := json.Marshal(Response{true, ""})
			w.Write(res)
		}
	}
}

func main() {
	http.HandleFunc("/tutorial", handler)
	http.ListenAndServe(":8080", nil)
}
