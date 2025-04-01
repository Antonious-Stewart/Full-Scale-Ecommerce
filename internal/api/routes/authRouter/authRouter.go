package auth_router

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Antonious-Stewart/Full-Scale-Ecommerce/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
	"time"
)

type AuthRouter struct {
	DB        db.DB
	Validator *validator.Validate
}

type registerRequest struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,gte=12"`
	Phone     string `json:"phone" validate:"required"`
}

func New(dbPool db.DB, validate *validator.Validate) *AuthRouter {
	return &AuthRouter{
		DB:        dbPool,
		Validator: validate,
	}
}

func (a *AuthRouter) Routes() http.Handler {
	router := chi.NewRouter()

	router.Post("/register", a.registerCustomer)

	return router
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), err
}

func generateToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	//secretKey := []byte("secret-tbd")
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}
	tokenString, err := token.SignedString(key)

	if err != nil {
		return "", err
	}

	return tokenString, err
}

func sliceToSQLString(strArr []string) string {
	return "{" + strings.Join(strArr, ",") + "}"
}

func (a *AuthRouter) registerCustomer(w http.ResponseWriter, r *http.Request) {
	// get customer data from the request body
	pool := a.DB.GetConnection()
	var req registerRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate data
	err = a.Validator.Struct(req)

	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			for _, e := range validateErrs {
				fmt.Println(e.Namespace())
				fmt.Println(e.Field())
				fmt.Println(e.StructNamespace())
				fmt.Println(e.StructField())
				fmt.Println(e.Tag())
				fmt.Println(e.ActualTag())
				fmt.Println(e.Kind())
				fmt.Println(e.Type())
				fmt.Println(e.Value())
				fmt.Println(e.Param())
				fmt.Println()
			}
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := uuid.New()
	hashedPassword, err := hashPassword(req.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := generateToken(userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tokens := sliceToSQLString([]string{token})

	// insert the data into the db
	query := "INSERT into customers (id, first_name, last_name, email, phone, password, tokens) values($1, $2, $3, $4, $5, $6, $7)"

	fmt.Println(query)
	_, err = pool.ExecContext(r.Context(), query, userID, req.FirstName, req.LastName, req.Email, req.Phone, hashedPassword, tokens)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Timeout", time.Now().Add(time.Second*7).String())
	w.Header().Add("Authorization", "Bearer "+token)

	w.WriteHeader(http.StatusCreated)
}
