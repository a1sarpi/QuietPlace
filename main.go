package main

import (
	"log"
	"net/http"
	"os"

	"github.com/a1srapi/QuietPlace/handlers"
)

func main() {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	hh := handlers.NewHello(l)
	gh := handlers.Goodbye(l)

	sm := http.NewServeMux()
	sm.Handler("/", hh)
	sm.Handler("/goodbye", gh)

	http.ListenAndServe(":9090", nil)
}
