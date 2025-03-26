package main

import (
	"errors"
	"github.com/Antonious-Stewart/Full-Scale-Ecommerce/internal/config"
	"log"
	"net/http"
	"time"
)

func main() {
	port, err := config.Get("PORT")

	if err != nil {
		log.Fatal(err)
	}

	serve := http.Server{
		Addr:         ":" + port,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	err = serve.ListenAndServe()

	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}
}
