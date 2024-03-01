package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var passwdDb = map[string]string{}
var serverSecret string
var apiKey string
var encryptionCost int

func init() {
	users := strings.Split(os.Getenv("AUTHORIZED_USERS"), ",")

	serverSecret = os.Getenv("SERVER_SECRET")
	if serverSecret == "" {
		log.Fatal("SERVER_SECRET is not defined")
	}

	apiKey = os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY not defined")
	}

	encryptionCostString := os.Getenv("ENCRYPTION_BCRYPT_COST")
	if encryptionCostString == "" {
		log.Fatal("ENCRYPTION_BCRYPT_COST is not defined")
	}
	var err error
	encryptionCost, err = strconv.Atoi(encryptionCostString)
	if err != nil {
		log.Fatal("ENCRYPTION_BCRYPT_COST is not a valid integer")
	}

	for _, user := range users {
		parts := strings.Split(user, ":")
		if len(parts) != 2 {
			continue
		}
		if hash, err := bcrypt.GenerateFromPassword([]byte(parts[1]), encryptionCost); err != nil {
			panic(err)
		} else {
			passwdDb[parts[0]] = string(hash)
		}
	}

	if len(passwdDb) == 0 {
		log.Fatal("AUTHORIZED_USERS is not defined with any valid users")
	}

	log.Info("user db ready")
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			log.Error("bad request", err)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("bad request"))
			return
		}

		if _, ok := passwdDb[body["username"]]; !ok {
			log.Warn("user not found in db")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("bad request"))
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(passwdDb[body["username"]]), []byte(body["password"])); err != nil {
			log.Error("error comparing hash and password", err)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("bad request"))
			return
		}
		jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": body["username"],
			"nbf":      time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		})

		tokenString, err := jwt.SignedString([]byte(serverSecret))
		if err != nil {
			log.Error("error signing token", err)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("bad request"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(tokenString))
	})

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		token, err := jwt.Parse(r.Header.Get("token"), func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(serverSecret), nil
		})
		if err != nil {
			log.Error("error parsing token", err)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("bad request"))
			return
		}

		if err, ok := token.Claims.(jwt.MapClaims); !ok {
			log.Error("invalid claims", err)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("bad request"))
			return
		}

		w.Header().Set("authorization", fmt.Sprintf("Bearer %s", apiKey))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	if err := http.ListenAndServe(":8000", r); err != nil {
		panic(err)
	}
}
