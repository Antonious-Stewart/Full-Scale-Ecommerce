package auth_router

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Antonious-Stewart/Full-Scale-Ecommerce/internal/db"
	"github.com/Antonious-Stewart/Full-Scale-Ecommerce/internal/shared"
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

type authRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=12"`
}

func New(dbPool db.DB, validate *validator.Validate) *AuthRouter {
	return &AuthRouter{
		DB:        dbPool,
		Validator: validate,
	}
}

func (a *AuthRouter) Routes() http.Handler {
	router := chi.NewRouter()

	router.Use(a.validateAuthReqBody)
	router.Post("/register", a.registerCustomer)
	router.With(a.validateUserCredentials).Post("/login", a.loginCustomer)

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

func (a *AuthRouter) validateAuthReqBody(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var req authRequest
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
			ctx := context.WithValue(r.Context(), "requestBody", req)
			next.ServeHTTP(w, r.WithContext(ctx))
		})

}
func (a *AuthRouter) validateUserCredentials(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		pool := a.DB.GetConnection()
		req := ctx.Value("requestBody").(authRequest)
		query := "SELECT email, password, id, first_name, last_name, phone, tokens, order_history from customers where email = $1"
		var user shared.AuthEntity
		err := pool.QueryRowContext(ctx, query, req.Email).Scan(
			&user.Email,
			&user.Password,
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Phone,
			&user.Tokens,
			&user.OrderHistory,
		)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Println(err)
				http.Error(w, "no records found with that user/password combination", http.StatusNotFound)
				return
			}

			log.Println(err)

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = bcrypt.CompareHashAndPassword(user.Password, []byte(req.Password))

		if err != nil {
			log.Println(err)

			http.Error(w, "No records found with that user/password combination", http.StatusNotFound)
			return
		}

		ctx = context.WithValue(ctx, "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *AuthRouter) registerCustomer(w http.ResponseWriter, r *http.Request) {
	// get customer data from the request body
	pool := a.DB.GetConnection()
	ctx := r.Context()
	req := ctx.Value("requestBody").(authRequest)

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
	query := "INSERT into customers (id, email, password, tokens) values($1, $2, $3, $4)"

	_, err = pool.ExecContext(ctx, query, userID, req.Email, hashedPassword, tokens)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Timeout", time.Now().Add(time.Second*7).String())
	w.Header().Add("Authorization", "Bearer "+token)

	w.WriteHeader(http.StatusCreated)
}

func (a *AuthRouter) loginCustomer(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")

	fmt.Println(user)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	return
	// validate credentials

	// sign new token
	// return
}
