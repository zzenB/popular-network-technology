package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
)

func baseHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Welcome. Please go to \"localhost:8080/hello\" or \"localhost:8080/pi\"")
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func piDigits(w http.ResponseWriter, r *http.Request) {
	pi := fmt.Sprintf("%.10f", math.Pi)
  fmt.Fprintf(w, "10 digits of pi is: %s", pi)
}

func main() {
	http.HandleFunc("/", baseHandler)
	http.HandleFunc("/hello", helloWorld)
	http.HandleFunc("/pi", piDigits)

	fmt.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
