package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	fiber "github.com/gofiber/fiber/v3"
)

const K = 1025
const SIZE = 2 * K

const MESSAGE = "Hello, World!"

type Msg struct {
	Message string `json:"message"`
}

var pool = sync.Pool{
	New: func() interface{} {
		arr := make([]Msg, SIZE)
		for i := 0; i < SIZE; i++ {
			arr[i].Message = MESSAGE
		}
		return &arr
	},
}

func buildMessageArray(num int) []Msg {
	arr := make([]Msg, num)
	for i := 0; i < num; i++ {
		arr[i].Message = MESSAGE
	}
	return arr
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, MESSAGE)
}

func httpHandlerJson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	msg := new(Msg)
	msg.Message = "Hello World"
	json.NewEncoder(w).Encode(msg)
}

func httpHandlerArray(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	arr := buildMessageArray(SIZE)
	json.NewEncoder(w).Encode(arr)
}

func httpHandlerArrayWithPool(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var arr = pool.Get().(*[]Msg)
	defer pool.Put(arr)
	json.NewEncoder(w).Encode(arr)
}

func fiberHandler(c fiber.Ctx) error {
	return c.SendString(MESSAGE)
}

func fiberHandlerJson(c fiber.Ctx) error {
	msg := new(Msg)
	msg.Message = "Hello World"
	return c.JSON(msg)
}

func fiberHandlerArray(c fiber.Ctx) error {
	arr := buildMessageArray(SIZE)
	return c.JSON(arr)
}

func main() {
	var wg = sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		http.HandleFunc("/", httpHandler)
		http.HandleFunc("/json", httpHandlerJson)
		http.HandleFunc("/array", httpHandlerArray)
		http.HandleFunc("/array_with_pool", httpHandlerArrayWithPool)
		log.Println("http started on http://127.0.0.1:3000")
		log.Fatal(http.ListenAndServe(":3000", nil))
	}()

	go func() {
		defer wg.Done()
		app := fiber.New()

		app.Get("/", fiberHandler)
		app.Get("/json", fiberHandlerJson)
		app.Get("/array", fiberHandlerArray)
		log.Fatal(app.Listen(":3001"))
	}()

	log.Println("Quit")

	wg.Wait()
}
