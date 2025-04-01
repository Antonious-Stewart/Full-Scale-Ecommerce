package main

import (
	"errors"
	"fmt"
	"github.com/Antonious-Stewart/Full-Scale-Ecommerce/internal/api"
	"github.com/Antonious-Stewart/Full-Scale-Ecommerce/internal/config"
	"github.com/Antonious-Stewart/Full-Scale-Ecommerce/internal/db"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
)

var validate *validator.Validate

func main() {
	validate = validator.New()
	port, err := config.Get("PORT")

	if err != nil {
		log.Fatal(fmt.Errorf("failed to get PORT: %w", err))
	}

	pgUrl, err := config.Get("PG_DATABASE_URL")

	if err != nil {
		log.Fatal(fmt.Errorf("failed to get PG_DATABASE_URL: %w", err))
	}

	dbPool, err := db.NewShopDb(pgUrl)

	if err != nil {
		log.Fatal(fmt.Errorf("failed to create a new shop db instance: %w", err))
	}

	serve := http.Server{
		Addr:    ":" + port,
		Handler: api.New(dbPool, validate).Routes(),
	}

	log.Println("Listening on port: " + port)
	err = serve.ListenAndServe()

	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}
}
