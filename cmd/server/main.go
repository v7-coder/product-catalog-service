package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Server starting on :8080")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		test := "test"
		test2 := "test"
		test3 := "sdsd"
		test4 := "sdsd"

		fmt.Println(test, test2, test3, test4)
		fmt.Fprintf(w, "Hello from product-catalog-service!")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
