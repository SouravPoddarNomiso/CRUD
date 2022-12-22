package main

import (
	"demo/router"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("This is a test")
	r := router.Router()
	http.ListenAndServe(":4000", r)
	fmt.Println("Listening at port 4000 ")
}
