package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt"
	"github.com/williamhaley/notes/api/data_store"
	"github.com/williamhaley/notes/api/middleware"
	"github.com/williamhaley/notes/api/sockets"
	"github.com/williamhaley/notes/api/www"
	"golang.org/x/crypto/bcrypt"
)

var passwdDb = map[string]string{}
var serverSecret string
var encryptionCost int

func init() {
	users := strings.Split(os.Getenv("AUTHORIZED_USERS"), ",")

	serverSecret = os.Getenv("SERVER_SECRET")
	if serverSecret == "" {
		log.Fatal("SERVER_SECRET is not defined")
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
	r.Use(chiMiddleware.Logger)

	socketServer := sockets.NewHub()
	go socketServer.Run()

	store := data_store.New()

	// Authenticated requests
	r.Group(func(r chi.Router) {
		r.Use(middleware.New(serverSecret))

		r.Delete("/api/notes/{id}", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			if err := store.DeleteNote(id); err != nil {
				log.Error("error deleting note", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("error deleting note"))
				return
			}

			socketServer.Broadcast("delete", id)

			w.WriteHeader(http.StatusOK)
		})

		r.Post("/api/notes", func(w http.ResponseWriter, r *http.Request) {
			var createNoteRequest data_store.CreateNoteRequest
			var err error
			var note *data_store.Note

			if err := json.NewDecoder(r.Body).Decode(&createNoteRequest); err != nil {
				log.Error("invalid note", err)
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("invalid note"))
				return
			}

			if note, err = store.CreateNote(&createNoteRequest); err != nil {
				log.Error("error saving", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("error saving"))
				return
			}

			socketServer.Broadcast("update", note)

			if err := json.NewEncoder(w).Encode(note); err != nil {
				log.Error("error serializing", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("error serializing"))
			}
		})

		r.Put("/api/notes", func(w http.ResponseWriter, r *http.Request) {
			var updateNoteRequest data_store.UpdateNoteRequest
			var err error
			var note *data_store.Note

			if err := json.NewDecoder(r.Body).Decode(&updateNoteRequest); err != nil {
				log.Error("invalid note", err)
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("invalid note"))
				return
			}

			if note, err = store.UpdateNote(&updateNoteRequest); err != nil {
				log.Error("error saving", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("error saving"))
				return
			}

			socketServer.Broadcast("update", note)

			if err := json.NewEncoder(w).Encode(note); err != nil {
				log.Error("error serializing", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("error serializing"))
			}
		})

		r.Get("/api/notes", func(w http.ResponseWriter, r *http.Request) {
			notes, err := store.GetNotes()
			if err != nil {
				log.Error("error getting notes", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("error getting notes"))
				return
			}

			w.Header().Set("Content-Type", "application/json")
			if err = json.NewEncoder(w).Encode(map[string]any{
				"results": notes,
			}); err != nil {
				log.Error("error getting notes", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("error getting notes"))
				return
			}
		})

		// Validate the current auth token
		r.Get("/api/check", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		})
	})

	r.Get("/socket", func(w http.ResponseWriter, r *http.Request) {
		conn, err := sockets.Upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error("unknown issue", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("unknown issue"))
			return
		}

		sockets.NewClient(socketServer, conn)
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	r.Post("/api/login", func(w http.ResponseWriter, r *http.Request) {
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

	fileServer := http.FileServer(http.FS(www.GetContent()))
	r.Get("/*", fileServer.ServeHTTP)

	if err := http.ListenAndServe(":8000", r); err != nil {
		panic(err)
	}
}
